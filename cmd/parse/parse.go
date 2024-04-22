package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/aattwwss/ihf-referee-rules/parser"
	"github.com/aattwwss/ihf-referee-rules/pdf"
	"github.com/aattwwss/ihf-referee-rules/token"
	"golang.org/x/exp/slog"
)

const (
	delimiter = '|'
)

func main() {
	// Define a flag for the file path, with a default value of the current directory.
	questionPath := flag.String("q", "./questions.pdf", "Path to the questions pdf file")
	answerPath := flag.String("a", "./answers.pdf", "Path to the answers pdf file")
	formatType := flag.String("f", "json", "format to output the parsed results")
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
	tokenizer := token.NewTokenizer()
	var tokens []token.Token
	qs, err := pdf.PdfToText(file)
	if err != nil {
		slog.Error("pdf to text error", slog.String("error", err.Error()))
		return
	}

	for _, s := range strings.Split(qs, "\n") {
		tokensFromLine, err := tokenizer.Tokenize(s)
		if err != nil {
			slog.Error("tokenise error", slog.String("error", err.Error()))
			return
		}
		tokens = append(tokens, *tokensFromLine)
	}

	allQuestions := parser.ParseQuestion(tokens, answerMap)
	handleOutput(allQuestions, *formatType)
}

func handleOutput(allQuestions []parser.Question, formatType string) {
	switch strings.ToLower(formatType) {
	case "sql":
		outputFile, err := os.Create("data.sql")
		if err != nil {
			slog.Error("error creating file: ", slog.String("error", err.Error()))
			return
		}
		defer outputFile.Close()

		var allChoices []parser.Choice
		var allReferences []parser.Reference
		var allRules []parser.Rule
		questionSQLText := "INSERT INTO question (id, text, rule, question_number) VALUES\n"

		for idx, q := range allQuestions {
			allChoices = append(allChoices, q.Choices...)
			allReferences = append(allReferences, q.References...)
			if !slices.ContainsFunc(allRules, func(r parser.Rule) bool {
				return r.Text == q.Rule.Text
			}) {
				allRules = append(allRules, q.Rule)
			}

			questionSQLText += fmt.Sprintf("(%d, '%s', '%s', %d)", q.ID, q.Text, q.Rule.Text, q.QuestionNumber)
			if (idx + 1) != len(allQuestions) {
				questionSQLText += ",\n"
			}
		}
		questionSQLText += ";\n\n"

		choiceSQLText := "INSERT INTO choice (question_id, option, text, is_answer) VALUES\n"
		for idx, c := range allChoices {
			choiceSQLText += fmt.Sprintf("(%d, '%s', '%s', %v)", c.QuestionID, c.Option, c.Text, c.IsAnswer)
			if (idx + 1) != len(allChoices) {
				choiceSQLText += ",\n"
			}
		}
		choiceSQLText += ";\n\n"

		referenceSQLText := "INSERT INTO reference (question_id, text) VALUES\n"
		for idx, ref := range allReferences {
			referenceSQLText += fmt.Sprintf("(%d, '%s')", ref.QuestionID, ref.Text)
			if (idx + 1) != len(allReferences) {
				referenceSQLText += ",\n"
			}
		}
		referenceSQLText += ";\n\n"

		ruleSQLText := "INSERT INTO rule (text, sort_order) VALUES\n"
		for idx, r := range allRules {
			r.SortOrder = idx
			ruleSQLText += fmt.Sprintf("('%s', %d)", r.Text, r.SortOrder)
			if (idx + 1) != len(allRules) {
				ruleSQLText += ",\n"
			}
		}
		ruleSQLText += ";\n\n"
		outputFile.WriteString(ruleSQLText + questionSQLText + choiceSQLText + referenceSQLText)

	case "json":

		b, _ := json.MarshalIndent(allQuestions, "", "  ")

		outputFile, err := os.Create("questions_answers.json")
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
