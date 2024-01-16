package sha256check

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

var (
	errCheck = errors.New("invalid sha256 checksum")
	errNoSum = errors.New("file sha256 checksum not found for current platform")
)

func Check(data []byte, dataSum []byte) error {
	hashed := sha256.Sum256(data)
	if !bytes.Equal(dataSum, hashed[:]) {
		return errCheck
	}
	return nil
}

func Extract(dataSums []byte, fileName string) ([]byte, error) {
	dataSumsStr := string(dataSums)
	for _, dataSumStr := range strings.Split(dataSumsStr, "\n") {
		dataSumStr, ok := strings.CutSuffix(dataSumStr, fileName)
		if ok {
			dataSumStr = strings.TrimSpace(dataSumStr)
			return hex.DecodeString(dataSumStr)
		}
	}
	return nil, errNoSum
}
