package sha256check

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
)

var (
	errCheck = errors.New("invalid sha256 checksum")
	errNoSig = errors.New("file sha256 not found for current platform")
)

func Check(data []byte, dataSig []byte) error {
	hasher := sha256.New()
	if _, err := hasher.Write(data); err != nil {
		return err
	}

	if !bytes.Equal(dataSig, hasher.Sum(nil)) {
		return errCheck
	}
	return nil
}

func Extract(dataSigs []byte, fileName string) ([]byte, error) {
	dataSigsStr := string(dataSigs)
	for _, dataSigStr := range strings.Split(dataSigsStr, "\n") {
		if dataSigStr, ok := strings.CutSuffix(dataSigStr, fileName); ok {
			dataSigStr = strings.TrimSpace(dataSigStr)
			return base64.StdEncoding.DecodeString(dataSigStr)
		}
	}
	return nil, errNoSig
}
