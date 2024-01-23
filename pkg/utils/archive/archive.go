package archive

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path"
)

func ExtractZipToDir(zipBodyReader io.Reader, dirPath string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}

	zipBody, err := io.ReadAll(zipBodyReader)
	if err != nil {
		return err
	}

	byteReader := bytes.NewReader(zipBody)
	zipReader, err := zip.NewReader(byteReader, int64(len(zipBody)))
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		if err = copyZipFileToDir(file, dirPath); err != nil {
			return err
		}
	}
	return nil
}

func copyZipFileToDir(zipFile *zip.File, dirPath string) error {
	destPath := path.Join(dirPath, zipFile.Name)
	if destPath[len(destPath)-1] == '/' {
		// trailing slash indicates a directory
		return os.MkdirAll(destPath, 0755)
	}

	reader, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return os.WriteFile(destPath, data, zipFile.Mode())
}
