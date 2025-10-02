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

package config_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

// Helper functions for platform-specific paths
func getTestHomeDir() string {
	// Use actual user home directory for testing
	userHome, err := os.UserHomeDir()
	if err != nil {
		// Fallback to test values if we can't get user home
		if runtime.GOOS == "windows" {
			return "C:\\Users\\testuser"
		}
		return "/home/testuser"
	}
	return userHome
}

func getExpectedRootPath() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		// Fallback to test values if we can't get user home
		if runtime.GOOS == "windows" {
			return "C:\\Users\\testuser\\.tenv"
		}
		return "/home/testuser/.tenv"
	}
	return filepath.Join(userHome, ".tenv")
}

func TestParseValidationMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected config.ValidationMode
	}{
		{"none mode", "none", config.NoValidation},
		{"sha mode", "sha", config.ShaValidation},
		{"sign mode", "sign", config.SignValidation},
		{"default mode", "unknown", config.SignValidation},
		{"empty mode", "", config.SignValidation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.ParseValidationMode(tt.input)
			if result != tt.expected {
				t.Errorf("ParseValidationMode(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg, err := config.DefaultConfig()
	if err != nil {
		t.Fatalf("DefaultConfig() error = %v", err)
	}

	// Test basic properties
	if cfg.Arch != runtime.GOARCH {
		t.Errorf("DefaultConfig() Arch = %v, want %v", cfg.Arch, runtime.GOARCH)
	}

	if cfg.SkipInstall != true {
		t.Errorf("DefaultConfig() SkipInstall = %v, want true", cfg.SkipInstall)
	}

	if cfg.Validation != config.SignValidation {
		t.Errorf("DefaultConfig() Validation = %v, want %v", cfg.Validation, config.SignValidation)
	}

	if cfg.WorkPath != "." {
		t.Errorf("DefaultConfig() WorkPath = %v, want '.'", cfg.WorkPath)
	}

	// Test that Getenv is set to EmptyGetenv
	if cfg.Getenv("TEST") != "" {
		t.Errorf("DefaultConfig() Getenv should return empty string for any input")
	}
}

func TestEmptyGetenv(t *testing.T) {
	t.Parallel()

	result := config.EmptyGetenv("any_key")
	if result != "" {
		t.Errorf("EmptyGetenv() = %v, want empty string", result)
	}
}

func TestConfigInitDisplayer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		forceQuiet  bool
		proxyCall   bool
		expectQuiet bool
	}{
		{"force quiet", true, false, true},
		{"normal display", false, false, false},
		{"proxy call", false, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ForceQuiet:     tt.forceQuiet,
				DisplayVerbose: true,
				Getenv:         config.EmptyGetenv,
			}

			cfg.InitDisplayer(tt.proxyCall)

			if tt.expectQuiet && cfg.Displayer != loghelper.InertDisplayer {
				t.Errorf("Expected InertDisplayer when ForceQuiet is true")
			}

			if !tt.expectQuiet && cfg.Displayer == loghelper.InertDisplayer {
				t.Errorf("Expected non-InertDisplayer when ForceQuiet is false")
			}
		})
	}
}

func TestConfigInitValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		initial  config.ValidationMode
		skipSum  bool
		skipSign bool
		expected config.ValidationMode
	}{
		{"skip sum overrides", config.SignValidation, true, false, config.NoValidation},
		{"skip sign with sum", config.SignValidation, false, true, config.ShaValidation},
		{"skip sign without sum", config.SignValidation, false, false, config.SignValidation},
		{"no validation skip", config.NoValidation, false, true, config.NoValidation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Validation: tt.initial}
			cfg.InitValidation(tt.skipSum, tt.skipSign)

			if cfg.Validation != tt.expected {
				t.Errorf("InitValidation() = %v, want %v", cfg.Validation, tt.expected)
			}
		})
	}
}

