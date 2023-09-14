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
	questionPath := flag.String("q", "./questions.txt", "Path to the questions text file")
	answerPath := flag.String("a", "./answers.txt", "Path to the answers text file")
	outputPath := flag.String("0", "./questions.json", "Path to the output file")
	flag.Parse()
	err := isValidFile(*questionPath)
	if err != nil {
		slog.Error("question path is invalid", slog.String("error", err.Error()))
		return
	}
	err = isValidFile(*answerPath)
	if err != nil {
		slog.Error("answer path is invalid", slog.String("error", err.Error()))
		return
	}

	aFile, _ := os.Open(*answerPath)
	defer aFile.Close()
	answerMap := parser.ParseAnswer(aFile)
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

	allQuestions := parser.ParseQuestion(tokens, answerMap)
	b, _ := json.MarshalIndent(allQuestions, "", "    ")
	fmt.Println(string(b))

	outputFile, err := os.Create(*outputPath)
	if err != nil {
		slog.Error("error creating file: ", slog.String("error", err.Error()))
		return
	}
	defer outputFile.Close()

	// Write the JSON data to the file
	_, err = outputFile.Write(b)
	if err != nil {
		slog.Error("error writing to file: ", slog.String("error", err.Error()))
		return
	}
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
