//go:build !windows

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

package detachproxy

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/tofuutils/tenv/v4/config"
	configutils "github.com/tofuutils/tenv/v4/config/utils"
)

const (
	msgDisableErr = msgStart + config.TenvDetachedProxyEnvName + " environment variable, disable behavior :"
	msgKeepErr    = msgStart + config.CiEnvName + " or " + config.PipelineWsEnvName + " environment variable, keep detached default behavior to true :"
	msgStart      = "Failed to read "
)

func InitBehaviorFromEnv(cmd *exec.Cmd, getenv configutils.GetenvFunc) {
	ciEnv, ciErr := getenv.BoolFallback(false, config.CiEnvName, config.PipelineWsEnvName)
	if ciErr != nil {
		fmt.Println(msgKeepErr, ciErr) //nolint
	}

	defaultDetached := !ciEnv
	detached, err := getenv.Bool(defaultDetached, config.TenvDetachedProxyEnvName)
	if err != nil {
		fmt.Println(msgDisableErr, err) //nolint
	}

	if !detached {
		return
	}

	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	cmd.SysProcAttr.Setpgid = true
	cmd.SysProcAttr.Foreground = true
}
