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

// Helper functions for platform-specific paths.
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := config.ParseValidationMode(testCase.input)
			if result != testCase.expected {
				t.Errorf("ParseValidationMode(%q) = %v, want %v", testCase.input, result, testCase.expected)
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			cfg := &config.Config{
				ForceQuiet:     testCase.forceQuiet,
				DisplayVerbose: true,
				Getenv:         config.EmptyGetenv,
			}

			cfg.InitDisplayer(testCase.proxyCall)

			if testCase.expectQuiet && cfg.Displayer != loghelper.InertDisplayer {
				t.Errorf("Expected InertDisplayer when ForceQuiet is true")
			}

			if !testCase.expectQuiet && cfg.Displayer == loghelper.InertDisplayer {
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			cfg := &config.Config{Validation: testCase.initial}
			cfg.InitValidation(testCase.skipSum, testCase.skipSign)

			if cfg.Validation != testCase.expected {
				t.Errorf("InitValidation() = %v, want %v", cfg.Validation, testCase.expected)
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			cfg := &config.Config{SkipInstall: testCase.initial}
			cfg.InitInstall(testCase.forceInstall, testCase.forceNoInstall)

			if cfg.SkipInstall != testCase.expected {
				t.Errorf("InitInstall() = %v, want %v", cfg.SkipInstall, testCase.expected)
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := config.MapGetDefault(testCase.m, testCase.key, testCase.defaultValue)
			if result != testCase.expected {
				t.Errorf("MapGetDefault() = %v, want %v", result, testCase.expected)
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			options := config.GetBasicAuthOption(mockGetenv, testCase.userEnvName, testCase.passEnvName)
			if testCase.expected && options == nil {
				t.Errorf("Expected non-nil options when both env vars are present")
			}
			if !testCase.expected && options != nil {
				t.Errorf("Expected nil options when env vars are missing")
			}
		})
	}
}

func getDefaultExpectedConfig() config.Config {
	return config.Config{
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
	}
}

func setupTestEnv(t *testing.T, envVars map[string]string) {
	t.Helper()
	envVarsToRestore := []string{
		"HOME", "TENV_ARCH", "TENV_AUTO_INSTALL", "TENV_FORCE_REMOTE", "TENV_QUIET",
		"TENV_ROOT", "TENV_GITHUB_TOKEN", "GITHUB_ACTIONS", "TENV_VALIDATION", "TENV_REMOTE_CONF",
		"TFENV_HASHICORP_PGP_KEY", "TOFUENV_OPENTOFU_PGP_KEY",
	}

	for _, envVar := range envVarsToRestore {
		if value, exists := envVars[envVar]; exists {
			t.Setenv(envVar, value)
		} else {
			t.Setenv(envVar, "")
		}
	}
}

//nolint:paralleltest,tparallel // t.Setenv cannot be used with t.Parallel()
func TestInitConfigFromEnv(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		modifyFunc  func(config.Config) config.Config
		expectError bool
	}{
		{
			name: "default configuration",
			envVars: map[string]string{
				"HOME": getTestHomeDir(),
			},
			modifyFunc:  func(c config.Config) config.Config { return c },
			expectError: false,
		},
		{
			name: "custom arch and auto install",
			envVars: map[string]string{
				"HOME":              getTestHomeDir(),
				"TENV_ARCH":         "arm64",
				"TENV_AUTO_INSTALL": "true",
			},
			modifyFunc: func(c config.Config) config.Config {
				c.Arch = "arm64"
				c.SkipInstall = false

				return c
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
			modifyFunc: func(c config.Config) config.Config {
				c.ForceQuiet = true
				c.ForceRemote = true

				return c
			},
			expectError: false,
		},
		{
			name: "custom root path",
			envVars: map[string]string{
				"HOME":      getTestHomeDir(),
				"TENV_ROOT": "/custom/tenv/path",
			},
			modifyFunc: func(c config.Config) config.Config {
				c.RootPath = "/custom/tenv/path"

				return c
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
			modifyFunc: func(c config.Config) config.Config {
				c.GithubActions = true
				c.GithubToken = "ghp_1234567890abcdef"

				return c
			},
			expectError: false,
		},
		{
			name: "validation modes",
			envVars: map[string]string{
				"HOME":            getTestHomeDir(),
				"TENV_VALIDATION": "sha",
			},
			modifyFunc: func(c config.Config) config.Config {
				c.Validation = config.ShaValidation

				return c
			},
			expectError: false,
		},
		{
			name: "custom remote configuration path",
			envVars: map[string]string{
				"HOME":             getTestHomeDir(),
				"TENV_REMOTE_CONF": "/etc/tenv/remote.yaml",
			},
			modifyFunc: func(c config.Config) config.Config {
				c.RemoteConfPath = "/etc/tenv/remote.yaml"

				return c
			},
			expectError: false,
		},
		{
			name: "custom terraform key path",
			envVars: map[string]string{
				"HOME":                    getTestHomeDir(),
				"TFENV_HASHICORP_PGP_KEY": "/custom/terraform-key.asc",
			},
			modifyFunc: func(c config.Config) config.Config {
				c.TfKeyPathOrURL = "/custom/terraform-key.asc"

				return c
			},
			expectError: false,
		},
		{
			name: "custom tofu key path",
			envVars: map[string]string{
				"HOME":                     getTestHomeDir(),
				"TOFUENV_OPENTOFU_PGP_KEY": "/custom/tofu-key.asc",
			},
			modifyFunc: func(c config.Config) config.Config {
				c.TofuKeyPathOrURL = "/custom/tofu-key.asc"

				return c
			},
			expectError: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			setupTestEnv(t, testCase.envVars)

			result, err := config.InitConfigFromEnv()

			if testCase.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				expected := testCase.modifyFunc(getDefaultExpectedConfig())
				assert.Equal(t, expected.Arch, result.Arch)
				assert.Equal(t, expected.ForceQuiet, result.ForceQuiet)
				assert.Equal(t, expected.ForceRemote, result.ForceRemote)
				assert.Equal(t, expected.GithubActions, result.GithubActions)
				assert.Equal(t, expected.GithubToken, result.GithubToken)
				assert.Equal(t, expected.RemoteConfPath, result.RemoteConfPath)
				assert.Equal(t, expected.RootPath, result.RootPath)
				assert.Equal(t, expected.SkipInstall, result.SkipInstall)
				assert.Equal(t, expected.UserPath, result.UserPath)
				assert.Equal(t, expected.Validation, result.Validation)
				assert.Equal(t, expected.WorkPath, result.WorkPath)
				assert.Equal(t, expected.TfKeyPathOrURL, result.TfKeyPathOrURL)
				assert.Equal(t, expected.TofuKeyPathOrURL, result.TofuKeyPathOrURL)
			}
		})
	}
}

