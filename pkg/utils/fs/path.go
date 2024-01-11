/*
 *
 * Copyright 2024 opentofuutils authors.
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
package fs

import (
	"fmt"
	"github.com/opentofuutils/tenv/pkg/utils/env"
	log "github.com/sirupsen/logrus"
)

func GetPath(name string) string {
	rootDir := env.GetEnv(env.RootEnv, "")

	switch name {
	case "root_dir":
		return rootDir
	case "bin_dir":
		return fmt.Sprintf("%s/bin", rootDir)
	case "misc_dir":
		return fmt.Sprintf("%s/misc", rootDir)
	case "tfenv_dir":
		return fmt.Sprintf("%s/bin/tfenv", rootDir)
	case "tofuenv_dir":
		return fmt.Sprintf("%s/bin/tofuenv", rootDir)
	case "tfenv_exec":
		return fmt.Sprintf("%s/bin/tfenv/bin/tfenv", rootDir)
	case "tofuenv_exec":
		return fmt.Sprintf("%s/bin/tofuenv/bin/tofuenv", rootDir)

	default:
		log.Warn("Unknown day")
		return ""
	}

}
