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

type QuestionRepository struct {
	db *pgxpool.Pool
}

func (r *QuestionRepository) GetRandomQuestion(ctx context.Context) (*Question, error) {
	rows, err := dbpool.Query(ctx, "SELECT * FROM question ORDER BY newid() LIMIT 1")
	if err != nil {
		return nil, err
	}
	questionEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[QuestionEntity])
	if err != nil {
		return nil, err
	}
	rows, err = dbpool.Query(ctx, "SELECT * FROM choice WHERE question_id = $1 order by option", question.ID)
	if err != nil {
		return nil, err
	}
	choiceEntities, err := pgx.CollectRows(rows, pgx.RowToStructByPos[ChoiceEntity])
	var choices []Choice
	for _, choiceEntity := range choiceEntities {
		choices = append(choices, Choice{
			ID:       choiceEntity.ID,
			Option:   choiceEntity.Option,
			Text:     choiceEntity.Text,
			IsAnswer: choiceEntity.IsAnswer,
		})
	}
	return &Question{
		ID:             questionEntity.ID,
		Text:           questionEntity.Text,
		Rule:           questionEntity.Rule,
		QuestionNumber: questionEntity.QuestionNumber,
		Choices:        choices,
	}, nil
}

type Repository interface {
	GetRandomQuestion(ctx context.Context) (Question, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return Service{repository: repository}
}

func (s *Service) GetRandomQuestion(ctx context.Context) (Question, error) {
	return s.repository.GetRandomQuestion(ctx)
}
