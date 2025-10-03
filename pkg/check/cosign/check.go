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
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/tofuutils/tenv/v4/pkg/fileperm"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

const (
	CosignExecName = "cosign"
	Verified       = "Verified OK"
)

var (
	ErrCheck        = errors.New("cosign check failed")
	ErrNotInstalled = errors.New("cosign executable not found")
)

func Check(ctx context.Context, data []byte, dataSig []byte, dataCert []byte, certIdentity string, certOidcIssuer string, displayer loghelper.Displayer) error {
	_, err := exec.LookPath(CosignExecName)
	if err != nil {
		return ErrNotInstalled
	}

	dataFileName, remove, err := TempFile("data", data)
	if err != nil {
		return err
	}
	defer remove()

	dataSigFileName, remove, err := TempFile("data.sig", dataSig)
	if err != nil {
		return err
	}
	defer remove()

	dataCertFileName, remove, err := TempFile("data.cert", dataCert)
	if err != nil {
		return err
	}
	defer remove()

	cmdArgs := []string{
		"verify-blob", "--certificate-identity", certIdentity, "--signature", dataSigFileName, "--certificate", dataCertFileName,
		"--certificate-oidc-issuer", certOidcIssuer, dataFileName,
	}

	var outBuffer, errBuffer strings.Builder
	cmd := exec.CommandContext(ctx, CosignExecName, cmdArgs...)
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	cmd.Run() //nolint

	stdOutContent, stdErrContent := outBuffer.String(), errBuffer.String()

	displayer.Log(hclog.Debug, "cosign output", "stdOut", stdOutContent, "stdErr", stdErrContent)

	if !strings.Contains(stdErrContent, Verified) {
		return ErrCheck
	}

	return nil
}

func TempFile(name string, data []byte) (string, func(), error) {
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		return "", nil, err
	}

	tmpFileName := tmpFile.Name()
	if err = os.WriteFile(tmpFileName, data, fileperm.RW); err != nil {
		tmpFile.Close()
		os.Remove(tmpFileName)

		return "", nil, err
	}

	return tmpFileName, func() {
		tmpFile.Close()
		os.Remove(tmpFileName)
	}, nil
}
