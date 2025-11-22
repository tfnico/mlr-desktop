package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mattn/go-shellwords"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SaveOutput opens a save file dialog and saves the content to the selected file
func (a *App) SaveOutput(content string) error {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "Save Output",
		DefaultFilename: "output.txt",
	})
	if err != nil {
		return err
	}
	if path == "" {
		return nil // User cancelled
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// SelectInputFile opens a file dialog to select an input file
func (a *App) SelectInputFile() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Input File",
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// VerbConfig holds the configuration for a single verb
type VerbConfig struct {
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

// Config holds the application state
type Config struct {
	InputPath      string       `json:"inputPath"`
	InputMode      string       `json:"inputMode"`
	InputFormat    string       `json:"inputFormat"`
	Ragged         bool         `json:"ragged"`
	Headerless     bool         `json:"headerless"`
	FieldSeparator string       `json:"fieldSeparator"`
	OutputFormat   string       `json:"outputFormat"`
	Verbs          []VerbConfig `json:"verbs"`
	Options        string       `json:"options"`
}

// getLastStatePath returns the path to the last state file
func getLastStatePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "mlr_desktop_state.json"
	}
	return filepath.Join(home, ".mlr_desktop_state.json")
}

// SaveLastState saves the current configuration to the auto-save file
func (a *App) SaveLastState(config Config) error {
	return a.SaveConfig(config, getLastStatePath())
}

// LoadLastState loads the configuration from the auto-save file
func (a *App) LoadLastState() (Config, error) {
	return a.LoadConfig(getLastStatePath())
}

// SaveConfig saves the current configuration to a file
func (a *App) SaveConfig(config Config, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadConfig loads the configuration from a file
func (a *App) LoadConfig(path string) (Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}




// constructArgs helper to build the argument list
func (a *App) constructArgs(verbs []VerbConfig, options string, inputFormat string, ragged bool, headerless bool, fieldSeparator string, outputFormat string) ([]string, error) {
	var finalArgs []string

	// Add input format first
	if inputFormat != "" {
		finalArgs = append(finalArgs, inputFormat)
	}

	// CSV/TSV specific options
	if ragged {
		finalArgs = append(finalArgs, "--ragged")
	}
	if headerless {
		finalArgs = append(finalArgs, "--headerless-csv-input")
	}
	if fieldSeparator != "" && fieldSeparator != "," {
		finalArgs = append(finalArgs, "--ifs")
		finalArgs = append(finalArgs, fieldSeparator)
	}

	// Add output format
	if outputFormat != "" {
		finalArgs = append(finalArgs, outputFormat)
	}

	// Parse options (global flags like --icsv, --opprint)
	if options != "" {
		tokens, err := shellwords.Parse(options)
		if err != nil {
			return nil, fmt.Errorf("error parsing options: %v", err)
		}
		finalArgs = append(finalArgs, tokens...)
	}

	first := true
	for _, verb := range verbs {
		if !verb.Enabled {
			continue
		}
		// Tokenize the verb string (e.g. "head -n 100" -> ["head", "-n", "100"])
		tokens, err := shellwords.Parse(verb.Value)
		if err != nil {
			return nil, fmt.Errorf("error parsing verb '%s': %v", verb.Value, err)
		}
		
		if !first {
			finalArgs = append(finalArgs, "then")
		}
		finalArgs = append(finalArgs, tokens...)
		first = false
	}
	return finalArgs, nil
}

// GetCommand returns the constructed mlr command string
func (a *App) GetCommand(verbs []VerbConfig, options string, inputFormat string, ragged bool, headerless bool, fieldSeparator string, outputFormat string, inputMode string, inputPath string) (string, error) {
	args, err := a.constructArgs(verbs, options, inputFormat, ragged, headerless, fieldSeparator, outputFormat)
	if err != nil {
		return "", err
	}

	// Quote arguments for display if they contain spaces or special characters
	// We can use shellwords.Join but it might not be perfect for display. 
	// Let's do a simple quoting for now or just join with spaces if simple.
	// shellwords doesn't have a Join.
	// Let's iterate and quote if needed.
	var displayArgs []string
	for _, arg := range args {
		if strings.Contains(arg, " ") || strings.Contains(arg, ";") || strings.Contains(arg, "\"") || strings.Contains(arg, "'") {
			displayArgs = append(displayArgs, fmt.Sprintf("%q", arg))
		} else {
			displayArgs = append(displayArgs, arg)
		}
	}

	cmdStr := "mlr " + strings.Join(displayArgs, " ")

	if inputMode == "file" && inputPath != "" {
		// If file mode, append the file path
		if strings.Contains(inputPath, " ") {
			cmdStr += fmt.Sprintf(" %q", inputPath)
		} else {
			cmdStr += " " + inputPath
		}
	} else {
		// If text mode, maybe indicate input comes from stdin?
		// Or just leave it as is, implying pipe.
		// Let's prepend "cat input | " for clarity? Or just leave it.
		// The user request is "what the resulting final mlr command looks like".
		// Usually one runs `mlr [args] [files...]`.
		// If no files, it reads from stdin.
	}

	return cmdStr, nil
}

// Preview executes the mlr command with the given arguments and input
func (a *App) Preview(input string, verbs []VerbConfig, options string, inputFormat string, ragged bool, headerless bool, fieldSeparator string, outputFormat string) (string, error) {
	finalArgs, err := a.constructArgs(verbs, options, inputFormat, ragged, headerless, fieldSeparator, outputFormat)
	if err != nil {
		return "", err
	}

	// We use "mlr" from the path. Ensure mlr is installed.
	cmd := exec.Command("mlr", finalArgs...)

	// Set up stdin
	cmd.Stdin = strings.NewReader(input)

	// Capture stdout and stderr
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run the command
	err = cmd.Run()
	if err != nil {
		return stderr.String(), fmt.Errorf("mlr error: %s", stderr.String())
	}

	return out.String(), nil
}

// ReadFileHead reads the first n lines of a file
func (a *App) ReadFileHead(path string, n int) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for i := 0; i < n && scanner.Scan(); i++ {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}

