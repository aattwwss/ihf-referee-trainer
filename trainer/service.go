package trainer

import "context"

type Repository interface {
	GetRandomQuestion(ctx context.Context, rules []string) (*Question, error)
	GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error)
	ListQuestions(ctx context.Context, rules []string, search string, limit int) ([]Question, error)
}

type QuestionService struct {
	repository Repository
}

func NewService(repository Repository) *QuestionService {
	return &QuestionService{repository: repository}
}

func (s *QuestionService) GetRandomQuestion(ctx context.Context, rules []string) (*Question, error) {
	return s.repository.GetRandomQuestion(ctx, rules)
}

func (s *QuestionService) GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error) {
	return s.repository.GetChoicesByQuestionID(ctx, questionID)
}

func (s *QuestionService) ListQuestions(ctx context.Context, rules []string, search string, limit int) ([]Question, error) {
	return s.repository.ListQuestions(ctx, rules, search, limit)
}
