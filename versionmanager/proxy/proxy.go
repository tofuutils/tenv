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

package proxy

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/versionmanager"
	terragruntparser "github.com/tofuutils/tenv/versionmanager/semantic/parser/terragrunt"
)

var errDelimiter = errors.New("key and value should not contains delimiter")

func Exec(builderFunc func(*config.Config, terragruntparser.TerragruntParser) versionmanager.VersionManager, execName string) {
	conf, err := config.InitConfigFromEnv()
	if err != nil {
		fmt.Println("Configuration error :", err) //nolint
		os.Exit(1)
	}

	conf.InitDisplayer(true)
	versionManager := builderFunc(&conf, terragruntparser.Make())
	detectedVersion, err := versionManager.Detect(true)
	if err != nil {
		fmt.Println("Failed to detect a version allowing to call", execName, ":", err) //nolint
		os.Exit(1)
	}

	RunCmd(versionManager.InstallPath(), detectedVersion, execName)
}

func RunCmd(installPath string, detectedVersion string, execName string) {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	cmdArgs := os.Args[1:]
	// proxy to selected version
	cmd := exec.Command(filepath.Join(installPath, detectedVersion, execName), cmdArgs...)
	done, err := initIO(cmd, execName, &exitCode)
	if err != nil {
		exitWithErrorMsg(execName, err, &exitCode)

		return
	}

	calledExitCode := 0
	defer func() {
		done(calledExitCode)
	}()

	if err = cmd.Start(); err != nil {
		exitWithErrorMsg(execName, err, &exitCode)

		return
	}

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	go transmitIncreasingSignal(signalChan, cmd.Process)

	if err = cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			calledExitCode = exitError.ExitCode()

			return
		}
		exitWithErrorMsg(execName, err, &exitCode)
	}
}

func exitWithErrorMsg(execName string, err error, pExitCode *int) {
	fmt.Println("Failure during", execName, "call :", err) //nolint
	*pExitCode = 1
}

func initIO(cmd *exec.Cmd, execName string, pExitCode *int) (func(int), error) {
	gha, err := config.GetenvBool(false, config.GithubActionsEnvName)
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
	outputFile, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
			process.Signal(os.Interrupt)
			first = false
		} else {
			process.Signal(os.Kill)
		}
	}
}

func writeMultiline(file *os.File, key string, value string) error {
	delimiter := "ghadelimeter_" + strconv.Itoa(rand.Int())
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
