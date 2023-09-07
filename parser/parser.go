package parser

import (
	"fmt"
	"strconv"
	"strings"

	"log/slog"

	"github.com/aattwwss/ihf-referee-rules/token"
)

type Question struct {
	ID          int
	Text        string
	Choices     []Choice
	Rule        string
	QuestionNum int
}

type Choice struct {
	ID         int
	QuestionID int
	Option     string
	Text       string
	IsAnswer   bool
}

func ParserQuestion(tokens []token.Token) []Question {
	allQuestions := []Question{}
	groups := groupByQuestions(tokens)
	for _, group := range groups {
		q, err := toQuestion(len(allQuestions), group)
		if err != nil {
			slog.Error("convert to question error", slog.String("error", err.Error()))
			return nil
		}
		allQuestions = append(allQuestions, *q)
	}
	return allQuestions
}

// given the raw question string, split into the rule,
// question number and the question text
func splitQuestion(s string) (string, int, string) {
	bracketIndex := 1
	for i, c := range s {
		if c == ')' {
			bracketIndex = i
		}
	}
	var rule string
	var qString string
	var text string

	text = s[bracketIndex+1:]
	if strings.HasPrefix(s, "SAR") {
		rule = "SAR"
		qString = s[3:4]
	} else {
		s = s[0:bracketIndex]
		arr := strings.Split(s, ".")
		rule = arr[0]
		qString = arr[1]
	}
	n, _ := strconv.Atoi(qString)
	return rule, n, strings.TrimSpace(text)
}

// given the raw choice string, split into the option and text
func splitChoice(s string) (string, string) {
	arr := strings.SplitN(s, ")", 2)
	return strings.TrimSpace(arr[0]), strings.TrimSpace(arr[1])
}

// filter away unneeded tokens, keeping only the questions and their choices.
// put each question in their own group
func groupByQuestions(tokens []token.Token) [][]token.Token {
	var groups [][]token.Token
	var group []token.Token
	isStart := true
	for _, t := range tokens {
		if t.Type == token.PAGE_NUMBER || t.Type == token.RULE_NUMBER {
			continue
		}

		if t.Type == token.QUESTION_START {
			if !isStart {
				groups = append(groups, mergeFreeText(group))
				group = []token.Token{}
			} else {
				isStart = false
			}
		}

		if !isStart {
			group = append(group, t)
		}
	}
	// add in the remaining group after going through all the tokens
	groups = append(groups, mergeFreeText(group))
	return groups
}

// merge the free text token with their main token to form a single value
func mergeFreeText(tokens []token.Token) []token.Token {
	var merged []token.Token
	for _, t := range tokens {
		if t.Type == token.QUESTION_START || t.Type == token.CHOICE_START {
			merged = append(merged, t)
		}

		if t.Type == token.FREE_TEXT {
			last := merged[len(merged)-1]
			last.Value = fmt.Sprintf("%s %s", last.Value, t.Value)
			merged[len(merged)-1] = last
		}
	}
	return merged
}

// given a token group of question and choices, construct the question object
func toQuestion(id int, tokens []token.Token) (*Question, error) {
	var q Question
	var choices []Choice
	for _, t := range tokens {
		if t.Type == token.QUESTION_START {
			rule, qNum, text := splitQuestion(t.Value)
			q.ID = id
			q.Rule = rule
			q.QuestionNum = qNum
			q.Text = text
		} else if t.Type == token.CHOICE_START {
			option, choiceText := splitChoice(t.Value)
			c := Choice{
				QuestionID: id,
				Option:     option,
				Text:       choiceText,
			}
			choices = append(choices, c)
		} else {
			return nil, fmt.Errorf("token not recognised: %s", t.Type)
		}
	}
	q.Choices = choices
	return &q, nil
}