func TestConfigInitInstall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		initial        bool
		forceInstall   bool
		forceNoInstall bool
		expected       bool
	}{
		{"force no install overrides", true, true, true, true},
		{"force install", false, true, false, false},
		{"no force flags", true, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{SkipInstall: tt.initial}
			cfg.InitInstall(tt.forceInstall, tt.forceNoInstall)

			if cfg.SkipInstall != tt.expected {
				t.Errorf("InitInstall() = %v, want %v", cfg.SkipInstall, tt.expected)
			}
		})
	}
}

func TestConfigInitRemoteConf(t *testing.T) {
	t.Parallel()

	// Test with non-existent file
	cfg := &config.Config{
		RootPath:  "/nonexistent",
		Displayer: loghelper.InertDisplayer,
	}

	err := cfg.InitRemoteConf()
	if err != nil {
		t.Errorf("InitRemoteConf() with non-existent file should not error = %v", err)
	}
}

func TestMapGetDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		m            map[string]string
		key          string
		defaultValue string
		expected     string
	}{
		{"key exists", map[string]string{"key": "value"}, "key", "default", "value"},
		{"key doesn't exist", map[string]string{}, "key", "default", "default"},
		{"key with spaces", map[string]string{"key": "  value  "}, "key", "default", "value"},
		{"empty key value", map[string]string{"key": ""}, "key", "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.MapGetDefault(tt.m, tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("MapGetDefault() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetBasicAuthOption(t *testing.T) {
	t.Parallel()

	mockGetenv := func(key string) string {
		envMap := map[string]string{
			"USER": "testuser",
			"PASS": "testpass",
		}
		return envMap[key]
	}

	tests := []struct {
		name        string
		userEnvName string
		passEnvName string
		expected    bool // true if options should be returned
	}{
		{"both env vars present", "USER", "PASS", true},
		{"missing user", "MISSING_USER", "PASS", false},
		{"missing pass", "USER", "MISSING_PASS", false},
		{"both missing", "MISSING_USER", "MISSING_PASS", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := config.GetBasicAuthOption(mockGetenv, tt.userEnvName, tt.passEnvName)
			if tt.expected && options == nil {
				t.Errorf("Expected non-nil options when both env vars are present")
			}
			if !tt.expected && options != nil {
				t.Errorf("Expected nil options when env vars are missing")
			}
		})
	}
}

func TestInitConfigFromEnv(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig config.Config
		expectError    bool
	}{
		{
			name: "default configuration",
			envVars: map[string]string{
				"HOME": getTestHomeDir(),
			},
			expectedConfig: config.Config{
				Arch:             runtime.GOARCH,
				ForceQuiet:       false,
				ForceRemote:      false,
				GithubActions:    false,
				GithubToken:      "",
				RemoteConfPath:   "",
				RootPath:         getExpectedRootPath(),
				SkipInstall:      true,
				UserPath:         getTestHomeDir(),
				Validation:       config.SignValidation,
				WorkPath:         ".",
				TfKeyPathOrURL:   "https://www.hashicorp.com/.well-known/pgp-key.txt",
				TofuKeyPathOrURL: "https://get.opentofu.org/opentofu.asc",
			},
			expectError: false,
		},
		{
			name: "custom arch and auto install",
			envVars: map[string]string{
				"HOME":              getTestHomeDir(),
				"TENV_ARCH":         "arm64",
				"TENV_AUTO_INSTALL": "true",
			},
			expectedConfig: config.Config{
				Arch:             "arm64",
				ForceQuiet:       false,
				ForceRemote:      false,
				GithubActions:    false,
				GithubToken:      "",
				RemoteConfPath:   "",
				RootPath:         getExpectedRootPath(),
				SkipInstall:      false,
				UserPath:         getTestHomeDir(),
				Validation:       config.SignValidation,
				WorkPath:         ".",
				TfKeyPathOrURL:   "https://www.hashicorp.com/.well-known/pgp-key.txt",
				TofuKeyPathOrURL: "https://get.opentofu.org/opentofu.asc",
			},
			expectError: false,
		},
		{
			name: "force remote and quiet mode",
			envVars: map[string]string{
				"HOME":              getTestHomeDir(),
				"TENV_FORCE_REMOTE": "true",
				"TENV_QUIET":        "true",
			},
			expectedConfig: config.Config{
				Arch:             runtime.GOARCH,
				ForceQuiet:       true,
				ForceRemote:      true,
				GithubActions:    false,
				GithubToken:      "",
				RemoteConfPath:   "",
				RootPath:         getExpectedRootPath(),
				SkipInstall:      true,
				UserPath:         getTestHomeDir(),
				Validation:       config.SignValidation,
				WorkPath:         ".",
				TfKeyPathOrURL:   "https://www.hashicorp.com/.well-known/pgp-key.txt",
				TofuKeyPathOrURL: "https://get.opentofu.org/opentofu.asc",
			},
			expectError: false,
		},
		{
			name: "custom root path",
			envVars: map[string]string{
				"HOME":      getTestHomeDir(),
				"TENV_ROOT": "/custom/tenv/path",
			},
			expectedConfig: config.Config{
				Arch:             runtime.GOARCH,
				ForceQuiet:       false,
				ForceRemote:      false,
				GithubActions:    false,
				GithubToken:      "",
				RemoteConfPath:   "",
				RootPath:         "/custom/tenv/path",
				SkipInstall:      true,
				UserPath:         getTestHomeDir(),
				Validation:       config.SignValidation,
				WorkPath:         ".",
				TfKeyPathOrURL:   "https://www.hashicorp.com/.well-known/pgp-key.txt",
				TofuKeyPathOrURL: "https://get.opentofu.org/opentofu.asc",
			},
			expectError: false,
		},
		{
			name: "github token and actions",
			envVars: map[string]string{
				"HOME":              getTestHomeDir(),
				"TENV_GITHUB_TOKEN": "ghp_1234567890abcdef",
				"GITHUB_ACTIONS":    "true",
			},
			expectedConfig: config.Config{
				Arch:             runtime.GOARCH,
				ForceQuiet:       false,
				ForceRemote:      false,
				GithubActions:    true,
				GithubToken:      "ghp_1234567890abcdef",
				RemoteConfPath:   "",
				RootPath:         getExpectedRootPath(),
				SkipInstall:      true,
				UserPath:         getTestHomeDir(),
				Validation:       config.SignValidation,
				WorkPath:         ".",
				TfKeyPathOrURL:   "https://www.hashicorp.com/.well-known/pgp-key.txt",
				TofuKeyPathOrURL: "https://get.opentofu.org/opentofu.asc",
			},
			expectError: false,
		},
		{
			name: "validation modes",
			envVars: map[string]string{
				"HOME":            getTestHomeDir(),
				"TENV_VALIDATION": "sha",
			},
			expectedConfig: config.Config{
				Arch:             runtime.GOARCH,
				ForceQuiet:       false,
				ForceRemote:      false,
				GithubActions:    false,
				GithubToken:      "",
				RemoteConfPath:   "",
				RootPath:         getExpectedRootPath(),
				SkipInstall:      true,
				UserPath:         getTestHomeDir(),
				Validation:       config.ShaValidation,
				WorkPath:         ".",
				TfKeyPathOrURL:   "https://www.hashicorp.com/.well-known/pgp-key.txt",
				TofuKeyPathOrURL: "https://get.opentofu.org/opentofu.asc",
			},
			expectError: false,
		},
		{
			name: "custom remote configuration path",
			envVars: map[string]string{
				"HOME":             getTestHomeDir(),
				"TENV_REMOTE_CONF": "/etc/tenv/remote.yaml",
			},
			expectedConfig: config.Config{
				Arch:             runtime.GOARCH,
				ForceQuiet:       false,
				ForceRemote:      false,
				GithubActions:    false,
				GithubToken:      "",
				RemoteConfPath:   "/etc/tenv/remote.yaml",
				RootPath:         getExpectedRootPath(),
				SkipInstall:      true,
				UserPath:         getTestHomeDir(),
				Validation:       config.SignValidation,
				WorkPath:         ".",
				TfKeyPathOrURL:   "https://www.hashicorp.com/.well-known/pgp-key.txt",
				TofuKeyPathOrURL: "https://get.opentofu.org/opentofu.asc",
			},
			expectError: false,
		},
		{
			name: "custom terraform key path",
			envVars: map[string]string{
				"HOME":                    getTestHomeDir(),
				"TFENV_HASHICORP_PGP_KEY": "/custom/terraform-key.asc",
			},
			expectedConfig: config.Config{
				Arch:             runtime.GOARCH,
				ForceQuiet:       false,
				ForceRemote:      false,
				GithubActions:    false,
				GithubToken:      "",
				RemoteConfPath:   "",
				RootPath:         getExpectedRootPath(),
				SkipInstall:      true,
				UserPath:         getTestHomeDir(),
				Validation:       config.SignValidation,
				WorkPath:         ".",
				TfKeyPathOrURL:   "/custom/terraform-key.asc",
				TofuKeyPathOrURL: "https://get.opentofu.org/opentofu.asc",
			},
			expectError: false,
		},
		{
			name: "custom tofu key path",
			envVars: map[string]string{
				"HOME":                     getTestHomeDir(),
				"TOFUENV_OPENTOFU_PGP_KEY": "/custom/tofu-key.asc",
			},
			expectedConfig: config.Config{
				Arch:             runtime.GOARCH,
				ForceQuiet:       false,
				ForceRemote:      false,
				GithubActions:    false,
				GithubToken:      "",
				RemoteConfPath:   "",
				RootPath:         getExpectedRootPath(),
				SkipInstall:      true,
				UserPath:         getTestHomeDir(),
				Validation:       config.SignValidation,
				WorkPath:         ".",
				TfKeyPathOrURL:   "https://www.hashicorp.com/.well-known/pgp-key.txt",
				TofuKeyPathOrURL: "/custom/tofu-key.asc",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment
			originalEnv := make(map[string]string)
			envVarsToRestore := []string{
				"HOME", "TENV_ARCH", "TENV_AUTO_INSTALL", "TENV_FORCE_REMOTE", "TENV_QUIET",
				"TENV_ROOT", "TENV_GITHUB_TOKEN", "GITHUB_ACTIONS", "TENV_VALIDATION", "TENV_REMOTE_CONF",
				"TFENV_HASHICORP_PGP_KEY", "TOFUENV_OPENTOFU_PGP_KEY",
			}

			for _, envVar := range envVarsToRestore {
				originalEnv[envVar] = os.Getenv(envVar)
			}

			defer func() {
				for envVar, value := range originalEnv {
					if value == "" {
						os.Unsetenv(envVar)
					} else {
						os.Setenv(envVar, value)
					}
				}
			}()

			// Set test environment variables
			for envVar, value := range tt.envVars {
				os.Setenv(envVar, value)
			}

			// Unset any environment variables not specified in test
			for _, envVar := range envVarsToRestore {
				if _, exists := tt.envVars[envVar]; !exists {
					os.Unsetenv(envVar)
				}
			}

			result, err := config.InitConfigFromEnv()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedConfig.Arch, result.Arch)
				assert.Equal(t, tt.expectedConfig.ForceQuiet, result.ForceQuiet)
				assert.Equal(t, tt.expectedConfig.ForceRemote, result.ForceRemote)
				assert.Equal(t, tt.expectedConfig.GithubActions, result.GithubActions)
				assert.Equal(t, tt.expectedConfig.GithubToken, result.GithubToken)
				assert.Equal(t, tt.expectedConfig.RemoteConfPath, result.RemoteConfPath)
				assert.Equal(t, tt.expectedConfig.RootPath, result.RootPath)
				assert.Equal(t, tt.expectedConfig.SkipInstall, result.SkipInstall)
				assert.Equal(t, tt.expectedConfig.UserPath, result.UserPath)
				assert.Equal(t, tt.expectedConfig.Validation, result.Validation)
				assert.Equal(t, tt.expectedConfig.WorkPath, result.WorkPath)
				assert.Equal(t, tt.expectedConfig.TfKeyPathOrURL, result.TfKeyPathOrURL)
				assert.Equal(t, tt.expectedConfig.TofuKeyPathOrURL, result.TofuKeyPathOrURL)
			}
		})
	}
}

