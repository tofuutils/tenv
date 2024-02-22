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

package types

import (
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/pkg/loghelper"
)

func DisplayDetectionInfo(displayer loghelper.Displayer, version string, source string) string {
	displayer.Display(loghelper.Concat("Resolved version from ", source, " : ", version))

	return version
}

type PredicateInfo struct {
	Predicate    func(string) bool
	ReverseOrder bool
}

type PredicateReader = func(*config.Config) (func(string) bool, error)

type VersionFile struct {
	Name   string
	Parser func(string, *config.Config) (string, error)
}
