package main

import "testing"

func TestParseCommand(t *testing.T) {
	app := NewApp()
	
	tests := []struct {
		name           string
		command        string
		wantFormat     string
		wantOutputFmt  string
		wantVerbsCount int
		wantOptions    string
		wantError      bool
	}{
		{
			name:           "Simple head command",
			command:        "mlr --icsv --opprint head -n 5",
			wantFormat:     "--icsv",
			wantOutputFmt:  "--opprint",
			wantVerbsCount: 1,
			wantOptions:    "",
			wantError:      false,
		},
		{
			name:           "Multiple verbs with then",
			command:        "mlr --icsv head -n 10 then cut -f name,price",
			wantFormat:     "--icsv",
			wantOutputFmt:  "",
			wantVerbsCount: 2,
			wantOptions:    "",
			wantError:      false,
		},
		{
			name:           "Command without mlr prefix",
			command:        "--icsv --opprint head -n 5",
			wantFormat:     "--icsv",
			wantOutputFmt:  "--opprint",
			wantVerbsCount: 1,
			wantOptions:    "",
			wantError:      false,
		},
		{
			name:           "Command with additional options",
			command:        "mlr --skip-comments --icsv head -n 5",
			wantFormat:     "--icsv",
			wantOutputFmt:  "",
			wantVerbsCount: 1,
			wantOptions:    "--skip-comments",
			wantError:      false,
		},
		{
			name:           "Complex command with file path",
			command:        "mlr --itsv --ragged --headerless-csv-input --opprint --implicit-csv-header head -n 50 then put \"$row_length = length($0)\" then label category /home/thomas/Downloads/strauss/import_de.txt",
			wantFormat:     "--itsv",
			wantOutputFmt:  "--opprint",
			wantVerbsCount: 3,
			wantOptions:    "--implicit-csv-header",
			wantError:      false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := app.ParseCommand(tt.command)
			
			if tt.wantError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if !tt.wantError {
				if config.InputFormat != tt.wantFormat {
					t.Errorf("InputFormat = %v, want %v", config.InputFormat, tt.wantFormat)
				}
				if config.OutputFormat != tt.wantOutputFmt {
					t.Errorf("OutputFormat = %v, want %v", config.OutputFormat, tt.wantOutputFmt)
				}
				if len(config.Verbs) != tt.wantVerbsCount {
					t.Errorf("Verbs count = %v, want %v", len(config.Verbs), tt.wantVerbsCount)
				}
				if config.Options != tt.wantOptions {
					t.Errorf("Options = %v, want %v", config.Options, tt.wantOptions)
				}
			}
		})
	}
}
