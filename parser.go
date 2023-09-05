package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"log/slog"

	"github.com/aattwwss/ihf-referee-rules/token"
)

type Question struct {
	ID         int64
	QuestionID string
	Text       string
	Choices    []Choice
	PageNum    int
	RuleNum    int
}

type Choice struct {
	QuestionID string
	Text       string
	IsAnswer   bool
}

func main() {
	// Define a flag for the file path, with a default value of the current directory.
	questionPath := flag.String("q", "./questions.txt", "Path to the file to read")
	answerPath := flag.String("a", "./answers.txt", "Path to the file to read")
	flag.Parse()

	err := isValidFile(*questionPath)
	if err != nil {
		slog.Error("file path is invalid", slog.String("error", err.Error()))
		return
	}

	err = isValidFile(*answerPath)
	if err != nil {
		slog.Error("file path is invalid", slog.String("error", err.Error()))
		return
	}

	file, _ := os.Open(*questionPath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	tokenizer := token.NewTokenizer()
	tokens := []token.Token{}
	for scanner.Scan() {
		s := scanner.Text()
		token, err := tokenizer.Tokenize(s)
		if err != nil {
			slog.Error("tokenise error", slog.String("error", err.Error()))
			return
		}
		tokens = append(tokens, *token)
	}
	for _, token := range tokens {
		fmt.Printf("%-18s %s\n", token.Type.String(), token.Value)
	}
}

func isValidFile(filePath string) error {
	// Check if the provided path exists.
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// Check if the provided path is a directory or a file.
	if fileInfo.IsDir() {
		return err
	}
	return nil
}
