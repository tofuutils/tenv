/*
 *
 * Copyright 2024 gotofuenv authors.
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

package tofuversion

import "github.com/dvaumoron/gotofuenv/config"

type Version struct {
	Name string
	Used bool
}

func Install(requestedVersion string, conf *config.Config) error {
	// TODO
	return nil
}

func List(conf *config.Config) ([]Version, error) {
	// TODO
	return nil, nil
}

func ListRemote(conf *config.Config) ([]string, error) {
	// TODO
	return nil, nil
}

func Uninstall(requestedVersion string, conf *config.Config) error {
	// TODO
	return nil
}

func Use(requestedVersion string, conf *config.Config) error {
	// TODO
	return nil
}
