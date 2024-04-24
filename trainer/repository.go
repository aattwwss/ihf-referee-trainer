package trainer

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *QuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

func (r *QuestionRepository) GetRandomQuestion(ctx context.Context, rules []string) (*Question, error) {
	if len(rules) == 0 {
		allRules, err := r.GetAllDistinctRuleIDs(ctx)
		if err != nil {
			return nil, err
		}
		rules = allRules
	}
	query := fmt.Sprintf("SELECT id, text, rule_id, question_number FROM question WHERE rule_id = ANY($1) ORDER BY RANDOM() LIMIT 1")
	rows, err := r.db.Query(ctx, query, rules)
	if err != nil {
		return nil, err
	}
	questionEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[QuestionEntity])
	if err != nil {
		return nil, err
	}
	query = fmt.Sprintf("SELECT * FROM choice WHERE question_id = $1 order by option")
	rows, err = r.db.Query(ctx, query, questionEntity.ID)
	if err != nil {
		return nil, err
	}

	choiceEntities, err := pgx.CollectRows(rows, pgx.RowToStructByPos[ChoiceEntity])
	var choices []Choice
	for _, choiceEntity := range choiceEntities {
		choices = append(choices, Choice{
			ID:         choiceEntity.ID,
			Option:     choiceEntity.Option,
			Text:       choiceEntity.Text,
			IsAnswer:   choiceEntity.IsAnswer,
			IsSelected: false,
		})
	}
	separator := "."
	if questionEntity.RuleID == "SAR" {
		separator = ""
	}
	ruleQuestionNumber := fmt.Sprintf("%s%s%d", questionEntity.RuleID, separator, questionEntity.QuestionNumber)
	rule, err := r.FindRuleByID(ctx, questionEntity.RuleID)
	if err != nil {
		return nil, err
	}
	return &Question{
		ID:                 questionEntity.ID,
		Text:               questionEntity.Text,
		Rule:               *rule,
		QuestionNumber:     questionEntity.QuestionNumber,
		RuleQuestionNumber: ruleQuestionNumber,
		Choices:            choices,
	}, nil
}

func (r *QuestionRepository) GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error) {
	query := fmt.Sprintf("SELECT * FROM choice WHERE question_id = $1 order by option")
	rows, err := r.db.Query(ctx, query, questionID)
	if err != nil {
		return nil, err
	}
	choiceEntities, err := pgx.CollectRows(rows, pgx.RowToStructByPos[ChoiceEntity])
	if err != nil {
		return nil, err
	}
	var choices []Choice
	for _, choiceEntity := range choiceEntities {
		choices = append(choices, Choice{
			ID:         choiceEntity.ID,
			Option:     choiceEntity.Option,
			Text:       choiceEntity.Text,
			IsAnswer:   choiceEntity.IsAnswer,
			IsSelected: false,
		})
	}
	return choices, nil
}

func (r *QuestionRepository) FindRuleByID(ctx context.Context, ruleID string) (*Rule, error) {
	rulesMap, err := r.FindRuleByIDs(ctx, ruleID)
	if err != nil {
		return nil, err
	}
	if rule, ok := rulesMap[ruleID]; ok {
		return &rule, nil
	}
	return nil, nil
}

// FindRuleByIDs finds rules by rule ids and returns a map of rule id to rule
func (r *QuestionRepository) FindRuleByIDs(ctx context.Context, ruleIDs ...string) (map[string]Rule, error) {
	query := fmt.Sprintf("SELECT * FROM rule WHERE id = ANY($1)")
	rows, err := r.db.Query(ctx, query, ruleIDs)
	if err != nil {
		return nil, err
	}
	ruleEntity, err := pgx.CollectRows(rows, pgx.RowToStructByPos[RuleEntity])
	if err != nil {
		return nil, err
	}
	var rulesMap = make(map[string]Rule)
	for _, ruleEntity := range ruleEntity {
		rulesMap[ruleEntity.ID] = Rule{
			ID:        ruleEntity.ID,
			Name:      ruleEntity.Name,
			SortOrder: ruleEntity.SortOrder,
		}
	}
	return rulesMap, nil
}

