package trainer

import "context"

type Service interface {
	GetRandomQuestion(ctx context.Context) (*Question, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{
		service: service,
	}
}
