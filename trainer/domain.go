package trainer

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

type ReferenceEntity struct {
	ID         int
	QuestionId int
	Text       string
}

type Question struct {
	ID                 int
	Text               string
	Rule               string
	QuestionNumber     int
	RuleQuestionNumber string
	Choices            []Choice
	References         []Reference
}

type Choice struct {
	ID         int
	Option     string
	Text       string
	IsAnswer   bool
	IsSelected bool
}

type Reference struct {
	ID   int
	Text string
}
