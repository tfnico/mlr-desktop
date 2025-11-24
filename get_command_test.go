package main

import (
	"testing"
)

// TestGetCommandQuoting tests that the GetCommand function properly quotes arguments
func TestGetCommandQuoting(t *testing.T) {
	app := NewApp()
	
	tests := []struct {
		name           string
		verbs          []VerbConfig
		expectedSubstr string // substring we expect to find in the command
		description    string
	}{
		{
			name: "Miller field reference with $",
			verbs: []VerbConfig{
				{Value: `put $total = $price * $quantity`, Enabled: true},
			},
			expectedSubstr: `'$total'`,
			description:    "Should use single quotes for tokens with $ to prevent shell variable expansion",
		},
		{
			name: "Filter with double quotes",
			verbs: []VerbConfig{
				{Value: `filter $1 == "ART"`, Enabled: true},
			},
			expectedSubstr: `filter '$1'`,
			description:    "Should use single quotes for $field references",
		},
		{
			name: "Put with double quotes in value",
			verbs: []VerbConfig{
				{Value: `put $category = "Electronics"`, Enabled: true},
			},
			expectedSubstr: `'$category'`,
			description:    "Should use single quotes for $field and quoted string values",
		},
		{
			name: "Simple verb without special chars",
			verbs: []VerbConfig{
				{Value: `head -n 10`, Enabled: true},
			},
			expectedSubstr: `head -n 10`,
			description:    "Should not quote simple arguments without special characters",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := app.GetCommand(tt.verbs, "", "--icsv", false, false, ",", "--opprint", "text", "")
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if !contains(cmd, tt.expectedSubstr) {
				t.Errorf("%s\nExpected command to contain: %s\nGot: %s", 
					tt.description, tt.expectedSubstr, cmd)
			}
			
			// Additional check: ensure no unquoted $ signs in the output
			// (they should all be within single quotes)
			if hasUnquotedDollar(cmd) {
				t.Errorf("Command contains unquoted $ which could be evaluated by shell: %s", cmd)
			}
		})
	}
}

// hasUnquotedDollar checks if there are any $ signs outside of quotes
func hasUnquotedDollar(s string) bool {
	inSingleQuote := false
	inDoubleQuote := false
	
	for i := 0; i < len(s); i++ {
		ch := s[i]
		
		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
		} else if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		} else if ch == '$' && !inSingleQuote && !inDoubleQuote {
			return true
		}
	}
	return false
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
