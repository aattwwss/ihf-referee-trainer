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
	GetQuestionByID(ctx context.Context, id int) (*Question, error)
	GetRandomQuestion(ctx context.Context, rules []string) (*Question, error)
	GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error)
	ListQuestions(ctx context.Context, rules []string, search string, lastRuleSortOrder int, lastQuestionNumber int, limit int) ([]Question, error)
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

func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(c.html, "base.tmpl", "home.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
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
	questions, err := c.service.ListQuestions(r.Context(), nil, search, lastRuleSortOrder, lastQuestionNumber, 10)
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
