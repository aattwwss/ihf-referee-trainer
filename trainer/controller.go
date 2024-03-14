package trainer

type Service interface {
	GetRandomQuestion(ctx context.Context) (Question, error)
}

type Controller struct {
}