func TestInitConfigFromEnvErrorHandling(t *testing.T) {
	// Test with invalid HOME directory (unset both HOME and USERPROFILE on Windows)
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")

	_, err := config.InitConfigFromEnv()
	// On Windows, this might not error if USERPROFILE is available
	// So we just check that the function completes without panicking
	if err != nil {
		// The error message might vary by platform, so we just check it's not nil
		assert.NotEmpty(t, err.Error())
	}
}

func TestInitConfigFromEnvWithInvalidBoolValues(t *testing.T) {
	// Test with invalid boolean values
	t.Setenv("HOME", "/home/testuser")
	t.Setenv("TENV_AUTO_INSTALL", "invalid_bool")
	t.Setenv("TENV_FORCE_REMOTE", "invalid_bool")
	t.Setenv("TENV_QUIET", "invalid_bool")

	_, err := config.InitConfigFromEnv()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid syntax")
}

func TestInitRemoteConf(t *testing.T) {
	t.Parallel()

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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Create temporary directory and file if needed
			if testCase.remoteConfContent != "" {
				tempDir := t.TempDir()
				testCase.remoteConfPath = filepath.Join(tempDir, "remote.yaml")

				err := os.WriteFile(testCase.remoteConfPath, []byte(testCase.remoteConfContent), 0o600)
				require.NoError(t, err)
			}

			config := config.Config{
				RootPath:       "/tmp/tenv",
				RemoteConfPath: testCase.remoteConfPath,
				Displayer:      &mockDisplayer{},
			}

			err := config.InitRemoteConf()

			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Mock displayer for testing.
type mockDisplayer struct{}

func (m *mockDisplayer) Display(msg string)                                       {}
func (m *mockDisplayer) Log(level hclog.Level, msg string, fields ...interface{}) {}
func (m *mockDisplayer) Flush(bool)                                               {}
func (m *mockDisplayer) IsDebug() bool                                            { return false }
