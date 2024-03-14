package trainer

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

