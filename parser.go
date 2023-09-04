package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"log/slog"
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

	qTokens = []TokenType{}
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
	questions := []Question{}
	for scanner.Scan() {
		s := scanner.Text()
		isPageNum, err := regexp.MatchString(strings.TrimSpace(s), "^[0-9]*$")
		if err != nil {
			slog.Error("parse line error", slog.String("error", err.Error()))
			return
		}
		if isPageNum {
			pageNumber, _ := strconv.Atoi(strings.TrimSpace(s))
		}
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

func tokenise(line string) (*TokenType, error) {
	isPageNum, err := regexp.MatchString(strings.TrimSpace(line), "^[0-9]*$")
	if err != nil {
		slog.Error("parse line error", slog.String("error", err.Error()))
		return nil, fmt.Errorf("tokenise: failed to parse page number: %w", err)
	}
}
