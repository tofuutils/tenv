package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	conf, err := DefaultConfig()
	if err != nil {
		t.Fatalf("DefaultConfig() error = %v", err)
	}

	// Verify default values
	if conf.Arch != runtime.GOARCH {
		t.Errorf("DefaultConfig().Arch = %v, want %v", conf.Arch, runtime.GOARCH)
	}
	if conf.SkipInstall != true {
		t.Errorf("DefaultConfig().SkipInstall = %v, want true", conf.SkipInstall)
	}
	if conf.remoteConfLoaded != true {
		t.Errorf("DefaultConfig().remoteConfLoaded = %v, want true", conf.remoteConfLoaded)
	}
	if conf.WorkPath != "." {
		t.Errorf("DefaultConfig().WorkPath = %v, want .", conf.WorkPath)
	}

	// Verify paths
	userPath, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() error = %v", err)
	}
	if conf.UserPath != userPath {
		t.Errorf("DefaultConfig().UserPath = %v, want %v", conf.UserPath, userPath)
	}
	if conf.RootPath != filepath.Join(userPath, defaultDirName) {
		t.Errorf("DefaultConfig().RootPath = %v, want %v", conf.RootPath, filepath.Join(userPath, defaultDirName))
	}
}

func TestInitConfigFromEnv(t *testing.T) {
	// Set up test environment variables
	os.Setenv("TENV_ARCH", "test-arch")
	os.Setenv("TENV_AUTO_INSTALL", "true")
	os.Setenv("TENV_FORCE_QUIET", "true")
	os.Setenv("TENV_FORCE_REMOTE", "true")
	os.Setenv("TENV_GITHUB_ACTIONS", "true")
	os.Setenv("TENV_GITHUB_TOKEN", "test-token")
	os.Setenv("TENV_REMOTE_CONF_PATH", "/test/path")
	os.Setenv("TENV_ROOT_PATH", "/test/root")
	os.Setenv("TENV_SKIP_INSTALL", "true")
	os.Setenv("TENV_SKIP_SIGNATURE", "true")
	os.Setenv("TENV_USER_PATH", "/test/user")
	os.Setenv("TENV_WORK_PATH", "/test/work")

	conf, err := InitConfigFromEnv()
	if err != nil {
		t.Fatalf("InitConfigFromEnv() error = %v", err)
	}

	// Verify environment variable values
	if conf.Arch != "test-arch" {
		t.Errorf("InitConfigFromEnv().Arch = %v, want test-arch", conf.Arch)
	}
	if conf.SkipInstall != true {
		t.Errorf("InitConfigFromEnv().SkipInstall = %v, want true", conf.SkipInstall)
	}
	if conf.ForceQuiet != true {
		t.Errorf("InitConfigFromEnv().ForceQuiet = %v, want true", conf.ForceQuiet)
	}
	if conf.ForceRemote != true {
		t.Errorf("InitConfigFromEnv().ForceRemote = %v, want true", conf.ForceRemote)
	}
	if conf.GithubActions != true {
		t.Errorf("InitConfigFromEnv().GithubActions = %v, want true", conf.GithubActions)
	}
	if conf.GithubToken != "test-token" {
		t.Errorf("InitConfigFromEnv().GithubToken = %v, want test-token", conf.GithubToken)
	}
	if conf.RemoteConfPath != "/test/path" {
		t.Errorf("InitConfigFromEnv().RemoteConfPath = %v, want /test/path", conf.RemoteConfPath)
	}
	if conf.RootPath != "/test/root" {
		t.Errorf("InitConfigFromEnv().RootPath = %v, want /test/root", conf.RootPath)
	}
	if conf.SkipSignature != true {
		t.Errorf("InitConfigFromEnv().SkipSignature = %v, want true", conf.SkipSignature)
	}
	if conf.UserPath != "/test/user" {
		t.Errorf("InitConfigFromEnv().UserPath = %v, want /test/user", conf.UserPath)
	}
	if conf.WorkPath != "/test/work" {
		t.Errorf("InitConfigFromEnv().WorkPath = %v, want /test/work", conf.WorkPath)
	}

	// Clean up environment variables
	os.Unsetenv("TENV_ARCH")
	os.Unsetenv("TENV_AUTO_INSTALL")
	os.Unsetenv("TENV_FORCE_QUIET")
	os.Unsetenv("TENV_FORCE_REMOTE")
	os.Unsetenv("TENV_GITHUB_ACTIONS")
	os.Unsetenv("TENV_GITHUB_TOKEN")
	os.Unsetenv("TENV_REMOTE_CONF_PATH")
	os.Unsetenv("TENV_ROOT_PATH")
	os.Unsetenv("TENV_SKIP_INSTALL")
	os.Unsetenv("TENV_SKIP_SIGNATURE")
	os.Unsetenv("TENV_USER_PATH")
	os.Unsetenv("TENV_WORK_PATH")
}

func TestInitDisplayer(t *testing.T) {
	tests := []struct {
		name       string
		proxyCall  bool
		forceQuiet bool
		verbose    bool
	}{
		{
			name:       "proxy call with verbose",
			proxyCall:  true,
			forceQuiet: false,
			verbose:    true,
		},
		{
			name:       "proxy call with quiet",
			proxyCall:  true,
			forceQuiet: true,
			verbose:    false,
		},
		{
			name:       "direct call with verbose",
			proxyCall:  false,
			forceQuiet: false,
			verbose:    true,
		},
		{
			name:       "direct call with quiet",
			proxyCall:  false,
			forceQuiet: true,
			verbose:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := Config{
				ForceQuiet:     tt.forceQuiet,
				DisplayVerbose: tt.verbose,
			}
			conf.InitDisplayer(tt.proxyCall)
			if conf.Displayer == nil {
				t.Error("InitDisplayer() did not initialize Displayer")
			}
		})
	}
}

func TestInitInstall(t *testing.T) {
	tests := []struct {
		name           string
		forceInstall   bool
		forceNoInstall bool
		initialSkip    bool
		expectedSkip   bool
	}{
		{
			name:           "force install",
			forceInstall:   true,
			forceNoInstall: false,
			initialSkip:    true,
			expectedSkip:   false,
		},
		{
			name:           "force no install",
			forceInstall:   false,
			forceNoInstall: true,
			initialSkip:    false,
			expectedSkip:   true,
		},
		{
			name:           "no force",
			forceInstall:   false,
			forceNoInstall: false,
			initialSkip:    true,
			expectedSkip:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := Config{
				SkipInstall: tt.initialSkip,
			}
			conf.InitInstall(tt.forceInstall, tt.forceNoInstall)
			if conf.SkipInstall != tt.expectedSkip {
				t.Errorf("InitInstall() SkipInstall = %v, want %v", conf.SkipInstall, tt.expectedSkip)
			}
		})
	}
}

func TestEmptyGetenv(t *testing.T) {
	if got := EmptyGetenv("TEST_KEY"); got != "" {
		t.Errorf("EmptyGetenv() = %v, want empty string", got)
	}
}
