package builder

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/versionmanager"
)

type mockRetriever struct {
	versionmanager.ReleaseRetriever
	installFunc      func(ctx context.Context, version string, targetPath string) error
	listVersionsFunc func(ctx context.Context) ([]string, error)
}

func (m *mockRetriever) Install(ctx context.Context, version string, targetPath string) error {
	return m.installFunc(ctx, version, targetPath)
}

func (m *mockRetriever) ListVersions(ctx context.Context) ([]string, error) {
	return m.listVersionsFunc(ctx)
}

func TestBuildVersionManager(t *testing.T) {
	tests := []struct {
		name      string
		conf      *config.Config
		retriever *mockRetriever
		wantErr   bool
	}{
		{
			name: "successful build",
			conf: &config.Config{
				Displayer: hclog.NewNullLogger(),
			},
			retriever: &mockRetriever{
				listVersionsFunc: func(ctx context.Context) ([]string, error) {
					return []string{"1.0.0", "1.1.0"}, nil
				},
			},
			wantErr: false,
		},
		{
			name: "invalid config",
			conf: nil,
			retriever: &mockRetriever{
				listVersionsFunc: func(ctx context.Context) ([]string, error) {
					return []string{"1.0.0", "1.1.0"}, nil
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildVersionManager(tt.conf, tt.retriever)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildVersionManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildRetriever(t *testing.T) {
	tests := []struct {
		name    string
		conf    *config.Config
		wantErr bool
	}{
		{
			name: "successful build",
			conf: &config.Config{
				Displayer: hclog.NewNullLogger(),
			},
			wantErr: false,
		},
		{
			name:    "invalid config",
			conf:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildRetriever(tt.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildRetriever() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildVersionManagerWithRetriever(t *testing.T) {
	tests := []struct {
		name      string
		conf      *config.Config
		retriever *mockRetriever
		wantErr   bool
	}{
		{
			name: "successful build",
			conf: &config.Config{
				Displayer: hclog.NewNullLogger(),
			},
			retriever: &mockRetriever{
				listVersionsFunc: func(ctx context.Context) ([]string, error) {
					return []string{"1.0.0", "1.1.0"}, nil
				},
			},
			wantErr: false,
		},
		{
			name: "invalid config",
			conf: nil,
			retriever: &mockRetriever{
				listVersionsFunc: func(ctx context.Context) ([]string, error) {
					return []string{"1.0.0", "1.1.0"}, nil
				},
			},
			wantErr: true,
		},
		{
			name: "invalid retriever",
			conf: &config.Config{
				Displayer: hclog.NewNullLogger(),
			},
			retriever: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildVersionManagerWithRetriever(tt.conf, tt.retriever)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildVersionManagerWithRetriever() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
