/*
 *
 * Copyright 2024 tofuutils authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package cosigncheck

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

const cosignExecName = "cosign"

var verified = "Verified OK"

var (
	ErrCheck        = errors.New("cosign check failed")
	ErrNotInstalled = errors.New("cosign executable not found")
)

func Check(data []byte, dataSig []byte, dataCert []byte, certIdentity string, certOidcIssuer string) error {
	_, err := exec.LookPath(cosignExecName)
	if err != nil {
		return ErrNotInstalled
	}

	dataFileName, remove, err := tempFile("data", data)
	if err != nil {
		return err
	}
	defer remove()

	dataSigFileName, remove, err := tempFile("data.sig", dataSig)
	if err != nil {
		return err
	}
	defer remove()

	dataCertFileName, remove, err := tempFile("data.cert", dataCert)
	if err != nil {
		return err
	}
	defer remove()

	cmdArgs := []string{
		"verify-blob", "--certificate-identity", certIdentity, "--signature", dataSigFileName, "--certificate", dataCertFileName,
		"--certificate-oidc-issuer", certOidcIssuer, dataFileName,
	}
	cmd := exec.Command(cosignExecName, cmdArgs...)

	if returnedData, _ := cmd.CombinedOutput(); !strings.Contains(string(returnedData), verified) {
		return ErrCheck
	}

	return nil
}

func tempFile(name string, data []byte) (string, func(), error) {
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		return "", nil, err
	}

	tmpFileName := tmpFile.Name()
	if err = os.WriteFile(tmpFileName, data, 0600); err != nil {
		return "", nil, err
	}

	return tmpFileName, func() {
		os.Remove(tmpFileName)
	}, nil
}
