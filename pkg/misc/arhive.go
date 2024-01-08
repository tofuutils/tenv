package misc

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnTarGz(source, destination string) error {
	// Open the source file
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Iterate through the tar archive and extract files
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			// End of tar archive
			break
		}

		if err != nil {
			return err
		}

		// Construct the destination path for the current file without the top-level directory
		target := filepath.Join(destination, strings.TrimPrefix(header.Name, filepath.Base(source)+"/"))

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directories
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}

		case tar.TypeReg, tar.TypeRegA:
			// Create regular files
			file, err := os.Create(target)
			if err != nil {
				return err
			}
			defer file.Close()

			// Copy file contents
			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}

			// Set file permissions
			if err := os.Chmod(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		}
	}

	return nil
}
