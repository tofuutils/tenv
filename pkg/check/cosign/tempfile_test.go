package cosigncheck

import (
	"errors"
	"os"
	"testing"
)

func TestTempFile(t *testing.T) {
	content := []byte("test content")

	fileName, cleanup, err := tempFile("test", content)
	if err != nil {
		t.Fatalf("tempFile() error = %v", err)
	}
	defer cleanup()

	got, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("failed to read temporary file: %v", err)
	}
	if string(got) != string(content) {
		t.Fatalf("file contents = %q, want %q", got, content)
	}

	cleanup()

	_, err = os.Stat(fileName)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected temporary file to be removed, stat error = %v", err)
	}
}
