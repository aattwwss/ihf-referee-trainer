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
		allRules, err := r.GetAllDistinctRules(ctx)
		if err != nil {
			return nil, err
		}
		rules = allRules
	}
	query := fmt.Sprintf("SELECT * FROM question WHERE rule =  ANY($1) ORDER BY RANDOM() LIMIT 1")
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
	if questionEntity.Rule == "SAR" {
		separator = ""
	}
	ruleQuestionNumber := fmt.Sprintf("%s%s%d", questionEntity.Rule, separator, questionEntity.QuestionNumber)
	return &Question{
		ID:                 questionEntity.ID,
		Text:               questionEntity.Text,
		Rule:               questionEntity.Rule,
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

func (r *QuestionRepository) GetAllDistinctRules(ctx context.Context) ([]string, error) {

	query := fmt.Sprintf("SELECT text FROM rule")
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

func (r *QuestionRepository) ListQuestions(ctx context.Context, rules []string, search string) ([]Question, error) {
	if len(rules) == 0 {
		allRules, err := r.GetAllDistinctRules(ctx)
		if err != nil {
			return nil, err
		}
		rules = allRules
	}
	query := fmt.Sprintf("SELECT * FROM question q join rule r on q.rule = r.text WHERE rule =  ANY($1) ORDER BY r.sort_number, q.question_number")
	rows, err := r.db.Query(ctx, query, rules)
	if err != nil {
		return nil, err
	}
	questionEntities, err := pgx.CollectRows(rows, pgx.RowToStructByPos[QuestionEntity])
	if err != nil {
		return nil, err
	}
	var questions []Question
	for _, questionEntity := range questionEntities {
		question := Question{
			ID:                 questionEntity.ID,
			Text:               questionEntity.Text,
			Rule:               questionEntity.Rule,
			QuestionNumber:     questionEntity.QuestionNumber,
			RuleQuestionNumber: ruleQuestionNumber,
			Choices:            choices,
		}
		questions = append(questions, question)
	}
	return questions, nil
}
