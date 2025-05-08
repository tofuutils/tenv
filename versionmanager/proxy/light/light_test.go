package lightproxy

import (
	"os"
	"testing"

	"github.com/tofuutils/tenv/v4/config/cmdconst"
)

func TestExec(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
	}{
		{
			name:     "valid tool name",
			toolName: cmdconst.TofuName,
		},
		{
			name:     "invalid tool name",
			toolName: "nonexistent",
		},
		{
			name:     "empty tool name",
			toolName: "",
		},
	}

	// Save original args
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test args
			os.Args = []string{"test", "--version"}

			// Execute command
			Exec(tt.toolName)
			// Note: Since Exec calls os.Exit, we can't actually test the return value
			// This test is mainly to ensure the function doesn't panic
		})
	}
}