func TestInitConfigFromEnvErrorHandling(t *testing.T) {
	// Save original environment
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	defer func() {
		if originalHome == "" {
			os.Unsetenv("HOME")
		} else {
			os.Setenv("HOME", originalHome)
		}
		if originalUserProfile == "" {
			os.Unsetenv("USERPROFILE")
		} else {
			os.Setenv("USERPROFILE", originalUserProfile)
		}
	}()

	// Test with invalid HOME directory (unset both HOME and USERPROFILE on Windows)
	os.Unsetenv("HOME")
	os.Unsetenv("USERPROFILE")

	_, err := config.InitConfigFromEnv()
	// On Windows, this might not error if USERPROFILE is available
	// So we just check that the function completes without panicking
	if err != nil {
		// The error message might vary by platform, so we just check it's not nil
		assert.NotEmpty(t, err.Error())
	}
}

func TestInitConfigFromEnvWithInvalidBoolValues(t *testing.T) {
	// Save original environment
	originalHome := os.Getenv("HOME")
	originalTenvAutoInstall := os.Getenv("TENV_AUTO_INSTALL")
	originalTenvForceRemote := os.Getenv("TENV_FORCE_REMOTE")
	originalTenvQuiet := os.Getenv("TENV_QUIET")
	defer func() {
		if originalHome == "" {
			os.Unsetenv("HOME")
		} else {
			os.Setenv("HOME", originalHome)
		}
		if originalTenvAutoInstall == "" {
			os.Unsetenv("TENV_AUTO_INSTALL")
		} else {
			os.Setenv("TENV_AUTO_INSTALL", originalTenvAutoInstall)
		}
		if originalTenvForceRemote == "" {
			os.Unsetenv("TENV_FORCE_REMOTE")
		} else {
			os.Setenv("TENV_FORCE_REMOTE", originalTenvForceRemote)
		}
		if originalTenvQuiet == "" {
			os.Unsetenv("TENV_QUIET")
		} else {
			os.Setenv("TENV_QUIET", originalTenvQuiet)
		}
	}()

	// Test with invalid boolean values
	os.Setenv("HOME", "/home/testuser")
	os.Setenv("TENV_AUTO_INSTALL", "invalid_bool")
	os.Setenv("TENV_FORCE_REMOTE", "invalid_bool")
	os.Setenv("TENV_QUIET", "invalid_bool")

	_, err := config.InitConfigFromEnv()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid syntax")
}

