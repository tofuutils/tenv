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

	"github.com/tofuutils/tenv/v4/config"
	configutils "github.com/tofuutils/tenv/v4/config/utils"
)

const (
	msgErr        = msgStart + "failed to read environment variable :"
	msgNotApplied = msgStart + "can not apply environment variable"
	msgStart      = config.TenvDetachedProxyEnvName + " behavior is always disabled on Windows OS, "
)

func InitBehaviorFromEnv(_ *exec.Cmd, getenv configutils.GetenvFunc) {
	switch detached, err := getenv.Bool(false, config.TenvDetachedProxyEnvName); {
	case err != nil:
		fmt.Println(msgErr, err) //nolint
	case detached:
		fmt.Println(msgNotApplied) //nolint
	}
}
