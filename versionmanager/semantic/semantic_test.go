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

package semantic

import (
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/go-version"

	"os"
	"path/filepath"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/versionmanager/lastuse"
	"github.com/tofuutils/tenv/v4/versionmanager/semantic/types"
	"github.com/tofuutils/tenv/v4/versionmanager/version"
)

func TestCmpVersion(t *testing.T) {
	t.Parallel()

	versions := []string{"1.6.0-beta5", "1.5.2", "1.6.0-alpha5", "1.6.0", "1.5.1", "1.5.0", "1.6.0-rc1"}
	slices.SortFunc(versions, CmpVersion)
	if !slices.Equal(versions, []string{"1.5.0", "1.5.1", "1.5.2", "1.6.0-alpha5", "1.6.0-beta5", "1.6.0-rc1", "1.6.0"}) {
		t.Error("Unmatching results, get :", versions)
	}
}

func TestStableVersion(t *testing.T) {
	t.Parallel()

	var filtered []string
	for _, version := range []string{"1.5.0", "1.5.1", "1.5.2", "1.6.0-alpha5", "1.6.0-beta5", "1.6.0-rc1", "1.6.0"} {
		if StableVersion(version) {
			filtered = append(filtered, version)
		}
	}

	if !slices.Equal(filtered, []string{"1.5.0", "1.5.1", "1.5.2", "1.6.0"}) {
		t.Error("Unmatching results, get :", filtered)
	}
}

type mockConstraintInfo struct {
	constraint string
}

func (m *mockConstraintInfo) ReadDefaultConstraint() string {
	return m.constraint
}

