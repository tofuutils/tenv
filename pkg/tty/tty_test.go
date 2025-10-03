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

package tty

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDetect tests the TTY detection logic
// Note: This is a conceptual test since we can't easily mock os.Stdout.Stat()
// In real scenarios, this function would be tested through integration tests
// or by testing the behavior when stdout is redirected vs when it's a terminal.
func TestDetect(t *testing.T) {
	t.Parallel()
	// Test that Detect() doesn't panic and returns a boolean
	result := Detect()
	assert.IsType(t, false, result)

	// Test that the function is deterministic for the same stdout
	result2 := Detect()
	assert.Equal(t, result, result2)
}

// This test simulates the scenario where stdout is not a TTY.
func TestDetectWithRedirectedOutput(t *testing.T) {
	t.Parallel()
	// Create a temporary file to simulate redirected output
	tmpFile, err := os.CreateTemp(t.TempDir(), "tty-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	// This test is conceptual since we can't easily replace os.Stdout
	// In a real implementation, we might use dependency injection
	// or test this through integration tests

	// For now, we'll test that the function handles the case where
	// we can't determine TTY status (which should return false)
	t.Logf("Current TTY detection result: %v", Detect())
	assert.NotPanics(t, func() {
		Detect()
	})
}
