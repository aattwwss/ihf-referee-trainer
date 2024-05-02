package trainer

import "context"

type Repository interface {
	GetQuestionByID(ctx context.Context, id int) (*Question, error)
	GetRandomQuestion(ctx context.Context, rules []string) (*Question, error)
	GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error)
	ListQuestions(ctx context.Context, rules []string, search string, lastRuleSortOrder int, lastQuestionNumber int, limit int) ([]Question, error)
}

type QuestionService struct {
	repository Repository
}

func NewService(repository Repository) *QuestionService {
	return &QuestionService{repository: repository}
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