func TestParsePredicate(t *testing.T) {
	tests := []struct {
		name       string
		constraint string
		want       types.PredicateInfo
		wantErr    bool
	}{
		{
			name:       "valid constraint",
			constraint: ">= 1.0.0",
			want: types.PredicateInfo{
				Predicate: func(v string) bool {
					ver, _ := version.NewVersion(v)
					return ver.Compare(version.Must(version.NewVersion("1.0.0"))) >= 0
				},
				ReverseOrder: false,
			},
			wantErr: false,
		},
		{
			name:       "invalid constraint",
			constraint: "invalid",
			want:       types.PredicateInfo{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePredicate(tt.constraint)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePredicate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Test with a sample version
				testVer := "1.1.0"
				if got.Predicate(testVer) != tt.want.Predicate(testVer) {
					t.Errorf("ParsePredicate() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestAddDefaultConstraint(t *testing.T) {
	tests := []struct {
		name       string
		constraint string
		want       string
	}{
		{
			name:       "empty constraint",
			constraint: "",
			want:       ">= 0.0.0",
		},
		{
			name:       "existing constraint",
			constraint: ">= 1.0.0",
			want:       ">= 1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &mockConstraintInfo{constraint: tt.constraint}
			if got := AddDefaultConstraint(info); got != tt.want {
				t.Errorf("AddDefaultConstraint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPredicateFromConstraint(t *testing.T) {
	tests := []struct {
		name       string
		constraint string
		want       types.PredicateInfo
		wantErr    bool
	}{
		{
			name:       "valid constraint",
			constraint: ">= 1.0.0",
			want: types.PredicateInfo{
				Predicate: func(v string) bool {
					ver, _ := version.NewVersion(v)
					return ver.Compare(version.Must(version.NewVersion("1.0.0"))) >= 0
				},
				ReverseOrder: false,
			},
			wantErr: false,
		},
		{
			name:       "invalid constraint",
			constraint: "invalid",
			want:       types.PredicateInfo{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PredicateFromConstraint(tt.constraint)
			if (err != nil) != tt.wantErr {
				t.Errorf("PredicateFromConstraint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Test with a sample version
				testVer := "1.1.0"
				if got.Predicate(testVer) != tt.want.Predicate(testVer) {
					t.Errorf("PredicateFromConstraint() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func mustNewConstraint(t *testing.T, constraint string) version.Constraints {
	t.Helper()
	c, err := version.NewConstraint(constraint)
	if err != nil {
		t.Fatalf("Failed to create constraint: %v", err)
	}
	return c
}

func TestSelectVersionsToUninstall(t *testing.T) {
	t.Parallel()

	testVersions := []string{"1.6.0", "1.5.2", "1.5.1", "1.5.0"}
	testPath := "/test/path"
	testConfig := &config.Config{}

	tests := []struct {
		name      string
		behaviour string
		versions  []string
		want      []string
		wantErr   bool
	}{
		{
			name:      "all versions",
			behaviour: allKey,
			versions:  testVersions,
			want:      testVersions,
			wantErr:   false,
		},
		{
			name:      "but last version",
			behaviour: butLast,
			versions:  testVersions,
			want:      []string{"1.5.2", "1.5.1", "1.5.0"},
			wantErr:   false,
		},
		{
			name:      "but last with empty list",
			behaviour: butLast,
			versions:  []string{},
			want:      nil,
			wantErr:   false,
		},
		{
			name:      "version constraint",
			behaviour: "< 1.5.2",
			versions:  testVersions,
			want:      []string{"1.5.1", "1.5.0"},
			wantErr:   false,
		},
		{
			name:      "invalid version constraint",
			behaviour: "invalid",
			versions:  testVersions,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "not used for days - invalid format",
			behaviour: "not-used-for:abc",
			versions:  testVersions,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "not used for days - valid format",
			behaviour: "not-used-for:30d",
			versions:  testVersions,
			want:      []string{}, // Empty because no files exist in test path
			wantErr:   false,
		},
		{
			name:      "not used for months - valid format",
			behaviour: "not-used-for:2m",
			versions:  testVersions,
			want:      []string{}, // Empty because no files exist in test path
			wantErr:   false,
		},
		{
			name:      "not used since - invalid date",
			behaviour: "not-used-since:invalid",
			versions:  testVersions,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "not used since - valid date",
			behaviour: "not-used-since:2024-01-01",
			versions:  testVersions,
			want:      []string{}, // Empty because no files exist in test path
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := SelectVersionsToUninstall(tt.behaviour, testPath, tt.versions, testConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectVersionsToUninstall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("SelectVersionsToUninstall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []string
		pred     func(string) bool
		expected []string
	}{
		{
			name:     "filter even numbers",
			input:    []string{"1", "2", "3", "4", "5"},
			pred:     func(s string) bool { n, _ := strconv.Atoi(s); return n%2 == 0 },
			expected: []string{"2", "4"},
		},
		{
			name:     "empty input",
			input:    []string{},
			pred:     func(s string) bool { return true },
			expected: []string{},
		},
		{
			name:     "no matches",
			input:    []string{"1", "3", "5"},
			pred:     func(s string) bool { n, _ := strconv.Atoi(s); return n%2 == 0 },
			expected: []string{},
		},
		{
			name:     "all matches",
			input:    []string{"2", "4", "6"},
			pred:     func(s string) bool { n, _ := strconv.Atoi(s); return n%2 == 0 },
			expected: []string{"2", "4", "6"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := filterStrings(tt.input, tt.pred)
			if !slices.Equal(result, tt.expected) {
				t.Errorf("filterStrings() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPredicateBeforeDate(t *testing.T) {
	t.Parallel()

	testPath := t.TempDir()
	testVersion := "1.0.0"
	versionPath := filepath.Join(testPath, testVersion)
	if err := os.MkdirAll(versionPath, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a last use file with a known date
	testDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	lastuse.Write(versionPath, testDate, &config.Config{})

	tests := []struct {
		name       string
		beforeDate time.Time
		want       bool
	}{
		{
			name:       "date before last use",
			beforeDate: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
			want:       false,
		},
		{
			name:       "date after last use",
			beforeDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			want:       true,
		},
		{
			name:       "same date as last use",
			beforeDate: testDate,
			want:       false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pred := predicateBeforeDate(testPath, tt.beforeDate, &config.Config{})
			if got := pred(testVersion); got != tt.want {
				t.Errorf("predicateBeforeDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
