package main

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/config/cmdconst"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
	"github.com/tofuutils/tenv/v4/versionmanager"
	"github.com/tofuutils/tenv/v4/versionmanager/builder"
	lightproxy "github.com/tofuutils/tenv/v4/versionmanager/proxy/light"
)

type mockRetriever struct {
	installFunc      func(ctx context.Context, version string, targetPath string) error
	listVersionsFunc func(ctx context.Context) ([]string, error)
}

func (m *mockRetriever) Install(ctx context.Context, version string, targetPath string) error {
	return m.installFunc(ctx, version, targetPath)
}

func (m *mockRetriever) ListVersions(ctx context.Context) ([]string, error) {
	return m.listVersionsFunc(ctx)
}

func TestTfCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		execErr     error
		expectError bool
	}{
		{
			name:        "successful execution",
			args:        []string{"version"},
			execErr:     nil,
			expectError: false,
		},
		{
			name:        "execution error",
			args:        []string{"invalid"},
			execErr:     assert.AnError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test configuration
			logger := hclog.New(&hclog.LoggerOptions{
				Output: os.Stderr,
				Level:  hclog.Info,
			})
			displayer := loghelper.MakeBasicDisplayer(logger, loghelper.StdDisplay)
			conf := &config.Config{
				Displayer: displayer,
			}

			// Create mock retriever
			retriever := &mockRetriever{
				installFunc: func(ctx context.Context, version string, targetPath string) error {
					return nil
				},
				listVersionsFunc: func(ctx context.Context) ([]string, error) {
					return []string{"1.0.0", "1.1.0"}, nil
				},
			}

			// Create version manager
			vm := versionmanager.Make(conf, "TEST_", "test", nil, retriever, nil)

			// Create a mock builder function
			builderFunc := func(conf *config.Config, parser *hclparse.Parser) versionmanager.VersionManager {
				return vm
			}

			// Save original builder
			originalBuilder := builder.Builders[cmdconst.TfName]
			builder.Builders[cmdconst.TfName] = builderFunc
			defer func() {
				builder.Builders[cmdconst.TfName] = originalBuilder
			}()

			// Execute command
			lightproxy.Exec(cmdconst.TfName)
		})
	}
}
