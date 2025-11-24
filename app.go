package main

import (
	"bufio"
	"bytes"
	"container/list"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/climain"
	"github.com/johnkerl/miller/v6/pkg/input"
	"github.com/johnkerl/miller/v6/pkg/output"
	"github.com/johnkerl/miller/v6/pkg/transformers"
	"github.com/johnkerl/miller/v6/pkg/types"
	"github.com/mattn/go-shellwords"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SaveOutput opens a save file dialog and saves the content to the selected file
func (a *App) SaveOutput(content string) error {
	defer RecoverFromPanic("SaveOutput")
	
	LogInfo("SaveOutput called", nil)
	
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "Save Output",
		DefaultFilename: "output.txt",
	})
	if err != nil {
		LogError(err, "Failed to open save dialog", nil)
		return err
	}
	if path == "" {
		LogInfo("User cancelled save operation", nil)
		return nil // User cancelled
	}
	
	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		LogError(err, "Failed to write output file", logrus.Fields{"path": path})
		return err
	}
	
	LogInfo("Output saved successfully", logrus.Fields{"path": path, "size_bytes": len(content)})
	return nil
}

// SelectInputFile opens a file dialog to select an input file
func (a *App) SelectInputFile() (string, error) {
	defer RecoverFromPanic("SelectInputFile")
	
	LogInfo("SelectInputFile called", nil)
	
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Input File",
	})
	if err != nil {
		LogError(err, "Failed to open file dialog", nil)
		return "", err
	}
	
	if path != "" {
		LogInfo("File selected", logrus.Fields{"path": path})
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
	defer RecoverFromPanic("startup")
	
	a.ctx = ctx
	LogInfo("App startup completed", nil)
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
	defer RecoverFromPanic("SaveConfig")
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		LogError(err, "Failed to marshal config", logrus.Fields{"path": path})
		return err
	}
	
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		LogError(err, "Failed to write config file", logrus.Fields{"path": path})
		return err
	}
	
	LogInfo("Config saved", logrus.Fields{"path": path})
	return nil
}

// LoadConfig loads the configuration from a file
func (a *App) LoadConfig(path string) (Config, error) {
	defer RecoverFromPanic("LoadConfig")
	
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		LogWarn("Failed to read config file", logrus.Fields{"path": path, "error": err.Error()})
		return config, err
	}
	
	err = json.Unmarshal(data, &config)
	if err != nil {
		LogError(err, "Failed to unmarshal config", logrus.Fields{"path": path})
		return config, err
	}
	
	LogInfo("Config loaded", logrus.Fields{"path": path})
	return config, nil
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
			LogError(err, "Failed to parse options", logrus.Fields{"options": options})
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
			LogError(err, "Failed to parse verb", logrus.Fields{"verb": verb.Value})
			return nil, fmt.Errorf("error parsing verb '%s': %v", verb.Value, err)
		}
		
		if !first {
			finalArgs = append(finalArgs, "then")
		}
		finalArgs = append(finalArgs, tokens...)
		first = false
	}
	
	// Log the constructed arguments for debugging
	LogInfo("Constructed Miller arguments", logrus.Fields{
		"args": finalArgs,
		"args_string": strings.Join(finalArgs, " "),
	})
	
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

// Preview executes the mlr transformation using the Miller library directly
func (a *App) Preview(input string, verbs []VerbConfig, options string, inputFormat string, ragged bool, headerless bool, fieldSeparator string, outputFormat string) (string, error) {
	defer RecoverFromPanic("Preview")
	
	LogInfo("Preview transformation started", logrus.Fields{
		"input_format": inputFormat,
		"output_format": outputFormat,
		"verbs_count": len(verbs),
		"input_size": len(input),
	})
	
	// Build the command-line arguments as we would pass to mlr
	args, err := a.constructArgs(verbs, options, inputFormat, ragged, headerless, fieldSeparator, outputFormat)
	if err != nil {
		LogError(err, "Failed to construct args", nil)
		return "", err
	}

	// Miller's ParseCommandLine expects args[0] to be the program name (like os.Args)
	// Prepend "mlr" to match the expected format
	argsWithProgramName := append([]string{"mlr"}, args...)
	
	LogInfo("Calling ParseCommandLine", logrus.Fields{
		"full_args": argsWithProgramName,
	})

	// Parse the command line to get options and transformers
	mlrOptions, recordTransformers, err := climain.ParseCommandLine(argsWithProgramName)
	if err != nil {
		LogError(err, "Failed to parse command", logrus.Fields{"args": argsWithProgramName})
		return "", fmt.Errorf("error parsing command: %v", err)
	}

	// Create a temporary file for the input data since Miller's input readers expect file names
	tmpFile, err := os.CreateTemp("", "mlr-input-*.txt")
	if err != nil {
		LogError(err, "Failed to create temp file", nil)
		return "", fmt.Errorf("error creating temp file: %v", err)
	}
	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)

	// Write input data to temp file
	if _, err := tmpFile.WriteString(input); err != nil {
		tmpFile.Close()
		LogError(err, "Failed to write to temp file", logrus.Fields{"file": tmpFileName})
		return "", fmt.Errorf("error writing to temp file: %v", err)
	}
	tmpFile.Close()

	// Set up output buffer
	var outputBuffer bytes.Buffer
	bufferedOutputStream := bufio.NewWriter(&outputBuffer)

	// Run the Miller transformation
	err = runMillerTransformation([]string{tmpFileName}, mlrOptions, recordTransformers, bufferedOutputStream)
	if err != nil {
		LogError(err, "Miller transformation failed", nil)
		return "", err
	}

	bufferedOutputStream.Flush()
	result := outputBuffer.String()
	
	LogInfo("Preview transformation completed", logrus.Fields{
		"output_size": len(result),
	})
	
	return result, nil
}

