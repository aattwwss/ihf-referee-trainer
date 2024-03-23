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
	GetRandomQuestion(ctx context.Context) (*Question, error)
	GetChoicesByQuestionID(ctx context.Context, questionID int) ([]Choice, error)
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
	question, err := c.service.GetRandomQuestion(r.Context())
	if err != nil {
		log.Printf("Error getting random question: %s", err)
	}
	tmpl, err := template.ParseFS(c.html, "base.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
	}
	err = tmpl.Execute(w, question)
	if err != nil {
		log.Printf("Error executing template: %s", err)
	}
}

func (c *Controller) NewQuestion(w http.ResponseWriter, r *http.Request) {
	question, err := c.service.GetRandomQuestion(r.Context())
	if err != nil {
		log.Printf("Error getting random question: %s", err)
	}
	tmpl, err := template.ParseFS(c.html, "game.tmpl")
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
