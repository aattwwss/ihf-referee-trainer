package parser

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"log/slog"

	"github.com/aattwwss/ihf-referee-rules/token"
	"github.com/aattwwss/ihf-referee-rules/pdf"
	"golang.org/x/exp/slices"
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

func ParseQuestion(tokens []token.Token, answerMap map[string]map[int][]string) []Question {
	allQuestions := []Question{}
	groups := groupByQuestions(tokens)
	for _, group := range groups {
		q, err := toQuestion(len(allQuestions)+1, group, answerMap)
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
	bracketIndex := strings.IndexRune(s, ')')
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
func toQuestion(id int, tokens []token.Token, answerMap map[string]map[int][]string) (*Question, error) {
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
				IsAnswer:   slices.Contains(answerMap[q.Rule][q.QuestionNum], option),
			}
			choices = append(choices, c)
		} else {
			return nil, fmt.Errorf("token not recognised: %s", t.Type)
		}
	}
	q.Choices = choices
	return &q, nil
}

// ParseAnswer returns the list of answers for a given rule and question number
func ParseAnswer(file io.Reader) map[string]map[int][]string {
	ansMap := map[string]map[int][]string{}
	s, _ := pdf.PdfToText(file)
	for _,s := range strings.Split(s, "\n") {
		s = strings.TrimSpace(s)
		if !hasAnswers(s) {
			continue
		}
		rule, questionNum, answers := splitAnswer(s)
		ruleMap, ok := ansMap[rule]
		if ok {
			ruleMap[questionNum] = answers
		} else {
			ansMap[rule] = map[int][]string{questionNum: answers}
		}
	}
	return ansMap
}

func hasAnswers(s string) bool {
	pattern := `\d\)`
	regex, _ := regexp.Compile(pattern)
	return regex.MatchString(s)
}

// given the raw question string, split into the rule,
// question number and the question text
func splitAnswer(s string) (string, int, []string) {
	bracketIndex := 1
	for i, c := range s {
		if c == ')' {
			bracketIndex = i
		}
	}
	var rule string
	var aString string
	var text string

	text = s[bracketIndex+1:]
	if strings.HasPrefix(s, "SAR") {
		rule = "SAR"
		aString = s[3:4]
	} else {
		s = s[0:bracketIndex]
		arr := strings.Split(s, ".")
		rule = arr[0]
		aString = arr[1]
	}
	n, _ := strconv.Atoi(aString)
	return rule, n, parseChoices(text)
}

// Given the line of text of the answer to a question
// return the correct answers
// We check if a character has a space before it, and a space or a comma after.
func parseChoices(s string) []string {
	s = " " + s
	var answers []string
	for i := 0; i < len(s)-3; i++ {
		if s[i] == ' ' && unicode.IsLower(rune(s[i+1])) && (s[i+2] == ' ' || s[i+2] == ',') {
			answers = append(answers, string(s[i+1:i+2]))
		}
	}
	return answers
}
