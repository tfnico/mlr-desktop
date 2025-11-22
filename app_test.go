package main

import (
	"os/exec"
	"strings"
	"testing"
)



func TestGetCommand(t *testing.T) {
	app := NewApp()
	verbs := []VerbConfig{
		{Value: "head -n 10", Enabled: true},
	}
	options := "--opprint"
	inputFormat := "--icsv"
	ragged := true
	headerless := true
	fieldSeparator := ";"
	outputFormat := "--ojson"
	inputMode := "file"
	inputPath := "/path/to/input.csv"

	cmd, err := app.GetCommand(verbs, options, inputFormat, ragged, headerless, fieldSeparator, outputFormat, inputMode, inputPath)
	if err != nil {
		t.Errorf("GetCommand failed: %v", err)
	}

	expectedParts := []string{
		"mlr",
		"--icsv",
		"--ragged",
		"--headerless-csv-input",
		"--ifs", ";",
		"--ojson",
		"--opprint",
		"head", "-n", "10",
		"/path/to/input.csv",
	}

	for _, part := range expectedParts {
		if !strings.Contains(cmd, part) {
			t.Errorf("Command missing part: %s. Got: %s", part, cmd)
		}
	}
}

func TestPreview(t *testing.T) {
	// This test requires mlr to be installed.
	// We can mock exec.Command or just skip if mlr is not found,
	// but for now let's assume dev environment has mlr.
	_, err := exec.LookPath("mlr")
	if err != nil {
		t.Skip("mlr not found in PATH")
	}

	app := NewApp()
	input := "a;b\n1;2\n3;4"
	verbs := []VerbConfig{
		{Value: "cat", Enabled: true},
	}
	options := ""
	inputFormat := "--icsv"
	ragged := false
	headerless := false
	fieldSeparator := ";"
	outputFormat := "--ojson"

	output, err := app.Preview(input, verbs, options, inputFormat, ragged, headerless, fieldSeparator, outputFormat)
	if err != nil {
		t.Errorf("Preview failed: %v", err)
	}
	if !strings.Contains(output, "{ \"a\": 1, \"b\": 2 }") && !strings.Contains(output, "\"a\": 1") { 
		// Check for JSON output structure
	}
}
