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

func Check(data []byte, dataSums []byte, fileName string) error {
	dataSumsStr := string(dataSums)
	for _, dataSumStr := range strings.Split(dataSumsStr, "\n") {
		dataSumStr, ok := strings.CutSuffix(dataSumStr, fileName)
		if ok {
			dataSumStr = strings.TrimSpace(dataSumStr)
			dataSum, err := hex.DecodeString(dataSumStr)
			if err != nil {
				return err
			}

			hashed := sha256.Sum256(data)
			if !bytes.Equal(dataSum, hashed[:]) {
				return errCheck
			}
			return nil
		}
	}
	return errNoSum
}