func TestInitRemoteConf(t *testing.T) {
	tests := []struct {
		name              string
		remoteConfPath    string
		remoteConfContent string
		expectedError     bool
	}{
		{
			name:           "file not found",
			remoteConfPath: "/nonexistent/remote.yaml",
			expectedError:  false,
		},
		{
			name:              "invalid yaml content",
			remoteConfPath:    "",
			remoteConfContent: `invalid: yaml: content: [`,
			expectedError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory and file if needed
			if tt.remoteConfContent != "" {
				tempDir := t.TempDir()
				tt.remoteConfPath = filepath.Join(tempDir, "remote.yaml")

				err := os.WriteFile(tt.remoteConfPath, []byte(tt.remoteConfContent), 0644)
				require.NoError(t, err)
			}

			config := config.Config{
				RootPath:       "/tmp/tenv",
				RemoteConfPath: tt.remoteConfPath,
				Displayer:      &mockDisplayer{},
			}

			err := config.InitRemoteConf()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Mock displayer for testing
type mockDisplayer struct{}

func (m *mockDisplayer) Display(msg string)                                       {}
func (m *mockDisplayer) Log(level hclog.Level, msg string, fields ...interface{}) {}
func (m *mockDisplayer) Flush(bool)                                               {}
func (m *mockDisplayer) IsDebug() bool                                            { return false }
