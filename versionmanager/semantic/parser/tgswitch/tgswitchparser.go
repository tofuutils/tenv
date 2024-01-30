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

package tgswitchparser

import (
	"fmt"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/tofuutils/tenv/config"
)

const (
	tomlName = ".tgswitch.toml"

	versionName = "version"

	msgTgSwitchErr = "Failed to read tgswitch file :"
)

func RetrieveTerraguntVersion(conf *config.Config) (string, error) {
	var parsed map[string]string
	data, err := os.ReadFile(tomlName)
	if err == nil {
		if _, err = toml.Decode(string(data), &parsed); err != nil {
			return "", err
		}

		return parsed[versionName], nil
	}
	if conf.Verbose {
		fmt.Println(msgTgSwitchErr, err) //nolint
	}

	data, err = os.ReadFile(path.Join(conf.UserPath, tomlName))
	if err == nil {
		if _, err = toml.Decode(string(data), &parsed); err != nil {
			return "", err
		}

		return parsed[versionName], nil
	}
	if conf.Verbose {
		fmt.Println(msgTgSwitchErr, err) //nolint
	}

	data, err = os.ReadFile(path.Join(conf.RootPath, tomlName))
	if err != nil {
		if conf.Verbose {
			fmt.Println(msgTgSwitchErr, err) //nolint
		}

		return "", nil
	}

	if _, err = toml.Decode(string(data), &parsed); err != nil {
		return "", err
	}

	return parsed[versionName], nil
}
