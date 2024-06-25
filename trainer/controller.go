package trainer

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	GetAllQuestions(ctx context.Context) ([]Question, error)
	GetAllRules(ctx context.Context) ([]Rule, error)
	GetQuestionByID(ctx context.Context, id int) (*Question, error)
	GetRandomQuestion(ctx context.Context, rules []string) (*Question, error)
	GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error)
	SubmitFeedback(ctx context.Context, feedback Feedback) error
	SubmitQuizConfig(ctx context.Context, quizConfig QuizConfig) (string, error)
	GetQuestionsByQuizConfigKey(ctx context.Context, key string) ([]Question, error)
	EvaluateQuizAnswer(ctx context.Context, quizConfigKey string, choiceIDs []int) ([]QuestionResult, error)
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
	Choices            []ChoiceDataV2
	QuestionNumber     int
	RuleName           string
}

type ChoiceDataV2 struct {
	ID         int
	Option     string
	Text       string
	Result     string
	IsSelected bool
}

func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(c.html, "base.tmpl", "home.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	allQuestions, err := c.service.GetAllQuestions(r.Context())
	var QuestionDataList []QuestionDataV2
	for i, question := range allQuestions {
		var choices []ChoiceDataV2
		var correctChoices []string
		for _, choice := range question.Choices {
			if choice.IsAnswer {
				correctChoices = append(correctChoices, choice.Option)
			}
			choices = append(choices, ChoiceDataV2{
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
		Topic: r.Form.Get("feedback-category"),
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

type QuizConfigData struct {
	MaxNumQuestions int
	Seed            string
	RulesFilter     []QuizConfigRule
}

type QuizConfigRule struct {
	ID   string
	Name string
}

func (c *Controller) QuizConfig(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(c.html, "base.tmpl", "quiz/quizConfig.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	rules, err := c.service.GetAllRules(r.Context())
	questions, err := c.service.GetAllQuestions(r.Context())
	var rulesFilter []QuizConfigRule
	for _, rule := range rules {
		rulesFilter = append(rulesFilter, QuizConfigRule{
			ID:   rule.ID,
			Name: rule.Name,
		})
	}
	quizConfigData := QuizConfigData{
		MaxNumQuestions: len(questions),
		Seed:            fmt.Sprintf("%d", time.Now().UnixMilli()),
		RulesFilter:     rulesFilter,
	}
	err = tmpl.Execute(w, quizConfigData)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func (c *Controller) SubmitQuizConfig(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing query: %s", err)
	}
	seed, _ := strconv.Atoi(r.FormValue("seed"))
	ruleIDs := r.Form["rules-filter"]
	numQuestions, _ := strconv.Atoi(r.FormValue("num-questions"))
	duration, _ := strconv.Atoi(r.FormValue("duration"))
	hasNegativeMarking, _ := strconv.ParseBool(r.FormValue("negative-marking"))
	quizConfig := QuizConfig{
		Seed:               seed,
		RuleIDs:            ruleIDs,
		NumQuestions:       numQuestions,
		DurationInMinutes:  duration,
		HasNegativeMarking: hasNegativeMarking,
	}
	key, err := c.service.SubmitQuizConfig(r.Context(), quizConfig)
	if err != nil {
		log.Printf("Error submitting quiz config: %s", err)
	}
	http.Redirect(w, r, fmt.Sprintf("/quiz/%s", key), http.StatusFound)
}

func (c *Controller) DoQuiz(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		http.Redirect(w, r, "/quiz/config", http.StatusFound)
	}
	questions, err := c.service.GetQuestionsByQuizConfigKey(r.Context(), key)
	if err != nil {
		log.Printf("Error getting quiz questions: %s", err)
	}

	tmpl, err := template.ParseFS(c.html, "base.tmpl", "quiz/quiz.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	var QuestionDataList []QuestionDataV2
	for i, question := range questions {
		var choices []ChoiceDataV2
		for _, choice := range question.Choices {
			choices = append(choices, ChoiceDataV2{
				ID:     choice.ID,
				Option: choice.Option,
				Text:   choice.Text,
			})
		}
		QuestionDataList = append(QuestionDataList, QuestionDataV2{
			ID:                 i + 1,
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

type QuestionResultData struct {
	QuestionDataV2
	TotalScore int
	Score      int
}

type QuizResultData struct {
	QuestionResults []QuestionResultData
	TotalScore      int
	Score           int
	Percent         int
}

func (c *Controller) SubmitQuiz(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(c.html, "base.tmpl", "quiz/result.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	err = r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %s", err)
	}
	var choiceIDs []int
	for k, _ := range r.Form {
		suffix := strings.TrimPrefix(k, "choice")
		//suffix = strings.TrimSpace(suffix)
		choiceID, err := strconv.Atoi(suffix)
		if err != nil {
			log.Printf("Error parsing choice ID: %s", err)
		}
		choiceIDs = append(choiceIDs, choiceID)
	}
	questions, err := c.service.EvaluateQuizAnswer(r.Context(), r.PathValue("key"), choiceIDs)
	if err != nil {
		log.Printf("Error evaluating quiz answer: %s", err)
	}
	var QuestionDataList []QuestionResultData
	quizTotalScore := 0
	quizScore := 0
	for i, question := range questions {
		var choices []ChoiceDataV2
		for _, choice := range question.Choices {
			var result string
			if choice.Result != nil {
				result = choice.Result.String()
			}
			choices = append(choices, ChoiceDataV2{
				ID:         choice.ID,
				Option:     choice.Option,
				Text:       choice.Text,
				Result:     result,
				IsSelected: slices.Contains(choiceIDs, choice.ID),
			})
		}
		QuestionDataList = append(QuestionDataList, QuestionResultData{
			QuestionDataV2: QuestionDataV2{
				ID:                 i + 1,
				RuleQuestionNumber: question.RuleQuestionNumber,
				Text:               question.Text,
				Choices:            choices,
				QuestionNumber:     question.QuestionNumber,
				RuleName:           question.Rule.Name,
			},
			TotalScore: question.TotalScore,
			Score:      question.Score,
		})
		quizTotalScore += question.TotalScore
		quizScore += question.Score
	}
	quizResultData := QuizResultData{
		QuestionResults: QuestionDataList,
		TotalScore:      quizTotalScore,
		Score:           quizScore,
		Percent:         quizScore * 100 / quizTotalScore,
	}

	err = tmpl.Execute(w, quizResultData)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func (c *Controller) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
