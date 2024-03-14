package trainer

type Repository interface {
	GetRandomQuestion(ctx context.Context) (Question, error)
}

type QuestionService struct {
	repository Repository
}

func NewService(repository Repository) *QuestionService {
	return QuestionService{repository: repository}
}

func (s *Service) GetRandomQuestion(ctx context.Context) (Question, error) {
	return s.repository.GetRandomQuestion(ctx)
}
