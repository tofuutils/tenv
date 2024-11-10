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

	"github.com/tofuutils/tenv/v3/config"
	configutils "github.com/tofuutils/tenv/v3/config/utils"
)

const rwPerm = 0o600

var errDelimiter = errors.New("key and value should not contains delimiter")

// Always call os.Exit.
func Run(cmd *exec.Cmd, gha bool, getenv configutils.GetenvFunc) {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	done := initIO(cmd, &exitCode, gha, getenv)
	defer done()

	err := cmd.Start()
	if err != nil {
		exitWithErrorMsg(cmd.Path, err, &exitCode)

		return
	}

	signalChan := make(chan os.Signal, 1)
	go transmitIncreasingSignal(signalChan, cmd.Process)
	signal.Notify(signalChan, os.Interrupt)

	if err = cmd.Wait(); err != nil {
		var exitError *exec.ExitError
		if ok := errors.As(err, &exitError); ok {
			exitCode = exitError.ExitCode()

			return
		}
		exitWithErrorMsg(cmd.Path, err, &exitCode)
	}
}

func exitWithErrorMsg(execPath string, err error, pExitCode *int) {
	fmt.Println("Failure during", execPath, "call :", err) //nolint
	if *pExitCode == 0 {
		*pExitCode = 1
	}
}

func initIO(cmd *exec.Cmd, pExitCode *int, gha bool, getenv configutils.GetenvFunc) func() {
	cmd.Stdin = os.Stdin
	if !gha {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return noAction
	}

	outputPath := getenv(config.GithubOutputEnvName)
	outputFile, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, rwPerm)
	if err != nil {
		fmt.Println("Ignore GITHUB_ACTIONS, fail to open GITHUB_OUTPUT :", err) //nolint

		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return noAction
	}

	var errBuffer, outBuffer strings.Builder
	cmd.Stderr = io.MultiWriter(&errBuffer, os.Stderr)
	cmd.Stdout = io.MultiWriter(&outBuffer, os.Stdout)

	return func() {
		defer outputFile.Close()

		err = writeMultiline(outputFile, "stderr", errBuffer.String())
		if err != nil {
			exitWithErrorMsg(cmd.Path, err, pExitCode)

			return
		}

		if err = writeMultiline(outputFile, "stdout", outBuffer.String()); err != nil {
			exitWithErrorMsg(cmd.Path, err, pExitCode)

			return
		}

		exitCode := *pExitCode
		if err = writeMultiline(outputFile, "exitcode", strconv.Itoa(exitCode)); err != nil {
			exitWithErrorMsg(cmd.Path, err, pExitCode)

			return
		}

		if exitCode != 0 && exitCode != 2 {
			err = fmt.Errorf("exited with code %d", exitCode)
			exitWithErrorMsg(cmd.Path, err, pExitCode)
		}
	}
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

func noAction() {}
