package trainer

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *QuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

func (r *QuestionRepository) GetRandomQuestion(ctx context.Context) (*Question, error) {
	query := fmt.Sprintf("SELECT * FROM question ORDER BY RANDOM() LIMIT 1")
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	questionEntity, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[QuestionEntity])
	if err != nil {
		return nil, err
	}
	query = fmt.Sprintf("SELECT * FROM choice WHERE question_id = $1 order by option")
	rows, err = r.db.Query(ctx, query, questionEntity.ID)
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