func (r *QuestionRepository) GetAllDistinctRuleIDs(ctx context.Context) ([]string, error) {

	query := fmt.Sprintf("SELECT id FROM rule")
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	// scan int string array
	var rules []string
	for rows.Next() {

		var rule string
		err = rows.Scan(&rule)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

// ListQuestions returns a list of questions
// supports pagination using the rule sort order and question number of the last question to offset
func (r *QuestionRepository) ListQuestions(ctx context.Context, ruleIDs []string, search string, lastRuleSortOrder int, lastQuestionNumber int, limit int) ([]Question, error) {
	if len(ruleIDs) == 0 {
		allRules, err := r.GetAllDistinctRuleIDs(ctx)
		if err != nil {
			return nil, err
		}
		ruleIDs = allRules
	}
	query := fmt.Sprintf(`
		SELECT q.id, q.text, q.rule_id, q.question_number 
		FROM question q join rule r on q.rule_id = r.id 
		WHERE r.id = ANY($1) 
			AND tsv @@ websearch_to_tsquery($2) 
			AND r.sort_order >= $3
			AND q.question_number > $4
		ORDER BY r.sort_order, q.question_number 
		LIMIT $5
	`)
	rows, err := r.db.Query(ctx, query, ruleIDs, search, lastRuleSortOrder, lastQuestionNumber, limit)
	if err != nil {
		return nil, err
	}
	questionEntities, err := pgx.CollectRows(rows, pgx.RowToStructByPos[QuestionEntity])
	if err != nil {
		return nil, err
	}
	questionIds := make([]int, 0, len(questionEntities))
	for _, questionEntity := range questionEntities {
		questionIds = append(questionIds, questionEntity.ID)
	}
	choiceMap, err := r.FindChoicesByQuestionIds(ctx, questionIds...)
	rulesMap, err := r.FindRuleByIDs(ctx, ruleIDs...)
	var questions []Question
	for _, questionEntity := range questionEntities {
		separator := "."
		if questionEntity.RuleID == "SAR" {
			separator = ""
		}
		ruleQuestionNumber := fmt.Sprintf("%s%s%d", questionEntity.RuleID, separator, questionEntity.QuestionNumber)
		question := Question{
			ID:                 questionEntity.ID,
			Text:               questionEntity.Text,
			Rule:               rulesMap[questionEntity.RuleID],
			QuestionNumber:     questionEntity.QuestionNumber,
			RuleQuestionNumber: ruleQuestionNumber,
			Choices:            choiceMap[questionEntity.ID],
		}
		questions = append(questions, question)
	}
	return questions, nil
}

// FindChoicesByQuestionIds finds choices by question ids and returns a map of question id to choices
func (r *QuestionRepository) FindChoicesByQuestionIds(ctx context.Context, questionIds ...int) (map[int][]Choice, error) {
	query := fmt.Sprintf("SELECT * FROM choice WHERE question_id = ANY($1) order by option")
	rows, err := r.db.Query(ctx, query, questionIds)
	if err != nil {
		return nil, err
	}
	choiceEntities, err := pgx.CollectRows(rows, pgx.RowToStructByPos[ChoiceEntity])
	if err != nil {
		return nil, err
	}
	var choiceMap = make(map[int][]Choice)
	for _, choiceEntity := range choiceEntities {
		choices, ok := choiceMap[choiceEntity.QuestionId]
		if !ok {
			choices = []Choice{}
		}

		choices = append(choices, Choice{
			ID:         choiceEntity.ID,
			Option:     choiceEntity.Option,
			Text:       choiceEntity.Text,
			IsAnswer:   choiceEntity.IsAnswer,
			IsSelected: false,
		})
		choiceMap[choiceEntity.QuestionId] = choices
	}
	return choiceMap, nil
}
