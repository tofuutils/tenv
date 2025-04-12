package sanitize

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Sanitize archive file pathing from "G305" (file traversal).
func ArchivePath(dirPath string, fileName string) (string, error) {
	destPath := filepath.Join(dirPath, fileName)
	if strings.HasPrefix(destPath, filepath.Clean(dirPath)) {
		return destPath, nil
	}

	return "", fmt.Errorf("content filepath is tainted: %s", fileName)
}
