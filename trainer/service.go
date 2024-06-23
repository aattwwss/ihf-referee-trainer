package trainer

import (
	"context"
	"golang.org/x/exp/rand"
	"slices"
)

type Repository interface {
	GetAllQuestions(ctx context.Context) ([]Question, error)
	GetAllRules(ctx context.Context) ([]Rule, error)
	GetQuestionByID(ctx context.Context, id int) (*Question, error)
	GetRandomQuestion(ctx context.Context, rules []string) (*Question, error)
	GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error)
	ListQuestions(ctx context.Context, rules []string, search string, lastRuleSortOrder int, lastQuestionNumber int, limit int) ([]Question, error)
	InsertFeedback(ctx context.Context, feedback Feedback) error
	InsertQuizConfig(ctx context.Context, config QuizConfig) (string, error)
	GetQuizConfigByKey(ctx context.Context, key string) (*QuizConfig, error)
}

type QuestionService struct {
	repository Repository
}

func NewService(repository Repository) *QuestionService {
	return &QuestionService{repository: repository}
}

func (s *QuestionService) GetAllQuestions(ctx context.Context) ([]Question, error) {
	return s.repository.GetAllQuestions(ctx)
}

func (s *QuestionService) GetAllRules(ctx context.Context) ([]Rule, error) {
	return s.repository.GetAllRules(ctx)
}

func (s *QuestionService) GetQuestionByID(ctx context.Context, id int) (*Question, error) {
	return s.repository.GetQuestionByID(ctx, id)
}

func (s *QuestionService) GetRandomQuestion(ctx context.Context, rules []string) (*Question, error) {
	return s.repository.GetRandomQuestion(ctx, rules)
}

func (s *QuestionService) GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error) {
	return s.repository.GetChoicesByQuestionID(ctx, questionID)
}

func (s *QuestionService) ListQuestions(ctx context.Context, rules []string, search string, lastRuleSortOrder int, lastQuestionNumber int, limit int) ([]Question, error) {
	return s.repository.ListQuestions(ctx, rules, search, lastRuleSortOrder, lastQuestionNumber, limit)
}

func (s *QuestionService) SubmitFeedback(ctx context.Context, feedback Feedback) error {
	return s.repository.InsertFeedback(ctx, feedback)
}

func (s *QuestionService) SubmitQuizConfig(ctx context.Context, config QuizConfig) (string, error) {
	return s.repository.InsertQuizConfig(ctx, config)
}

func (s *QuestionService) GetQuestionsByQuizConfigKey(ctx context.Context, key string) ([]Question, error) {
	quizConfig, err := s.repository.GetQuizConfigByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	allQuestions, err := s.repository.GetAllQuestions(ctx)
	if err != nil {
		return nil, err
	}

	var questions []Question
	for _, q := range allQuestions {
		if slices.Contains(quizConfig.RuleIDs, q.Rule.ID) {
			questions = append(questions, q)
		}
	}

	rand.Seed(uint64(quizConfig.Seed))
	rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
	questions = questions[:quizConfig.NumQuestions]
	slices.SortFunc(questions, func(a, b Question) int {
		if a.Rule.SortOrder == b.Rule.SortOrder {
			return a.QuestionNumber - b.QuestionNumber
		}
		return a.Rule.SortOrder - b.Rule.SortOrder
	})
	return questions[:quizConfig.NumQuestions], nil
}

func (s *QuestionService) EvaluateQuizAnswer(ctx context.Context, quizConfigKey string, choiceIDs []int) ([]Question, error) {
	questions, err := s.GetQuestionsByQuizConfigKey(ctx, quizConfigKey)
	if err != nil {
		return nil, err
	}
	for _, q := range questions {
		for _, c := range q.Choices {
			isSelected := slices.Contains(choiceIDs, c.ID)
			if isSelected && c.IsAnswer {
				c.Result = refOf(ChoiceResultCorrect)
			} else if isSelected && !c.IsAnswer {
				c.Result = refOf(ChoiceResultWrong)
			} else if !isSelected && c.IsAnswer {
				c.Result = refOf(ChoiceResultMissing)
			}
		}
	}
	return questions, nil
}

func refOf[E any](e E) *E {
	return &e
}
