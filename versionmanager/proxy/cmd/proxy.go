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

package cmdproxy

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"

	"github.com/tofuutils/tenv/v2/config/cmdconst"
	configutils "github.com/tofuutils/tenv/v2/config/utils"
)

var errDelimiter = errors.New("key and value should not contains delimiter")

func Run(execPath string, cmdArgs []string) {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	// proxy to selected version
	cmd := exec.Command(execPath, cmdArgs...)
	done, err := initIO(cmd, execPath, &exitCode)
	if err != nil {
		exitWithErrorMsg(execPath, err, &exitCode)

		return
	}

	calledExitCode := 0
	defer func() {
		done(calledExitCode)
	}()

	if err = cmd.Start(); err != nil {
		exitWithErrorMsg(execPath, err, &exitCode)

		return
	}

	signalChan := make(chan os.Signal)
	go transmitIncreasingSignal(signalChan, cmd.Process)
	signal.Notify(signalChan, os.Interrupt) //nolint

	if err = cmd.Wait(); err != nil {
		var exitError *exec.ExitError
		if ok := errors.As(err, &exitError); ok {
			calledExitCode = exitError.ExitCode()

			return
		}
		exitWithErrorMsg(execPath, err, &exitCode)
	}
}

func exitWithErrorMsg(execName string, err error, pExitCode *int) {
	fmt.Println("Failure during", execName, "call :", err) //nolint
	*pExitCode = 1
}

func initIO(cmd *exec.Cmd, execName string, pExitCode *int) (func(int), error) {
	gha, err := configutils.GetenvBool(false, cmdconst.GithubActionsEnvName)
	if err != nil {
		return nil, err
	}

	cmd.Stdin = os.Stdin
	if !gha {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return func(calledExitCode int) {
			if calledExitCode != 0 {
				*pExitCode = calledExitCode
			}
		}, nil
	}

	outputPath := os.Getenv("GITHUB_OUTPUT")
	outputFile, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644) //nolint
	if err != nil {
		return nil, err
	}

	var errBuffer, outBuffer strings.Builder
	cmd.Stderr = io.MultiWriter(&errBuffer, os.Stderr)
	cmd.Stdout = io.MultiWriter(&outBuffer, os.Stdout)

	return func(calledExitCode int) {
		var err error
		defer func() {
			if err != nil {
				exitWithErrorMsg(execName, err, pExitCode)
			}
		}()
		defer outputFile.Close()

		if err = writeMultiline(outputFile, "stderr", errBuffer.String()); err != nil {
			return
		}

		if err = writeMultiline(outputFile, "stdout", outBuffer.String()); err != nil {
			return
		}

		if err = writeMultiline(outputFile, "exitcode", strconv.Itoa(calledExitCode)); err != nil {
			return
		}

		if calledExitCode != 0 && calledExitCode != 2 {
			err = fmt.Errorf("exited with code %d", calledExitCode)
		}
	}, nil
}

func transmitIncreasingSignal(signalReceiver <-chan os.Signal, process *os.Process) {
	first := true
	for range signalReceiver {
		if first {
			_ = process.Signal(os.Interrupt)
			first = false
		} else {
			_ = process.Signal(os.Kill)
		}
	}
}

func writeMultiline(file *os.File, key string, value string) error {
	delimiter := "ghadelimeter_" + strconv.Itoa(rand.Int()) //nolint
	if strings.Contains(key, delimiter) || strings.Contains(value, delimiter) {
		return errDelimiter
	}

	var builder strings.Builder
	builder.WriteString(key)
	builder.WriteString("<<")
	builder.WriteString(delimiter)
	builder.WriteRune('\n')
	builder.WriteString(value)
	builder.WriteRune('\n')
	builder.WriteString(delimiter)
	builder.WriteRune('\n')
	_, err := file.WriteString(builder.String())

	return err
}
