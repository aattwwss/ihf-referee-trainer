package trainer

import "context"

type Repository interface {
	GetAllQuestions(ctx context.Context) ([]Question, error)
	GetAllRules(ctx context.Context) ([]Rule, error)
	GetQuestionByID(ctx context.Context, id int) (*Question, error)
	GetRandomQuestion(ctx context.Context, rules []string) (*Question, error)
	GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error)
	ListQuestions(ctx context.Context, rules []string, search string, lastRuleSortOrder int, lastQuestionNumber int, limit int) ([]Question, error)
	InsertFeedback(ctx context.Context, feedback Feedback) error
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
