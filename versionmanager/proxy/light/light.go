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

package lightproxy

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/tofuutils/tenv/v4/config/cmdconst"
	detachproxy "github.com/tofuutils/tenv/v4/versionmanager/proxy/detach"
)

func Exec(execName string) {
	cmdArgs := make([]string, len(os.Args)+1)
	cmdArgs[0], cmdArgs[1] = cmdconst.CallSubCmd, execName
	copy(cmdArgs[2:], os.Args[1:])

	// proxy to selected version
	cmd := exec.Command(cmdconst.TenvName, cmdArgs...) //nolint
	detachproxy.InitBehaviorFromEnv(cmd, os.Getenv)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		exitWithErrorMsg(execName, err)
	}

	signalChan := make(chan os.Signal, 1)
	go transmitSignal(signalChan, cmd.Process)
	signal.Notify(signalChan, os.Interrupt)

	if err = cmd.Wait(); err != nil {
		var exitError *exec.ExitError
		if ok := errors.As(err, &exitError); ok {
			os.Exit(exitError.ExitCode())
		}
		exitWithErrorMsg(execName, err)
	}
}

func exitWithErrorMsg(execName string, err error) {
	fmt.Println("Failure during", execName, "call :", err) //nolint
	os.Exit(1)
}

func transmitSignal(signalReceiver <-chan os.Signal, process *os.Process) {
	for range signalReceiver {
		_ = process.Signal(os.Interrupt)
	}
}
