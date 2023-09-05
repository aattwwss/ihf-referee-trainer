package token

import (
	"fmt"
	"regexp"
	"strings"
)

type TokenType int

const (
	PAGE_NUMBER TokenType = iota
	RULE_NUMBER
	QUESTION_START
	CHOICE_START
	FREE_TEXT
)

func (tt TokenType) String() string {
	return []string{"PAGE_NUMBER", "RULE_NUMBER", "QUESTION_START", "CHOICE_START", "FREE_TEXT"}[tt]

}

type Token struct {
	Type  TokenType
	Value string
}

type TokenPattern struct {
	Type  TokenType
	Match func(s string) bool
}

type Tokenizer struct {
	Matchers []TokenPattern
}

func NewTokenizer() Tokenizer {
	// if the line contains only the word Rule and a number
	ruleNumberMatcher := func(s string) bool {
		s = strings.TrimSpace(s)
		regex := regexp.MustCompile(`^Rule \d+$`)
		return regex.MatchString(s)
	}
	ruleNumberTokenPattern := TokenPattern{
		Type:  RULE_NUMBER,
		Match: ruleNumberMatcher,
	}

	// line contains the number only
	pageNumberMatcher := func(s string) bool {
		s = strings.TrimSpace(s)
		regex := regexp.MustCompile(`^\d+$`)
		return regex.MatchString(s)
	}
	pageNumberTokenPattern := TokenPattern{
		Type:  PAGE_NUMBER,
		Match: pageNumberMatcher,
	}

	// line starts with the rule number and question number follow by a close bracket
	// e.g. 12.34) some other text
	questionStartMatcher := func(s string) bool {
		s = strings.TrimSpace(s)
		regex := regexp.MustCompile(`^(\d+\.\d+\)|SAR\d\)) .+$`)
		return regex.MatchString(s)
	}
	questionStartTokenPattern := TokenPattern{
		Type:  QUESTION_START,
		Match: questionStartMatcher,
	}

	// line starts with the choice character follow by a close bracket
	// e.g. a) some other text
	choiceStartMatcher := func(s string) bool {
		s = strings.TrimSpace(s)
		regex := regexp.MustCompile(`^.\) .+$`)
		return regex.MatchString(s)
	}
	choiceStartTokenPattern := TokenPattern{
		Type:  CHOICE_START,
		Match: choiceStartMatcher,
	}

	// free text is anything else
	freeTextMatcher := func(s string) bool {
		return true
	}
	freeTextTokenPattern := TokenPattern{
		Type:  FREE_TEXT,
		Match: freeTextMatcher,
	}

	// return the tokenizer where the order of the token pattern matters!
	return Tokenizer{
		Matchers: []TokenPattern{
			ruleNumberTokenPattern,
			pageNumberTokenPattern,
			questionStartTokenPattern,
			choiceStartTokenPattern,
			freeTextTokenPattern,
		},
	}
}

func (t Tokenizer) Tokenize(s string) (*Token, error) {
	for _, matcher := range t.Matchers {
		if matcher.Match(s) {
			return &Token{matcher.Type, s}, nil
		}
	}
	return nil, fmt.Errorf("no token found for %s", s)
}
