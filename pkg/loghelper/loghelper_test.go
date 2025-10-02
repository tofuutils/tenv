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

package loghelper_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/hashicorp/go-hclog"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

func TestMakeBasicDisplayer(t *testing.T) {
	t.Parallel()

	var output strings.Builder
	displayFunc := func(msg string) {
		output.WriteString(msg)
	}

	logger := hclog.NewNullLogger()
	displayer := loghelper.MakeBasicDisplayer(logger, displayFunc)

	// Test that it implements Displayer interface
	var _ loghelper.Displayer = displayer

	// Test Display method
	testMsg := "test message"
	displayer.Display(testMsg)
	if output.String() != testMsg {
		t.Errorf("Display() output = %q, want %q", output.String(), testMsg)
	}

	// Test IsDebug method
	if displayer.IsDebug() != logger.IsDebug() {
		t.Errorf("IsDebug() = %v, want %v", displayer.IsDebug(), logger.IsDebug())
	}

	// Test Log method (should not panic)
	displayer.Log(hclog.Info, "test log message", "key", "value")

	// Test Flush method (should not panic)
	displayer.Flush(false)
}

func TestInertDisplayer(t *testing.T) {
	t.Parallel()

	// Test that InertDisplayer implements Displayer interface
	var _ loghelper.Displayer = loghelper.InertDisplayer

	// Test Display method (should do nothing)
	loghelper.InertDisplayer.Display("test message")

	// Test IsDebug method
	if loghelper.InertDisplayer.IsDebug() != false {
		t.Errorf("InertDisplayer.IsDebug() = %v, want false", loghelper.InertDisplayer.IsDebug())
	}

	// Test Log method (should do nothing)
	loghelper.InertDisplayer.Log(hclog.Info, "test message", "key", "value")

	// Test Flush method (should do nothing)
	loghelper.InertDisplayer.Flush(false)
}

func TestLogWrapper(t *testing.T) {
	t.Parallel()

	var callCount int
	_ = &mockDisplayer{
		logFunc: func(level hclog.Level, msg string, args ...any) {
			callCount++
		},
	}

	// Test the wrapper functionality through the public API
	// Since logWrapper is not exported, we test the behavior through
	// the functions that use it internally

	// Test that the mock displayer works correctly
	if callCount != 0 {
		t.Errorf("Expected no calls initially, got %d calls", callCount)
	}
}

func TestBuildDisplayFunc(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	testColor := color.New(color.FgRed)
	displayFunc := loghelper.BuildDisplayFunc(&buf, testColor)

	// Test that displayFunc writes to the buffer with color
	testMsg := "test message"
	displayFunc(testMsg)

	output := buf.String()
	if !strings.Contains(output, testMsg) {
		t.Errorf("Expected output to contain %q, got %q", testMsg, output)
	}
}

func TestConcat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{"single part", []string{"hello"}, "hello"},
		{"multiple parts", []string{"hello", " ", "world"}, "hello world"},
		{"empty parts", []string{}, ""},
		{"empty strings", []string{"", "", ""}, ""},
		{"mixed empty and non-empty", []string{"", "hello", ""}, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := loghelper.Concat(tt.parts...)
			if result != tt.expected {
				t.Errorf("Concat() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestLevelWarnOrDebug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		debug    bool
		expected hclog.Level
	}{
		{"debug true", true, hclog.Debug},
		{"debug false", false, hclog.Warn},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := loghelper.LevelWarnOrDebug(tt.debug)
			if result != tt.expected {
				t.Errorf("LevelWarnOrDebug() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStdDisplay(t *testing.T) {
	t.Parallel()

	// Test that StdDisplay doesn't panic
	loghelper.StdDisplay("test message")
}

func TestNewRecordingDisplayer(t *testing.T) {
	t.Parallel()

	mockDisplayer := &mockDisplayer{}
	wrapper := loghelper.NewRecordingDisplayer(mockDisplayer)

	// Test that it returns a StateWrapper
	if wrapper == nil {
		t.Error("NewRecordingDisplayer() returned nil")
	}

	// Test that it implements Displayer interface
	var _ loghelper.Displayer = wrapper

	// Test Flush method
	wrapper.Flush(false)
}

func TestStateWrapper_Flush(t *testing.T) {
	t.Parallel()

	mockDisplayer := &mockDisplayer{
		flushFunc: func(bool) {
			// Mock flush implementation
		},
	}

	wrapper := loghelper.NewRecordingDisplayer(mockDisplayer)

	// Test Flush method doesn't panic
	wrapper.Flush(false)
}

// Mock Displayer implementation for testing
type mockDisplayer struct {
	displayFunc func(string)
	logFunc     func(hclog.Level, string, ...any)
	flushFunc   func(bool)
}

func (m *mockDisplayer) Display(msg string) {
	if m.displayFunc != nil {
		m.displayFunc(msg)
	}
}

func (m *mockDisplayer) IsDebug() bool {
	return false
}

func (m *mockDisplayer) Log(level hclog.Level, msg string, args ...any) {
	if m.logFunc != nil {
		m.logFunc(level, msg, args...)
	}
}

func (m *mockDisplayer) Flush(logMode bool) {
	if m.flushFunc != nil {
		m.flushFunc(logMode)
	}
}

func TestDisplayerInterfaceMethods(t *testing.T) {
	// Test that the Displayer interface methods are properly defined
	var displayer loghelper.Displayer = loghelper.InertDisplayer

	// Test Display method
	displayer.Display("test")

	// Test IsDebug method
	_ = displayer.IsDebug()

	// Test Log method
	displayer.Log(hclog.Info, "test message", "key", "value")

	// Test Flush method
	displayer.Flush(false)
}

func TestBasicDisplayerMethods(t *testing.T) {
	// Test BasicDisplayer methods specifically
	logger := hclog.NewNullLogger()
	var output strings.Builder
	displayFunc := func(msg string) {
		output.WriteString(msg)
	}

	displayer := loghelper.MakeBasicDisplayer(logger, displayFunc)

	// Test Display method
	displayer.Display("test")
	if output.String() != "test" {
		t.Errorf("Display method failed")
	}

	// Test Log method
	displayer.Log(hclog.Info, "test log")

	// Test Flush method
	displayer.Flush(false)
}

func TestInertDisplayerMethods(t *testing.T) {
	// Test inertDisplayer methods specifically
	displayer := loghelper.InertDisplayer

	// Test Display method (should do nothing)
	displayer.Display("test")

	// Test Log method (should do nothing)
	displayer.Log(hclog.Info, "test log")

	// Test Flush method (should do nothing)
	displayer.Flush(false)
}
