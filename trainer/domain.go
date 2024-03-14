package trainer

import (
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionEntity struct {
	ID             int
	Text           string
	Rule           string
	QuestionNumber int
}

type ChoiceEntity struct {
	ID         int
	QuestionId int
	Option     string
	Text       string
	IsAnswer   bool
}

type Question struct {
	ID             int
	Text           string
	Rule           string
	QuestionNumber int
	Choices        []Choice
}

type Choice struct {
	ID       int
	Option   string
	Text     string
	IsAnswer bool
}
