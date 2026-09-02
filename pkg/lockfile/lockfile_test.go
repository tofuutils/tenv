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

package lockfile_test

import (
	_ "embed"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tofuutils/tenv/v4/pkg/fileperm"
	"github.com/tofuutils/tenv/v4/pkg/lockfile"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

// recordingDisplayer captures Log calls so tests can assert on them.
// Safe for concurrent use since WriteWithCustomLockPath's retry loop logs from the caller goroutine only,
// but TestParallelWriteRead exercises several goroutines against a shared displayer.
type recordingDisplayer struct {
	mu   sync.Mutex
	logs []logCall
}

type logCall struct {
	Level hclog.Level
	Msg   string
	Args  []any
}

func (d *recordingDisplayer) Display(string) {}

func (d *recordingDisplayer) IsDebug() bool { return false }

func (d *recordingDisplayer) Log(level hclog.Level, msg string, args ...any) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.logs = append(d.logs, logCall{Level: level, Msg: msg, Args: args})
}

func (d *recordingDisplayer) Flush(bool) {}

func (d *recordingDisplayer) Calls() []logCall {
	d.mu.Lock()
	defer d.mu.Unlock()

	return slices.Clone(d.logs)
}

//go:embed testdata/data1.txt
var data1 []byte

//go:embed testdata/data2.txt
var data2 []byte

//go:embed testdata/data3.txt
var data3 []byte

func TestParallelWriteRead(t *testing.T) {
	t.Parallel()

	parallelDirPath := t.TempDir()
	parallelFilePath := filepath.Join(parallelDirPath, "rw_test")

	var err1, err2, err3 error
	var res1, res2, res3 []byte
	done1, done2, done3 := make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		res1, err1 = writeReadFile(parallelDirPath, parallelFilePath, data1, loghelper.InertDisplayer)
		done1 <- struct{}{}
	}()
	go func() {
		res2, err2 = writeReadFile(parallelDirPath, parallelFilePath, data2, loghelper.InertDisplayer)
		done2 <- struct{}{}
	}()
	go func() {
		res3, err3 = writeReadFile(parallelDirPath, parallelFilePath, data3, loghelper.InertDisplayer)
		done3 <- struct{}{}
	}()

	<-done1
	<-done2
	<-done3

	if err1 != nil {
		t.Error("Unexpected error with call 1 :", err1)
	}
	if err2 != nil {
		t.Error("Unexpected error with call 2 :", err2)
	}
	if err3 != nil {
		t.Error("Unexpected error with call 3 :", err1)
	}

	if !slices.Equal(res1, data1) || !slices.Equal(res2, data2) || !slices.Equal(res3, data3) {
		t.Error("Read data does not match written data")
	}
}

func writeReadFile(dirPath string, filePath string, data []byte, displayer loghelper.Displayer) ([]byte, error) {
	deleteLock := lockfile.WriteWithCustomLockPath(dirPath, ".", displayer)
	defer deleteLock()

	if err := os.WriteFile(filePath, data, fileperm.RW); err != nil {
		return nil, err
	}

	time.Sleep(100 * time.Millisecond)

	return os.ReadFile(filePath)
}

func TestCleanAndExitOnInterrupt(t *testing.T) {
	t.Parallel()

	var cleaned bool
	disableExit := lockfile.CleanAndExitOnInterrupt(func() { cleaned = true })

	// No signal was sent, so disabling must return without invoking clean and without hanging.
	disableExit()

	assert.False(t, cleaned, "clean should not run unless an interrupt signal is received")
}

func TestWriteWithCustomLockPath_CreatesAndRemovesLockFile(t *testing.T) {
	t.Parallel()

	lockDir := t.TempDir()
	displayer := &recordingDisplayer{}

	deleteLock := lockfile.WriteWithCustomLockPath(lockDir, "terraform", displayer)

	lockPath := filepath.Join(lockDir, "terraform.lock")
	_, err := os.Stat(lockPath)
	require.NoError(t, err, "lock file should exist after WriteWithCustomLockPath")

	deleteLock()

	_, err = os.Stat(lockPath)
	assert.True(t, os.IsNotExist(err), "lock file should be removed after calling the cleanup function")
}

func TestWriteWithCustomLockPath_CleanupIsIdempotent(t *testing.T) {
	t.Parallel()

	lockDir := t.TempDir()
	displayer := &recordingDisplayer{}

	deleteLock := lockfile.WriteWithCustomLockPath(lockDir, "tofu", displayer)

	assert.NotPanics(t, func() {
		deleteLock()
		deleteLock()
	}, "cleanup function must be safe to call more than once")
}

func TestWriteWithCustomLockPath_MissingLockDir(t *testing.T) {
	t.Parallel()

	missingDir := filepath.Join(t.TempDir(), "does-not-exist")
	displayer := &recordingDisplayer{}

	deleteLock := lockfile.WriteWithCustomLockPath(missingDir, "terraform", displayer)

	// No lock file must be created when the lock directory itself is missing.
	_, err := os.Stat(filepath.Join(missingDir, "terraform.lock"))
	assert.True(t, os.IsNotExist(err))

	calls := displayer.Calls()
	require.Len(t, calls, 1)
	assert.Equal(t, hclog.Error, calls[0].Level)
	assert.Equal(t, "lock directory does not exist", calls[0].Msg)
	assert.Equal(t, []any{"dir", missingDir}, calls[0].Args)

	// The returned cleanup function must still be safe to call (no-op).
	assert.NotPanics(t, deleteLock)
}

func TestWriteWithCustomLockPath_SeparateFoldersDoNotConflict(t *testing.T) {
	t.Parallel()

	lockDir := t.TempDir()
	displayer := &recordingDisplayer{}

	// Two different folder names in the same lock directory must not block each other,
	// since the lock file name is scoped by folderName.
	deleteTerraformLock := lockfile.WriteWithCustomLockPath(lockDir, "terraform", displayer)
	deleteTofuLock := lockfile.WriteWithCustomLockPath(lockDir, "tofu", displayer)

	assert.Empty(t, displayer.Calls(), "acquiring locks for distinct folders should not retry or log")

	deleteTerraformLock()
	deleteTofuLock()
}