// runMillerTransformation runs the Miller transformation pipeline
// This is based on Miller's streaming architecture from the library examples
func runMillerTransformation(
	fileNames []string,
	options *cli.TOptions,
	recordTransformers []transformers.IRecordTransformer,
	outputStream io.Writer,
) error {
	defer RecoverFromPanic("runMillerTransformation")
	
	outputIsStdout := false

	// Create initial context
	initialContext := types.NewContext()

	// Create the record reader
	recordReader, err := input.Create(&options.ReaderOptions, options.ReaderOptions.RecordsPerBatch)
	if err != nil {
		LogError(err, "Failed to create record reader", nil)
		return fmt.Errorf("error creating record reader: %v", err)
	}

	// Create the record writer
	recordWriter, err := output.Create(&options.WriterOptions)
	if err != nil {
		LogError(err, "Failed to create record writer", nil)
		return fmt.Errorf("error creating record writer: %v", err)
	}

	// Set up channels for the pipeline
	readerChannel := make(chan *list.List, 2)              // reader -> transformer
	writerChannel := make(chan *list.List, 1)              // transformer -> writer
	inputErrorChannel := make(chan error, 1)               // reader errors
	doneWritingChannel := make(chan bool, 1)               // writer done signal
	dataProcessingErrorChannel := make(chan bool, 1)       // data processing errors
	readerDownstreamDoneChannel := make(chan bool, 1)      // downstream done signal

	bufferedOutputStream := bufio.NewWriter(outputStream)

	// Start the pipeline goroutines
	go recordReader.Read(fileNames, *initialContext, readerChannel, inputErrorChannel, readerDownstreamDoneChannel)
	go transformers.ChainTransformer(readerChannel, readerDownstreamDoneChannel, recordTransformers, writerChannel, options)
	go output.ChannelWriter(writerChannel, recordWriter, &options.WriterOptions, doneWritingChannel, dataProcessingErrorChannel, bufferedOutputStream, outputIsStdout)

	// Wait for completion or error
	var retval error
	done := false
	for !done {
		select {
		case ierr := <-inputErrorChannel:
			LogError(ierr, "Input error during Miller transformation", nil)
			retval = ierr
			done = true
		case <-dataProcessingErrorChannel:
			retval = errors.New("data processing error")
			LogError(retval, "Data processing error during Miller transformation", nil)
			done = true
		case <-doneWritingChannel:
			done = true
		}
	}

	bufferedOutputStream.Flush()
	if retval != nil {
		LogError(retval, "Miller transformation completed with errors", nil)
	}
	return retval
}

// ReadFileHead reads the first n lines of a file
func (a *App) ReadFileHead(path string, n int) (string, error) {
	defer RecoverFromPanic("ReadFileHead")
	
	LogInfo("Reading file head", logrus.Fields{"path": path, "lines": n})
	
	file, err := os.Open(path)
	if err != nil {
		LogError(err, "Failed to open file", logrus.Fields{"path": path})
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for i := 0; i < n && scanner.Scan(); i++ {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		LogError(err, "Scanner error while reading file", logrus.Fields{"path": path})
		return "", err
	}
	
	LogInfo("File head read successfully", logrus.Fields{"path": path, "lines_read": len(lines)})
	return strings.Join(lines, "\n"), nil
}

