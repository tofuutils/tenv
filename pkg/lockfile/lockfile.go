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

package lockfile

import (
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/tofuutils/tenv/pkg/loghelper"
)

const (
	msgWrite  = "can not write .lock file, will retry"
	msgDelete = "can not remove .lock file"
)

func Write(dirPath string, displayer loghelper.Displayer) func() {
	lockPath := filepath.Join(dirPath, ".lock")

	continueB := true
	for continueB {
		if _, err := os.OpenFile(lockPath, os.O_WRONLY|os.O_CREATE, 0644); err == nil { //nolint
			continueB = false
		} else {
			displayer.Log(hclog.Debug, msgWrite, loghelper.Error, err)
			time.Sleep(time.Second)
		}
	}

	return func() {
		if err := os.RemoveAll(lockPath); err != nil {
			displayer.Log(hclog.Warn, msgDelete, loghelper.Error, err)
		}
	}
}
