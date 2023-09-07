package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/aattwwss/ihf-referee-rules/parser"
	"github.com/aattwwss/ihf-referee-rules/token"
	"golang.org/x/exp/slog"
)

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

	allQuestions := parser.ParseQuestion(tokens)
	b, _ := json.Marshal(allQuestions)
	fmt.Println(string(b))
}

// check if path exists and is a file, not a folder
func isValidFile(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("path is not a file")
	}
	return nil
}
