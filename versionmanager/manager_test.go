package versionmanager

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/parser/iac"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/types"
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

func TestVersionManager_Detect(t *testing.T) {
	tests := []struct {
		name        string
		conf        *config.Config
		retriever   *mockRetriever
		proxyCall   bool
		wantVersion string
		wantErr     bool
	}{
		{
			name: "successful detection",
			conf: &config.Config{
				SkipInstall: true,
				Displayer:   hclog.NewNullLogger(),
			},
			retriever: &mockRetriever{
				listVersionsFunc: func(ctx context.Context) ([]string, error) {
					return []string{"1.0.0", "1.1.0"}, nil
				},
			},
			proxyCall:   false,
			wantVersion: "1.1.0",
			wantErr:     false,
		},
		{
			name: "error in retriever",
			conf: &config.Config{
				SkipInstall: true,
				Displayer:   hclog.NewNullLogger(),
			},
			retriever: &mockRetriever{
				listVersionsFunc: func(ctx context.Context) ([]string, error) {
					return nil, errors.New("retriever error")
				},
			},
			proxyCall:   false,
			wantVersion: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := Make(tt.conf, "TEST", "test", []iac.ExtDescription{}, tt.retriever, []types.VersionFile{})
			got, err := vm.Detect(context.Background(), tt.proxyCall)
			if (err != nil) != tt.wantErr {
				t.Errorf("VersionManager.Detect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantVersion {
				t.Errorf("VersionManager.Detect() = %v, want %v", got, tt.wantVersion)
			}
		})
	}
}

func TestVersionManager_InstallPath(t *testing.T) {
	tests := []struct {
		name    string
		conf    *config.Config
		folder  string
		want    string
		wantErr bool
	}{
		{
			name: "successful path creation",
			conf: &config.Config{
				RootPath: t.TempDir(),
			},
			folder:  "test",
			want:    filepath.Join(t.TempDir(), "test"),
			wantErr: false,
		},
		{
			name: "invalid root path",
			conf: &config.Config{
				RootPath: "/invalid/path",
			},
			folder:  "test",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := Make(tt.conf, "TEST", tt.folder, []iac.ExtDescription{}, &mockRetriever{}, []types.VersionFile{})
			got, err := vm.InstallPath()
			if (err != nil) != tt.wantErr {
				t.Errorf("VersionManager.InstallPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want && !tt.wantErr {
				t.Errorf("VersionManager.InstallPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionManager_LocalSet(t *testing.T) {
	tests := []struct {
		name   string
		conf   *config.Config
		folder string
		setup  func(string) error
		want   map[string]struct{}
	}{
		{
			name: "empty directory",
			conf: &config.Config{
				RootPath:  t.TempDir(),
				Displayer: hclog.NewNullLogger(),
			},
			folder: "test",
			setup:  func(path string) error { return nil },
			want:   map[string]struct{}{},
		},
		{
			name: "with versions",
			conf: &config.Config{
				RootPath:  t.TempDir(),
				Displayer: hclog.NewNullLogger(),
			},
			folder: "test",
			setup: func(path string) error {
				versions := []string{"1.0.0", "1.1.0"}
				for _, v := range versions {
					if err := os.MkdirAll(filepath.Join(path, v), 0755); err != nil {
						return err
					}
				}
				return nil
			},
			want: map[string]struct{}{
				"1.0.0": {},
				"1.1.0": {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installPath := filepath.Join(tt.conf.RootPath, tt.folder)
			if err := tt.setup(installPath); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			vm := Make(tt.conf, "TEST", tt.folder, []iac.ExtDescription{}, &mockRetriever{}, []types.VersionFile{})
			got := vm.LocalSet()

			if len(got) != len(tt.want) {
				t.Errorf("VersionManager.LocalSet() = %v, want %v", got, tt.want)
			}

			for version := range tt.want {
				if _, ok := got[version]; !ok {
					t.Errorf("VersionManager.LocalSet() missing version %v", version)
				}
			}
		})
	}
}

func TestVersionManager_ListLocal(t *testing.T) {
	tests := []struct {
		name         string
		conf         *config.Config
		folder       string
		setup        func(string) error
		reverseOrder bool
		want         []DatedVersion
		wantErr      bool
	}{
		{
			name: "empty directory",
			conf: &config.Config{
				RootPath:  t.TempDir(),
				Displayer: hclog.NewNullLogger(),
			},
			folder:       "test",
			setup:        func(path string) error { return nil },
			reverseOrder: false,
			want:         []DatedVersion{},
			wantErr:      false,
		},
		{
			name: "with versions",
			conf: &config.Config{
				RootPath:  t.TempDir(),
				Displayer: hclog.NewNullLogger(),
			},
			folder: "test",
			setup: func(path string) error {
				versions := []string{"1.0.0", "1.1.0"}
				for _, v := range versions {
					if err := os.MkdirAll(filepath.Join(path, v), 0755); err != nil {
						return err
					}
				}
				return nil
			},
			reverseOrder: false,
			want: []DatedVersion{
				{Version: "1.0.0", UseDate: time.Time{}},
				{Version: "1.1.0", UseDate: time.Time{}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installPath := filepath.Join(tt.conf.RootPath, tt.folder)
			if err := tt.setup(installPath); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			vm := Make(tt.conf, "TEST", tt.folder, []iac.ExtDescription{}, &mockRetriever{}, []types.VersionFile{})
			got, err := vm.ListLocal(tt.reverseOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("VersionManager.ListLocal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("VersionManager.ListLocal() = %v, want %v", got, tt.want)
			}

			for i, version := range got {
				if version.Version != tt.want[i].Version {
					t.Errorf("VersionManager.ListLocal() version = %v, want %v", version.Version, tt.want[i].Version)
				}
			}
		})
	}
}

func TestVersionManager_ListRemote(t *testing.T) {
	tests := []struct {
		name         string
		retriever    *mockRetriever
		reverseOrder bool
		want         []string
		wantErr      bool
	}{
		{
			name: "successful list",
			retriever: &mockRetriever{
				listVersionsFunc: func(ctx context.Context) ([]string, error) {
					return []string{"1.0.0", "1.1.0"}, nil
				},
			},
			reverseOrder: false,
			want:         []string{"1.0.0", "1.1.0"},
			wantErr:      false,
		},
		{
			name: "error in retriever",
			retriever: &mockRetriever{
				listVersionsFunc: func(ctx context.Context) ([]string, error) {
					return nil, errors.New("retriever error")
				},
			},
			reverseOrder: false,
			want:         nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := Make(&config.Config{}, "TEST", "test", []iac.ExtDescription{}, tt.retriever, []types.VersionFile{})
			got, err := vm.ListRemote(context.Background(), tt.reverseOrder)
			if (err != nil) != tt.wantErr {
				t.Errorf("VersionManager.ListRemote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("VersionManager.ListRemote() = %v, want %v", got, tt.want)
			}

			for i, version := range got {
				if version != tt.want[i] {
					t.Errorf("VersionManager.ListRemote() version = %v, want %v", version, tt.want[i])
				}
			}
		})
	}
}
