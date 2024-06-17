package trainer

import (
	"context"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

type Service interface {
	GetAllQuestions(ctx context.Context) ([]Question, error)
	GetQuestionByID(ctx context.Context, id int) (*Question, error)
	GetRandomQuestion(ctx context.Context, rules []string) (*Question, error)
	GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error)
	ListQuestions(ctx context.Context, rules []string, search string, lastRuleSortOrder int, lastQuestionNumber int, limit int) ([]Question, error)
	SubmitFeedback(ctx context.Context, feedback Feedback) error
}

type Controller struct {
	service Service
	html    fs.FS
}

func NewController(service Service, html fs.FS) *Controller {
	return &Controller{
		service: service,
		html:    html,
	}
}

type QuestionDataV2 struct {
	ID                 int
	CorrectChoices     string
	RuleQuestionNumber string
	Text               string
	Choices            []ChoiceDateV2
	QuestionNumber     int
	RuleName           string
}

type ChoiceDateV2 struct {
	ID     int
	Option string
	Text   string
}

func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(c.html, "base.tmpl", "home.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	allQuestions, err := c.service.GetAllQuestions(r.Context())
	var QuestionDataList []QuestionDataV2
	for i, question := range allQuestions {
		var choices []ChoiceDateV2
		var correctChoices []string
		for _, choice := range question.Choices {
			if choice.IsAnswer {
				correctChoices = append(correctChoices, choice.Option)
			}
			choices = append(choices, ChoiceDateV2{
				ID:     choice.ID,
				Option: choice.Option,
				Text:   choice.Text,
			})
		}
		QuestionDataList = append(QuestionDataList, QuestionDataV2{
			ID:                 i + 1,
			CorrectChoices:     strings.Join(correctChoices, ","),
			RuleQuestionNumber: question.RuleQuestionNumber,
			Text:               question.Text,
			Choices:            choices,
			QuestionNumber:     question.QuestionNumber,
			RuleName:           question.Rule.Name,
		})
	}
	err = tmpl.Execute(w, QuestionDataList)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func (c *Controller) Feedback(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFS(c.html, "base.tmpl", "feedback/feedback.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func (c *Controller) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(c.html, "feedback/submitFeedback.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	err = r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %s", err)
	}
	err = c.service.SubmitFeedback(r.Context(), Feedback{
		Name:  r.Form.Get("name"),
		Email: r.Form.Get("email"),
		Topic: r.Form.Get("topic"),
		Text:  r.Form.Get("feedback"),
	})
	if err != nil {
		log.Printf("Error submitting feedback: %s", err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}

}

type QuestionListPageData struct {
	Questions     []QuestionData
	LoadMoreParam LoadMoreParam
}

type QuestionData struct {
	ID                 int
	RuleID             string
	RuleName           string
	RuleQuestionNumber string
	Text               string
}

type LoadMoreParam struct {
	Search             string
	LastRuleSortOrder  int
	LastQuestionNumber int
	LastIndex          int
	Limit              int
}

func (c *Controller) QuestionByID(w http.ResponseWriter, r *http.Request) {
	id := queryParamInt(r, "id", 0)
	tmpl, err := template.ParseFS(c.html, "base.tmpl", "questionByID.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	err = tmpl.Execute(w, id)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func (c *Controller) QuestionList(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(c.html, "questionList.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	search := strings.TrimSpace(queryParamString(r, "search", ""))
	lastRuleSortOrder := queryParamInt(r, "lastRuleSortOrder", 0)
	lastQuestionNumber := queryParamInt(r, "lastQuestionNumber", 0)
	limit := 10
	questions, err := c.service.ListQuestions(r.Context(), nil, search, lastRuleSortOrder, lastQuestionNumber, limit)
	if err != nil {
		log.Printf("Error getting questions: %s", err)
	}

	if len(questions) > 0 {
		lastRuleSortOrder = questions[len(questions)-1].Rule.SortOrder
		lastQuestionNumber = questions[len(questions)-1].QuestionNumber
	}
	var questionsData []QuestionData
	for _, question := range questions {
		questionsData = append(questionsData, QuestionData{
			ID:                 question.ID,
			RuleID:             question.Rule.ID,
			RuleName:           question.Rule.Name,
			RuleQuestionNumber: question.RuleQuestionNumber,
			Text:               question.Text,
		})
	}
	data := QuestionListPageData{
		Questions: questionsData,
		LoadMoreParam: LoadMoreParam{
			LastRuleSortOrder:  lastRuleSortOrder,
			LastQuestionNumber: lastQuestionNumber,
			Limit:              10,
		},
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func (c *Controller) RandomQuestion(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFS(c.html, "base.tmpl", "questionByID.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func (c *Controller) NewQuestion(w http.ResponseWriter, r *http.Request) {
	id := queryParamInt(r, "id", 0)
	var question *Question
	var err error
	if id == 0 {
		rules, err := getQueryStrings(r, "rules")
		if err != nil {
			log.Printf("Error getting query strings: %s", err)
		}
		question, err = c.service.GetRandomQuestion(r.Context(), rules)
		if err != nil {
			log.Printf("Error getting random question: %s", err)
		}
	} else {
		question, err = c.service.GetQuestionByID(r.Context(), id)
		if err != nil {
			log.Printf("Error getting question: %s", err)
		}
	}
	tmpl, err := template.ParseFS(c.html, "question.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	err = tmpl.Execute(w, question)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func (c *Controller) Result(w http.ResponseWriter, r *http.Request) {
	questionID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/submit/"))
	if err != nil {
		log.Printf("Error parsing question ID: %s", err)
	}
	err = r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %s", err)
	}
	selected := r.Form["choices"]
	choices, err := c.service.GetChoicesByQuestionID(r.Context(), questionID)
	for i, choice := range choices {
		if slices.Contains(selected, choice.Option) {
			choices[i].IsSelected = true
		}
	}
	if err != nil {
		log.Printf("Error getting choices: %s", err)
	}

	tmpl, err := template.ParseFS(c.html, "result.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	err = tmpl.Execute(w, choices)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func (c *Controller) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// getQueryStrings parses a list of query strings from the request.
func getQueryStrings(r *http.Request, query string) ([]string, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	var ss []string
	for _, s := range r.Form[query] {
		ss = append(ss, strings.Split(s, ",")...)
	}
	return ss, nil
}

func queryParamString(r *http.Request, query string, defaultValue string) string {
	if s := r.URL.Query().Get(query); s != "" {
		return s
	}
	return defaultValue
}

func queryParamInt(r *http.Request, query string, defaultValue int) int {
	s := r.URL.Query().Get(query)
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Print("Error parsing int: ", err)
		return defaultValue
	}
	return i
}
