package token

import (
	"fmt"
)

type TokenType int

const (
	PAGE_NUMBER TokenType = iota
	RULE_NUMBER
	QUESTION_START
	CHOICE_START
	FREE_TEXT
)

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

func (t Tokenizer) Tokenize(s string) (*Token, error) {
	for _, matcher := range t.Matchers {
		if matcher.Match(s) {
			return &Token{matcher.Type, s}, nil
		}
	}
	return nil, fmt.Errorf("no token found for %s", s)
}
