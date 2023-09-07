package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aattwwss/ihf-referee-rules/parser"
	"github.com/aattwwss/ihf-referee-rules/token"
	"golang.org/x/exp/slog"
)

const (
	QUESTION_DOC_TYPE = "QUESTION"
	ANSWER_DOC_TYPE   = "ANSWER"
)

func main() {
	// Define a flag for the file path, with a default value of the current directory.
	filePath := flag.String("file", "./questions.txt", "Path to the file to read")
	docType := flag.String("docType", "QUESTION", "Type of document to parse")
	flag.Parse()
	err := validateFlags(*filePath, *docType)
	if err != nil {
		slog.Error("flag error", slog.String("error", err.Error()))
		return
	}

	err = isValidFile(*filePath)
	if err != nil {
		slog.Error("file path is invalid", slog.String("error", err.Error()))
		return
	}

	file, _ := os.Open(*filePath)
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

	switch strings.ToUpper(*docType) {
	case QUESTION_DOC_TYPE:
		allQuestions := parser.ParseQuestion(tokens)
		b, _ := json.Marshal(allQuestions)
		fmt.Println(string(b))
	case ANSWER_DOC_TYPE:
		slog.Info("answers parser not implemented yet")
	default:
	}
}

func validateFlags(filepath string, docType string) error {
	err := isValidFile(filepath)
	if err != nil {
		return err
	}
	err = isValidDocType(docType)
	if err != nil {
		return err
	}
	return nil
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

func isValidDocType(docType string) error {
	s := strings.ToUpper(docType)
	if s != QUESTION_DOC_TYPE && s != ANSWER_DOC_TYPE {
		return fmt.Errorf("unrecognised document type %s", docType)
	}
	return nil
}
