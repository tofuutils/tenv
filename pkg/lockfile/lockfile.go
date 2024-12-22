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
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

const (
	msgWrite  = "can not write .lock file, will retry"
	msgDelete = "can not remove .lock file"

	rwPerm = 0o600
)

// ! dirPath must already exist (no mkdir here).
// the returned function must be used to delete the lock.
func Write(dirPath string, displayer loghelper.Displayer) func() {
	lockPath := filepath.Join(dirPath, ".lock")
	for logLevel := hclog.Warn; true; logLevel = hclog.Info {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL, rwPerm)
		if err == nil {
			f.Close()

			break
		}

		displayer.Log(logLevel, msgWrite, loghelper.Error, err)
		time.Sleep(time.Second)
	}

	return sync.OnceFunc(func() {
		if err := os.RemoveAll(lockPath); err != nil {
			displayer.Log(hclog.Warn, msgDelete, loghelper.Error, err)
		}
	})
}

// the returned function may be used to avoid goroutine leak
// (also avoid conflicting behavior with versionmanager/proxy.transmitIncreasingSignal).
func CleanAndExitOnInterrupt(clean func()) func() {
	signalChan := make(chan os.Signal, 1)
	endChan := make(chan struct{})
	go listenToClean(signalChan, endChan, clean)
	signal.Notify(signalChan, os.Interrupt)

	return sync.OnceFunc(func() {
		endChan <- struct{}{}
	})
}

func listenToClean(signalChan chan os.Signal, endChan <-chan struct{}, clean func()) {
	for {
		select {
		case <-signalChan:
			clean()
			os.Exit(1)
		case <-endChan:
			signal.Stop(signalChan)

			break
		}
	}
}
