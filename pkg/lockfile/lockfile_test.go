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
	"testing"
	"time"

	"github.com/tofuutils/tenv/v4/pkg/fileperm"
	"github.com/tofuutils/tenv/v4/pkg/lockfile"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

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
	deleteLock := lockfile.Write(dirPath, displayer)
	defer deleteLock()

	if err := os.WriteFile(filePath, data, fileperm.RW); err != nil {
		return nil, err
	}

	time.Sleep(100 * time.Millisecond)

	return os.ReadFile(filePath)
}

func TestCleanAndExitOnInterrupt(t *testing.T) {
	t.Parallel()

	cleanCalled := false
	clean := func() {
		cleanCalled = true
	}

	// Start the interrupt handler
	stop := lockfile.CleanAndExitOnInterrupt(clean)

	// Stop the handler without sending interrupt
	stop()

	if cleanCalled {
		t.Error("Clean function should not have been called")
	}
}

func TestWriteLockFileExists(t *testing.T) {
	t.Parallel()

	testDirPath := filepath.Join(os.TempDir(), "locktest")
	err := os.MkdirAll(testDirPath, 0o755)
	if err != nil {
		t.Fatal("Failed to create test directory:", err)
	}
	defer os.RemoveAll(testDirPath)

	// Create a lock file manually
	lockPath := filepath.Join(testDirPath, ".lock")
	if err := os.WriteFile(lockPath, []byte("test"), fileperm.RW); err != nil {
		t.Fatal("Failed to create test lock file:", err)
	}

	// Try to acquire lock with a timeout to avoid infinite wait
	done := make(chan bool)
	go func() {
		deleteLock := lockfile.Write(testDirPath, loghelper.InertDisplayer)
		defer deleteLock()
		done <- true
	}()

	// Wait for a short time to ensure the lock is attempted
	select {
	case <-done:
		// Lock was acquired (after the existing one was considered stale)
	case <-time.After(3 * time.Second):
		t.Error("Lock acquisition took too long")
	}
}

func TestWriteLockFilePermissions(t *testing.T) {
	t.Parallel()

	testDirPath := filepath.Join(os.TempDir(), "locktest-perm")
	err := os.MkdirAll(testDirPath, 0o755)
	if err != nil {
		t.Fatal("Failed to create test directory:", err)
	}
	defer os.RemoveAll(testDirPath)

	deleteLock := lockfile.Write(testDirPath, loghelper.InertDisplayer)
	defer deleteLock()

	// Check if lock file exists with correct permissions
	lockPath := filepath.Join(testDirPath, ".lock")
	info, err := os.Stat(lockPath)
	if err != nil {
		t.Fatal("Lock file was not created:", err)
	}

	// Check permissions (considering umask)
	if info.Mode().Perm()&0o600 != 0o600 {
		t.Errorf("Lock file has wrong permissions. Got: %v, want at least: %v", info.Mode().Perm(), 0o600)
	}
}
