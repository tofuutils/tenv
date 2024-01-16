package sha256check

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/dvaumoron/gotofuenv/pkg/apierrors"
)

func Check(data []byte, dataSum []byte) error {
	hashed := sha256.Sum256(data)
	if !bytes.Equal(dataSum, hashed[:]) {
		return apierrors.ErrCheck
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
	return nil, apierrors.ErrNoSum
}
