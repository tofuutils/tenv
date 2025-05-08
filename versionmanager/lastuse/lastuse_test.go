package lastuse

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/tofuutils/tenv/v4/config"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		setup   func(string) error
		conf    *config.Config
		want    time.Time
		wantErr bool
	}{
		{
			name: "valid lastuse file",
			path: "test",
			setup: func(path string) error {
				now := time.Now()
				return os.WriteFile(filepath.Join(path, ".lastuse"), []byte(now.Format(time.RFC3339)), 0644)
			},
			conf: &config.Config{
				Displayer: hclog.NewNullLogger(),
			},
			want:    time.Now(),
			wantErr: false,
		},
		{
			name: "invalid lastuse file",
			path: "test",
			setup: func(path string) error {
				return os.WriteFile(filepath.Join(path, ".lastuse"), []byte("invalid"), 0644)
			},
			conf: &config.Config{
				Displayer: hclog.NewNullLogger(),
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "missing lastuse file",
			path: "test",
			setup: func(path string) error {
				return nil
			},
			conf: &config.Config{
				Displayer: hclog.NewNullLogger(),
			},
			want:    time.Time{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := filepath.Join(tempDir, tt.path)
			if err := os.MkdirAll(path, 0755); err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}

			if err := tt.setup(path); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			got := Read(path, tt.conf)
			if !tt.wantErr && got.IsZero() {
				t.Errorf("Read() returned zero time when it shouldn't")
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		conf    *config.Config
		wantErr bool
	}{
		{
			name: "successful write",
			path: "test",
			conf: &config.Config{
				Displayer: hclog.NewNullLogger(),
			},
			wantErr: false,
		},
		{
			name: "invalid path",
			path: "/invalid/path",
			conf: &config.Config{
				Displayer: hclog.NewNullLogger(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := filepath.Join(tempDir, tt.path)
			if err := os.MkdirAll(path, 0755); err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}

			err := Write(path, tt.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify the file was created
				if _, err := os.Stat(filepath.Join(path, ".lastuse")); os.IsNotExist(err) {
					t.Error("Write() did not create .lastuse file")
				}

				// Verify the content is a valid timestamp
				content, err := os.ReadFile(filepath.Join(path, ".lastuse"))
				if err != nil {
					t.Errorf("Failed to read .lastuse file: %v", err)
				}

				_, err = time.Parse(time.RFC3339, string(content))
				if err != nil {
					t.Errorf("Invalid timestamp in .lastuse file: %v", err)
				}
			}
		})
	}
}
