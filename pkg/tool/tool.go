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
