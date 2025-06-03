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

package reversecmp_test

import (
	"cmp"
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/reversecmp"
)

func TestReverserFalse(t *testing.T) {
	t.Parallel()

	reversed := reversecmp.Reverser[int](cmp.Compare[int], false)
	if reversed(0, 5) != -1 {
		t.Error("Not ordered")
	}
	if reversed(1, 1) != 0 {
		t.Error("WTF")
	}
	if reversed(10, 5) != 1 {
		t.Error("Not ordered again")
	}
}

func TestReverserTrue(t *testing.T) {
	t.Parallel()

	reversed := reversecmp.Reverser[int](cmp.Compare[int], true)
	if reversed(0, 5) != 1 {
		t.Error("Not inversed")
	}
	if reversed(1, 1) != 0 {
		t.Error("WTF")
	}
	if reversed(10, 5) != -1 {
		t.Error("Not inversed again")
	}
}

func TestReverserString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		a            string
		b            string
		reverseOrder bool
		want         int
	}{
		{
			name:         "forward empty strings",
			a:            "",
			b:            "",
			reverseOrder: false,
			want:         0,
		},
		{
			name:         "forward lexicographic order",
			a:            "abc",
			b:            "def",
			reverseOrder: false,
			want:         -1,
		},
		{
			name:         "reverse lexicographic order",
			a:            "abc",
			b:            "def",
			reverseOrder: true,
			want:         1,
		},
		{
			name:         "forward with unicode",
			a:            "世界",
			b:            "你好",
			reverseOrder: false,
			want:         1,
		},
		{
			name:         "reverse with unicode",
			a:            "世界",
			b:            "你好",
			reverseOrder: true,
			want:         -1,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reversed := reversecmp.Reverser[string](cmp.Compare[string], tt.reverseOrder)
			got := reversed(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Reverser(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestReverserFloat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		a            float64
		b            float64
		reverseOrder bool
		want         int
	}{
		{
			name:         "forward equal floats",
			a:            0.0,
			b:            0.0,
			reverseOrder: false,
			want:         0,
		},
		{
			name:         "forward positive floats",
			a:            1.5,
			b:            2.5,
			reverseOrder: false,
			want:         -1,
		},
		{
			name:         "reverse positive floats",
			a:            1.5,
			b:            2.5,
			reverseOrder: true,
			want:         1,
		},
		{
			name:         "forward negative floats",
			a:            -2.5,
			b:            -1.5,
			reverseOrder: false,
			want:         -1,
		},
		{
			name:         "reverse negative floats",
			a:            -2.5,
			b:            -1.5,
			reverseOrder: true,
			want:         1,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reversed := reversecmp.Reverser[float64](cmp.Compare[float64], tt.reverseOrder)
			got := reversed(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Reverser(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// Custom type to test with custom comparison function
type Version struct {
	Major, Minor, Patch int
}

func compareVersion(a, b Version) int {
	if a.Major != b.Major {
		return cmp.Compare(a.Major, b.Major)
	}
	if a.Minor != b.Minor {
		return cmp.Compare(a.Minor, b.Minor)
	}
	return cmp.Compare(a.Patch, b.Patch)
}

func TestReverserCustomType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		a            Version
		b            Version
		reverseOrder bool
		want         int
	}{
		{
			name:         "forward equal versions",
			a:            Version{1, 0, 0},
			b:            Version{1, 0, 0},
			reverseOrder: false,
			want:         0,
		},
		{
			name:         "forward major version diff",
			a:            Version{1, 0, 0},
			b:            Version{2, 0, 0},
			reverseOrder: false,
			want:         -1,
		},
		{
			name:         "reverse major version diff",
			a:            Version{1, 0, 0},
			b:            Version{2, 0, 0},
			reverseOrder: true,
			want:         1,
		},
		{
			name:         "forward minor version diff",
			a:            Version{1, 1, 0},
			b:            Version{1, 2, 0},
			reverseOrder: false,
			want:         -1,
		},
		{
			name:         "forward patch version diff",
			a:            Version{1, 0, 1},
			b:            Version{1, 0, 2},
			reverseOrder: false,
			want:         -1,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reversed := reversecmp.Reverser[Version](compareVersion, tt.reverseOrder)
			got := reversed(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Reverser(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestReverserInt8(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		a            int8
		b            int8
		reverseOrder bool
		want         int
	}{
		{
			name:         "forward equal values",
			a:            0,
			b:            0,
			reverseOrder: false,
			want:         0,
		},
		{
			name:         "forward a < b",
			a:            1,
			b:            2,
			reverseOrder: false,
			want:         -1,
		},
		{
			name:         "reverse a < b",
			a:            1,
			b:            2,
			reverseOrder: true,
			want:         1,
		},
		{
			name:         "forward a > b",
			a:            2,
			b:            1,
			reverseOrder: false,
			want:         1,
		},
		{
			name:         "reverse a > b",
			a:            2,
			b:            1,
			reverseOrder: true,
			want:         -1,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reversed := reversecmp.Reverser[int8](cmp.Compare[int8], tt.reverseOrder)
			got := reversed(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Reverser(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestReverserUint16(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		a            uint16
		b            uint16
		reverseOrder bool
		want         int
	}{
		{
			name:         "forward equal values",
			a:            0,
			b:            0,
			reverseOrder: false,
			want:         0,
		},
		{
			name:         "forward a < b",
			a:            1,
			b:            2,
			reverseOrder: false,
			want:         -1,
		},
		{
			name:         "reverse a < b",
			a:            1,
			b:            2,
			reverseOrder: true,
			want:         1,
		},
		{
			name:         "forward a > b",
			a:            2,
			b:            1,
			reverseOrder: false,
			want:         1,
		},
		{
			name:         "reverse a > b",
			a:            2,
			b:            1,
			reverseOrder: true,
			want:         -1,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reversed := reversecmp.Reverser[uint16](cmp.Compare[uint16], tt.reverseOrder)
			got := reversed(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Reverser(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestReverserFloat32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		a            float32
		b            float32
		reverseOrder bool
		want         int
	}{
		{
			name:         "forward equal values",
			a:            0.0,
			b:            0.0,
			reverseOrder: false,
			want:         0,
		},
		{
			name:         "forward a < b",
			a:            1.5,
			b:            2.5,
			reverseOrder: false,
			want:         -1,
		},
		{
			name:         "reverse a < b",
			a:            1.5,
			b:            2.5,
			reverseOrder: true,
			want:         1,
		},
		{
			name:         "forward a > b",
			a:            2.5,
			b:            1.5,
			reverseOrder: false,
			want:         1,
		},
		{
			name:         "reverse a > b",
			a:            2.5,
			b:            1.5,
			reverseOrder: true,
			want:         -1,
		},
		{
			name:         "forward negative values",
			a:            -2.5,
			b:            -1.5,
			reverseOrder: false,
			want:         -1,
		},
		{
			name:         "reverse negative values",
			a:            -2.5,
			b:            -1.5,
			reverseOrder: true,
			want:         1,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			reversed := reversecmp.Reverser[float32](cmp.Compare[float32], tt.reverseOrder)
			got := reversed(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Reverser(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
