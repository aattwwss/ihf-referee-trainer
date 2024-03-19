package pdf

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// run the pdftotext command
func PdfToText(file io.Reader) (string, error) {
	params := []string{
		"-layout",
		"-nopgbrk",
		"-",
		"-",
	}

	cmd := exec.Command("pdftotext", params...)
	b, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	cmd.Stdin = bytes.NewReader(b)

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error executing pdftotext binary: %w", err)
	}
    return out.String(), nil
}
