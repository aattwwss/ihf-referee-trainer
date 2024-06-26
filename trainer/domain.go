package trainer

type QuestionEntity struct {
	ID             int
	Text           string
	RuleID         string
	QuestionNumber int
}

type ChoiceEntity struct {
	ID         int
	QuestionId int
	Option     string
	Text       string
	IsAnswer   bool
}

type RuleEntity struct {
	ID        string
	Name      string
	SortOrder int
}

type ReferenceEntity struct {
	ID         int
	QuestionId int
	Text       string
}

type FeedbackEntity struct {
	ID             int
	Name           string
	Email          string
	Topic          string
	Text           string
	IsAcknowledged bool
	IsCompleted    bool
}

type Question struct {
	ID                 int
	Text               string
	Rule               Rule
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

type Rule struct {
	ID        string
	Name      string
	SortOrder int
}

type Reference struct {
	ID   int
	Text string
}

type Feedback struct {
	ID             int
	Name           string
	Email          string
	Topic          string
	Text           string
	IsAcknowledged bool
	IsCompleted    bool
}
