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

package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tofuutils/tenv/v4/config"
	"github.com/tofuutils/tenv/v4/pkg/loghelper"
)

func TestInitRemoteConf(t *testing.T) {
	t.Parallel()

	confPath := filepath.Join(t.TempDir(), "remote.yaml")
	data := "terragrunt:\n  url: https://mirror.example.com/tg\nterramate:\n  url: https://mirror.example.com/tm\n"
	if err := os.WriteFile(confPath, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}

	conf := config.Config{Displayer: loghelper.InertDisplayer, RemoteConfPath: confPath}
	if err := conf.InitRemoteConf(); err != nil {
		t.Fatal("Unexpected error : ", err)
	}

	if url := conf.Tg.GetRemoteURL(); url != "https://mirror.example.com/tg" {
		t.Error("Unexpected terragrunt remote url : ", url)
	}
	if url := conf.Tm.GetRemoteURL(); url != "https://mirror.example.com/tm" {
		t.Error("Unexpected terramate remote url : ", url)
	}
}
