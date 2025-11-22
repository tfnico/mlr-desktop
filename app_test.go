package main

import (
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
	// Test now uses the Miller library directly - no external binary needed
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
	
	// Check for JSON output structure - should contain field "a" and "b"
	if !strings.Contains(output, "\"a\"") || !strings.Contains(output, "\"b\"") {
		t.Errorf("Expected JSON output with fields 'a' and 'b', got: %s", output)
	}
	
	// Should contain values 1, 2, 3, 4
	if !strings.Contains(output, "1") || !strings.Contains(output, "2") {
		t.Errorf("Expected output to contain values from input, got: %s", output)
	}
}
