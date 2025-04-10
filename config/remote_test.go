package config

import (
	"testing"

	configutils "github.com/tofuutils/tenv/v4/config/utils"
	"github.com/tofuutils/tenv/v4/pkg/download"
)

func TestMakeDefaultRemoteConfig(t *testing.T) {
	defaultURL := "https://example.com"
	defaultBaseURL := "https://base.example.com"

	config := makeDefaultRemoteConfig(defaultURL, defaultBaseURL)

	if config.defaultURL != defaultURL {
		t.Errorf("makeDefaultRemoteConfig().defaultURL = %v, want %v", config.defaultURL, defaultURL)
	}
	if config.defaultBaseURL != defaultBaseURL {
		t.Errorf("makeDefaultRemoteConfig().defaultBaseURL = %v, want %v", config.defaultBaseURL, defaultBaseURL)
	}
}

func TestMakeRemoteConfig(t *testing.T) {
	tests := []struct {
		name            string
		getenv          configutils.GetenvFunc
		remoteURLEnv    string
		listURLEnv      string
		installModeEnv  string
		listModeEnv     string
		defaultURL      string
		defaultBaseURL  string
		wantRemoteURL   string
		wantListURL     string
		wantInstallMode string
		wantListMode    string
	}{
		{
			name:            "all defaults",
			getenv:          configutils.EmptyGetenv,
			remoteURLEnv:    "TEST_REMOTE_URL",
			listURLEnv:      "TEST_LIST_URL",
			installModeEnv:  "TEST_INSTALL_MODE",
			listModeEnv:     "TEST_LIST_MODE",
			defaultURL:      "https://example.com",
			defaultBaseURL:  "https://base.example.com",
			wantRemoteURL:   "https://example.com",
			wantListURL:     "https://example.com",
			wantInstallMode: InstallModeDirect,
			wantListMode:    ListModeHTML,
		},
		{
			name: "all custom",
			getenv: func(key string) string {
				switch key {
				case "TEST_REMOTE_URL":
					return "https://custom.example.com"
				case "TEST_LIST_URL":
					return "https://list.example.com"
				case "TEST_INSTALL_MODE":
					return "custom"
				case "TEST_LIST_MODE":
					return "custom"
				}
				return ""
			},
			remoteURLEnv:    "TEST_REMOTE_URL",
			listURLEnv:      "TEST_LIST_URL",
			installModeEnv:  "TEST_INSTALL_MODE",
			listModeEnv:     "TEST_LIST_MODE",
			defaultURL:      "https://example.com",
			defaultBaseURL:  "https://base.example.com",
			wantRemoteURL:   "https://custom.example.com",
			wantListURL:     "https://list.example.com",
			wantInstallMode: "custom",
			wantListMode:    "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := makeRemoteConfig(tt.getenv, tt.remoteURLEnv, tt.listURLEnv, tt.installModeEnv, tt.listModeEnv, tt.defaultURL, tt.defaultBaseURL)

			if config.GetRemoteURL() != tt.wantRemoteURL {
				t.Errorf("makeRemoteConfig().GetRemoteURL() = %v, want %v", config.GetRemoteURL(), tt.wantRemoteURL)
			}
			if config.GetListURL() != tt.wantListURL {
				t.Errorf("makeRemoteConfig().GetListURL() = %v, want %v", config.GetListURL(), tt.wantListURL)
			}
			if config.GetInstallMode() != tt.wantInstallMode {
				t.Errorf("makeRemoteConfig().GetInstallMode() = %v, want %v", config.GetInstallMode(), tt.wantInstallMode)
			}
			if config.GetListMode() != tt.wantListMode {
				t.Errorf("makeRemoteConfig().GetListMode() = %v, want %v", config.GetListMode(), tt.wantListMode)
			}
		})
	}
}

func TestRemoteConfig_GetRewriteRule(t *testing.T) {
	tests := []struct {
		name     string
		config   RemoteConfig
		wantRule download.URLTransformer
	}{
		{
			name: "no rewrite rule",
			config: RemoteConfig{
				Data: map[string]string{},
			},
			wantRule: nil,
		},
		{
			name: "with rewrite rule",
			config: RemoteConfig{
				Data: map[string]string{
					"rewrite_rule": "s/old/new/",
				},
			},
			wantRule: download.URLTransformer(func(url string) string {
				return "new"
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := tt.config.GetRewriteRule()
			if tt.wantRule == nil && rule != nil {
				t.Error("GetRewriteRule() returned non-nil rule when expected nil")
			}
			if tt.wantRule != nil && rule == nil {
				t.Error("GetRewriteRule() returned nil rule when expected non-nil")
			}
		})
	}
}

func TestMapGetDefault(t *testing.T) {
	tests := []struct {
		name         string
		m            map[string]string
		key          string
		defaultValue string
		want         string
	}{
		{
			name:         "key exists",
			m:            map[string]string{"test": "value"},
			key:          "test",
			defaultValue: "default",
			want:         "value",
		},
		{
			name:         "key does not exist",
			m:            map[string]string{},
			key:          "test",
			defaultValue: "default",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapGetDefault(tt.m, tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("MapGetDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBasicAuthOption(t *testing.T) {
	tests := []struct {
		name       string
		getenv     configutils.GetenvFunc
		userEnv    string
		passEnv    string
		wantOption bool
	}{
		{
			name:       "no auth",
			getenv:     configutils.EmptyGetenv,
			userEnv:    "TEST_USER",
			passEnv:    "TEST_PASS",
			wantOption: false,
		},
		{
			name: "with auth",
			getenv: func(key string) string {
				switch key {
				case "TEST_USER":
					return "user"
				case "TEST_PASS":
					return "pass"
				}
				return ""
			},
			userEnv:    "TEST_USER",
			passEnv:    "TEST_PASS",
			wantOption: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := GetBasicAuthOption(tt.getenv, tt.userEnv, tt.passEnv)
			if tt.wantOption && len(options) == 0 {
				t.Error("GetBasicAuthOption() returned no options when expected some")
			}
			if !tt.wantOption && len(options) > 0 {
				t.Error("GetBasicAuthOption() returned options when expected none")
			}
		})
	}
}
