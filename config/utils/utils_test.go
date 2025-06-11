package configutils

import (
	"testing"
)

func TestGetenvFunc_Bool(t *testing.T) {
	tests := []struct {
		name         string
		getenv       GetenvFunc
		defaultValue bool
		key          string
		want         bool
		wantErr      bool
	}{
		{
			name:         "true value",
			getenv:       func(string) string { return "true" },
			defaultValue: false,
			key:          "TEST_KEY",
			want:         true,
			wantErr:      false,
		},
		{
			name:         "false value",
			getenv:       func(string) string { return "false" },
			defaultValue: true,
			key:          "TEST_KEY",
			want:         false,
			wantErr:      false,
		},
		{
			name:         "empty value",
			getenv:       func(string) string { return "" },
			defaultValue: true,
			key:          "TEST_KEY",
			want:         true,
			wantErr:      false,
		},
		{
			name:         "invalid value",
			getenv:       func(string) string { return "invalid" },
			defaultValue: false,
			key:          "TEST_KEY",
			want:         false,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.getenv.Bool(tt.defaultValue, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetenvFunc.Bool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetenvFunc.Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetenvFunc_BoolFallback(t *testing.T) {
	tests := []struct {
		name         string
		getenv       GetenvFunc
		defaultValue bool
		keys         []string
		want         bool
		wantErr      bool
	}{
		{
			name: "first key has value",
			getenv: func(key string) string {
				if key == "KEY1" {
					return "true"
				}
				return ""
			},
			defaultValue: false,
			keys:         []string{"KEY1", "KEY2"},
			want:         true,
			wantErr:      false,
		},
		{
			name: "second key has value",
			getenv: func(key string) string {
				if key == "KEY2" {
					return "true"
				}
				return ""
			},
			defaultValue: false,
			keys:         []string{"KEY1", "KEY2"},
			want:         true,
			wantErr:      false,
		},
		{
			name:         "no keys have value",
			getenv:       func(string) string { return "" },
			defaultValue: true,
			keys:         []string{"KEY1", "KEY2"},
			want:         true,
			wantErr:      false,
		},
		{
			name: "invalid value",
			getenv: func(key string) string {
				if key == "KEY1" {
					return "invalid"
				}
				return ""
			},
			defaultValue: false,
			keys:         []string{"KEY1", "KEY2"},
			want:         false,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.getenv.BoolFallback(tt.defaultValue, tt.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetenvFunc.BoolFallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetenvFunc.BoolFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetenvFunc_Fallback(t *testing.T) {
	tests := []struct {
		name   string
		getenv GetenvFunc
		keys   []string
		want   string
	}{
		{
			name: "first key has value",
			getenv: func(key string) string {
				if key == "KEY1" {
					return "value1"
				}
				return ""
			},
			keys: []string{"KEY1", "KEY2"},
			want: "value1",
		},
		{
			name: "second key has value",
			getenv: func(key string) string {
				if key == "KEY2" {
					return "value2"
				}
				return ""
			},
			keys: []string{"KEY1", "KEY2"},
			want: "value2",
		},
		{
			name:   "no keys have value",
			getenv: func(string) string { return "" },
			keys:   []string{"KEY1", "KEY2"},
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.getenv.Fallback(tt.keys...); got != tt.want {
				t.Errorf("GetenvFunc.Fallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetenvFunc_Present(t *testing.T) {
	tests := []struct {
		name   string
		getenv GetenvFunc
		key    string
		want   bool
	}{
		{
			name:   "key present",
			getenv: func(string) string { return "value" },
			key:    "TEST_KEY",
			want:   true,
		},
		{
			name:   "key not present",
			getenv: func(string) string { return "" },
			key:    "TEST_KEY",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.getenv.Present(tt.key); got != tt.want {
				t.Errorf("GetenvFunc.Present() = %v, want %v", got, tt.want)
			}
		})
	}
}
