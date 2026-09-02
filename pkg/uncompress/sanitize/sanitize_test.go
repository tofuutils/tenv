package sanitize_test

import (
	"path/filepath"
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/uncompress/sanitize"
)

func TestArchivePathClean(t *testing.T) {
	t.Parallel()

	path, err := sanitize.ArchivePath("/home/test", "index.json")
	if err != nil {
		t.Fatal("Unexpected error :", err)
	}

	expected := filepath.Join("/home/test", "index.json")
	if path != expected {
		t.Errorf("Unexpected result, get: %s, want: %s", path, expected)
	}
}

func TestArchivePathTainted(t *testing.T) {
	t.Parallel()

	if path, err := sanitize.ArchivePath("/home/test", "../index.json"); err == nil {
		t.Error("Should fail on tainted path, get :", path)
	}
}
