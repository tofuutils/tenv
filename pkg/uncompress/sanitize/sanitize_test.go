package sanitize_test

import (
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/uncompress/sanitize"
)

func TestArchivePathClean(t *testing.T) {
	t.Parallel()

	path, err := sanitize.ArchivePath("/home/test", "index.json")
	if err != nil {
		t.Fatal("Unexpected error :", err)
	}

	if path != "/home/test/index.json" {
		t.Error("Unexpected result, get :", path)
	}
}

func TestArchivePathTainted(t *testing.T) {
	t.Parallel()

	if path, err := sanitize.ArchivePath("/home/test", "../index.json"); err == nil {
		t.Error("Should fail on tainted path, get :", path)
	}
}

func TestArchivePathSiblingPrefixTraversal(t *testing.T) {
	t.Parallel()

	// /home/testing shares the /home/test prefix but is outside the extraction root.
	if path, err := sanitize.ArchivePath("/home/test", "../testing/index.json"); err == nil {
		t.Error("Should fail on tainted path, get :", path)
	}
}
