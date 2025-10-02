/*
 *
 * Copyright 2024 tofuutils authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package versionmanager

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/parser/iac"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/types"
)

// MockReleaseRetriever is a mock implementation of ReleaseRetriever
type MockReleaseRetriever struct {
	mock.Mock
}

func (m *MockReleaseRetriever) Install(ctx context.Context, version string, targetPath string) error {
	args := m.Called(ctx, version, targetPath)
	return args.Error(0)
}

func (m *MockReleaseRetriever) ListVersions(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

// MockDisplayer is a mock implementation of the displayer interface
type MockDisplayer struct {
	mock.Mock
	messages []string
}

func (m *MockDisplayer) Display(msg string) {
	m.Called(msg)
	m.messages = append(m.messages, msg)
}

func (m *MockDisplayer) Log(level hclog.Level, msg string, args ...any) {
	m.Called(level, msg, args)
}

func (m *MockDisplayer) Flush(proxyCall bool) {
	m.Called(proxyCall)
}

func (m *MockDisplayer) IsDebug() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockDisplayer) GetMessages() []string {
	return m.messages
}

func TestMake(t *testing.T) {
	tests := []struct {
		name         string
		conf         *config.Config
		envPrefix    string
		folderName   string
		iacExts      []iacparser.ExtDescription
		retriever    ReleaseRetriever
		versionFiles []types.VersionFile
	}{
		{
			name:       "basic_make",
			conf:       &config.Config{},
			envPrefix:  "TF",
			folderName: "terraform",
			iacExts:    []iacparser.ExtDescription{},
			retriever:  &MockReleaseRetriever{},
			versionFiles: []types.VersionFile{
				{Name: ".terraform-version"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Make(tt.conf, tt.envPrefix, tt.folderName, tt.iacExts, tt.retriever, tt.versionFiles)

			assert.NotNil(t, result)
			assert.Equal(t, tt.conf, result.Conf)
			assert.Equal(t, EnvPrefix(tt.envPrefix), result.EnvNames)
			assert.Equal(t, tt.folderName, result.FolderName)
			assert.Equal(t, tt.iacExts, result.iacExts)
			assert.Equal(t, tt.retriever, result.retriever)
			assert.Equal(t, tt.versionFiles, result.VersionFiles)
		})
	}
}

func TestVersionManager_InstallPath(t *testing.T) {
	tests := []struct {
		name       string
		rootPath   string
		folderName string
		expected   string
	}{
		{
			name:       "basic_path",
			rootPath:   "/tmp",
			folderName: "terraform",
			expected:   "/tmp/terraform",
		},
		{
			name:       "nested_path",
			rootPath:   "/opt/tools",
			folderName: "tofu",
			expected:   "/opt/tools/tofu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for testing
			tempDir, err := os.MkdirTemp("", "tenv_test")
			assert.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Override root path to use temp directory
			originalRootPath := tt.rootPath
			if originalRootPath == "/tmp" || originalRootPath == "/opt/tools" {
				originalRootPath = tempDir
			}

			displayer := &MockDisplayer{}
			conf := &config.Config{
				RootPath:  originalRootPath,
				Displayer: displayer,
			}

			m := VersionManager{
				Conf:       conf,
				FolderName: tt.folderName,
			}

			result, err := m.InstallPath()

			assert.NoError(t, err)
			expectedPath := filepath.Join(originalRootPath, tt.folderName)
			assert.Equal(t, expectedPath, result)

			// Verify directory was created
			_, err = os.Stat(result)
			assert.NoError(t, err)
		})
	}
}

func TestVersionManager_ListLocal(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "tenv_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create version directories
	versions := []string{"1.0.0", "1.1.0", "1.2.0"}
	for _, version := range versions {
		versionDir := filepath.Join(tempDir, "terraform", version)
		err := os.MkdirAll(versionDir, 0755)
		assert.NoError(t, err)
	}

	displayer := &MockDisplayer{}
	displayer.On("Log", mock.Anything, mock.Anything, mock.Anything).Maybe() // Allow any log calls

	conf := &config.Config{
		RootPath:  tempDir,
		Displayer: displayer,
	}

	m := VersionManager{
		Conf:       conf,
		FolderName: "terraform",
	}

	result, err := m.ListLocal(false)

	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// Check that versions are sorted (should be in ascending order by default)
	expectedVersions := []string{"1.0.0", "1.1.0", "1.2.0"}
	for i, expected := range expectedVersions {
		assert.Equal(t, expected, result[i].Version)
		// UseDate might be zero if the file doesn't exist, so we just check it's a valid time
		assert.True(t, result[i].UseDate.Year() >= 1, "UseDate should be a valid time")
	}

	displayer.AssertExpectations(t)
}

func TestVersionManager_ListRemote(t *testing.T) {
	mockRetriever := &MockReleaseRetriever{}
	expectedVersions := []string{"1.5.0", "1.4.0", "1.3.0"}

	mockRetriever.On("ListVersions", mock.Anything).Return(expectedVersions, nil)

	displayer := &MockDisplayer{}
	conf := &config.Config{
		Displayer: displayer,
	}

	m := VersionManager{
		Conf:      conf,
		retriever: mockRetriever,
	}

	ctx := context.Background()
	result, err := m.ListRemote(ctx, false)

	assert.NoError(t, err)
	assert.Equal(t, expectedVersions, result)
	mockRetriever.AssertExpectations(t)
}

func TestVersionManager_LocalSet(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "tenv_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create version directories
	versions := []string{"1.0.0", "1.1.0", "1.2.0"}
	for _, version := range versions {
		versionDir := filepath.Join(tempDir, "terraform", version)
		err := os.MkdirAll(versionDir, 0755)
		assert.NoError(t, err)
	}

	displayer := &MockDisplayer{}
	conf := &config.Config{
		RootPath:  tempDir,
		Displayer: displayer,
	}

	m := VersionManager{
		Conf:       conf,
		FolderName: "terraform",
	}

	result := m.LocalSet()

	assert.NotNil(t, result)
	assert.Len(t, result, 3)

	// Check that all expected versions are in the set
	for _, version := range versions {
		_, exists := result[version]
		assert.True(t, exists, "Version %s should be in the set", version)
	}
}

func TestVersionManager_Resolve(t *testing.T) {
	tests := []struct {
		name             string
		envVersion       string
		versionFiles     []types.VersionFile
		expectedVersion  string
		expectedStrategy string
	}{
		{
			name:             "resolve_from_env",
			envVersion:       "1.2.0",
			expectedVersion:  "1.2.0",
			expectedStrategy: "",
		},
		{
			name:             "resolve_fallback_to_strategy",
			envVersion:       "",
			versionFiles:     []types.VersionFile{},
			expectedVersion:  semantic.LatestAllowedKey,
			expectedStrategy: semantic.LatestAllowedKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			displayer := &MockDisplayer{}
			displayer.On("Display", mock.Anything).Maybe()
			displayer.On("Log", mock.Anything, mock.Anything, mock.Anything).Maybe()

			// Create a proper config with a mock Getenv function
			conf := &config.Config{
				Displayer: displayer,
				Getenv:    func(key string) string { return os.Getenv(key) },
			}

			// Mock environment variable by setting it directly
			if tt.envVersion != "" {
				originalEnv := os.Getenv("TFVERSION")
				os.Setenv("TFVERSION", tt.envVersion)
				defer func() {
					if originalEnv == "" {
						os.Unsetenv("TFVERSION")
					} else {
						os.Setenv("TFVERSION", originalEnv)
					}
				}()
			}

			m := VersionManager{
				Conf:         conf,
				EnvNames:     "TF",
				FolderName:   "terraform",
				VersionFiles: tt.versionFiles,
			}

			result, err := m.Resolve(tt.expectedStrategy)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedVersion, result)

			displayer.AssertExpectations(t)
		})
	}
}

func TestVersionManager_ResolveWithVersionFiles(t *testing.T) {
	displayer := &MockDisplayer{}
	conf := &config.Config{
		Displayer: displayer,
	}

	m := VersionManager{
		Conf: conf,
	}

	result, err := m.ResolveWithVersionFiles()

	// This should not error, but may return empty string if no version files
	assert.NoError(t, err)
	assert.Equal(t, "", result) // No version files configured
}

func TestVersionManager_RootConstraintFilePath(t *testing.T) {
	displayer := &MockDisplayer{}
	conf := &config.Config{
		RootPath:  "/opt/tenv",
		Displayer: displayer,
	}

	m := VersionManager{
		Conf:       conf,
		FolderName: "terraform",
	}

	result := m.RootConstraintFilePath()

	expected := filepath.Join("/opt/tenv", "terraform", "constraint")
	assert.Equal(t, expected, result)
}

func TestVersionManager_RootVersionFilePath(t *testing.T) {
	displayer := &MockDisplayer{}
	conf := &config.Config{
		RootPath:  "/opt/tenv",
		Displayer: displayer,
	}

	m := VersionManager{
		Conf:       conf,
		FolderName: "terraform",
	}

	result := m.RootVersionFilePath()

	expected := filepath.Join("/opt/tenv", "terraform", "version")
	assert.Equal(t, expected, result)
}

func TestVersionManager_SetConstraint(t *testing.T) {
	tests := []struct {
		name        string
		constraint  string
		expectError bool
	}{
		{
			name:        "valid_constraint",
			constraint:  ">= 1.0.0",
			expectError: false,
		},
		{
			name:        "invalid_constraint",
			constraint:  "invalid-constraint",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for testing
			tempDir, err := os.MkdirTemp("", "tenv_test")
			assert.NoError(t, err)
			defer os.RemoveAll(tempDir)

			displayer := &MockDisplayer{}
			displayer.On("Display", mock.Anything).Maybe()
			conf := &config.Config{
				RootPath:  tempDir,
				Displayer: displayer,
			}

			// Create the terraform directory first
			terraformDir := filepath.Join(tempDir, "terraform")
			err = os.MkdirAll(terraformDir, 0755)
			assert.NoError(t, err)

			m := VersionManager{
				Conf:       conf,
				FolderName: "terraform",
			}

			err = m.SetConstraint(tt.constraint)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify file was created with correct content
				constraintFile := filepath.Join(tempDir, "terraform", "constraint")
				content, err := os.ReadFile(constraintFile)
				assert.NoError(t, err)
				assert.Equal(t, tt.constraint, string(content))
			}
		})
	}
}

func TestVersionManager_UninstallMultiple(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "tenv_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create version directories
	versions := []string{"1.0.0", "1.1.0"}
	for _, version := range versions {
		versionDir := filepath.Join(tempDir, "terraform", version)
		err := os.MkdirAll(versionDir, 0755)
		assert.NoError(t, err)
	}

	displayer := &MockDisplayer{}
	displayer.On("Display", mock.Anything).Maybe()
	conf := &config.Config{
		RootPath:  tempDir,
		Displayer: displayer,
	}

	m := VersionManager{
		Conf:       conf,
		FolderName: "terraform",
	}

	err = m.UninstallMultiple(versions)

	assert.NoError(t, err)

	// Verify directories were removed
	for _, version := range versions {
		versionDir := filepath.Join(tempDir, "terraform", version)
		_, err := os.Stat(versionDir)
		assert.True(t, os.IsNotExist(err), "Directory %s should be removed", versionDir)
	}
}

func TestVersionManager_checkVersionInstallation(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "tenv_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a version directory
	versionDir := filepath.Join(tempDir, "terraform", "1.2.0")
	err = os.MkdirAll(versionDir, 0755)
	assert.NoError(t, err)

	displayer := &MockDisplayer{}
	conf := &config.Config{
		RootPath:  tempDir,
		Displayer: displayer,
	}

	m := VersionManager{
		Conf:       conf,
		FolderName: "terraform",
	}

	tests := []struct {
		name         string
		installPath  string
		version      string
		expectedPath string
		expectedBool bool
	}{
		{
			name:         "version_exists",
			installPath:  "",
			version:      "1.2.0",
			expectedPath: filepath.Join(tempDir, "terraform"),
			expectedBool: true,
		},
		{
			name:         "version_does_not_exist",
			installPath:  "",
			version:      "1.3.0",
			expectedPath: filepath.Join(tempDir, "terraform"),
			expectedBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultPath, resultBool, err := m.checkVersionInstallation(tt.installPath, tt.version)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPath, resultPath)
			assert.Equal(t, tt.expectedBool, resultBool)
		})
	}
}

func TestVersionManager_innerListLocal(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "tenv_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create version directories
	versions := []string{"1.0.0", "1.1.0", "1.2.0"}
	for _, version := range versions {
		versionDir := filepath.Join(tempDir, version)
		err := os.MkdirAll(versionDir, 0755)
		assert.NoError(t, err)
	}

	m := VersionManager{}

	result, err := m.innerListLocal(tempDir, false)

	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// Check that versions are sorted
	expectedVersions := []string{"1.0.0", "1.1.0", "1.2.0"}
	assert.Equal(t, expectedVersions, result)
}
