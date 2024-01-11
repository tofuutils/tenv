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
package tool

import (
	"fmt"
	"github.com/opentofuutils/tenv/pkg/utils/archive"
	"github.com/opentofuutils/tenv/pkg/utils/fs"
	"github.com/opentofuutils/tenv/pkg/utils/github"
	log "github.com/sirupsen/logrus"

	"os"
)

func CheckToolInstalled(name string) bool {

	path := fs.GetPath("tofuenv_exec")
	_, err := os.Stat(path)

	return !os.IsExist(err)
}

func PrepareTool(owner, repo, rootDir string) error {
	binDir := fs.GetPath("bin_dir")
	miscDir := fs.GetPath("misc_dir")

	// Create temporary directory where tarballs will be stored
	err := fs.CreateFolder(miscDir)
	if err != nil {
		return err
	}

	defer func() {
		err = fs.DeleteFolder(miscDir)
		if err != nil {
			log.Error("Error removing temporary directory:", err)
		}
	}()

	tarballPath := fmt.Sprintf("%s/%s-%s", miscDir, owner, repo)
	if err := github.DownloadLatestRelease(owner, repo, tarballPath); err != nil {
		log.Error("Error:", err)
		return err
	}
	log.Info(fmt.Sprintf("Latest %s release owned by %s downloaded successfully", repo, owner))

	err = archive.ExtractTarGz(tarballPath, fmt.Sprintf("%s/%s", binDir, repo))
	if err != nil {
		log.Warn("Error:", err)
	} else {
		log.Info("Archive untarred successfully.")
	}

	log.Info(fmt.Sprintf("Latest %s release owned by %s prepared successfully", repo, owner))

	return nil
}
