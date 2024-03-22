package trainer

import "context"

type Repository interface {
	GetRandomQuestion(ctx context.Context) (*Question, error)
	GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error)
}

type QuestionService struct {
	repository Repository
}

func NewService(repository Repository) *QuestionService {
	return &QuestionService{repository: repository}
}

func (s *QuestionService) GetRandomQuestion(ctx context.Context) (*Question, error) {
	return s.repository.GetRandomQuestion(ctx)
}

func (s *QuestionService) GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error) {
	return s.repository.GetChoicesByQuestionID(ctx, questionID)
}
