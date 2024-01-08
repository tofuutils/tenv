package misc

import (
	"fmt"
	"github.com/opentofuutils/tenv/pkg/github"
	log "github.com/sirupsen/logrus"
	"os"
)

func PrepareTool(owner, repo, rootDir string) error {
	binDir := fmt.Sprintf("%s/bin", rootDir)
	miscDir := fmt.Sprintf("%s/misc", rootDir)

	// Create temporary directory where tarballs will be stored
	err := CreateFolder(miscDir)
	if err != nil {
		return err
	}
	defer func() {
		// Clean up: Remove the temporary directory when done
		err := os.RemoveAll(miscDir)
		if err != nil {
			log.Error("Error removing temporary directory:", err)
		}
	}()

	tarballPath := fmt.Sprintf("%s/%s-%s", miscDir, owner, repo)
	if err := github.DownloadLatestRelease(owner, repo, tarballPath); err != nil {
		fmt.Println("Error:", err)
		return err
	}
	log.Info(fmt.Sprintf("Latest %s release owned by %s downloaded successfully", repo, owner))

	err = UnTarGz(tarballPath, fmt.Sprintf("%s/%s", binDir, repo))
	if err != nil {
		log.Warn("Error:", err)
	} else {
		log.Info("Archive untarred successfully.")
	}

	log.Info(fmt.Sprintf("Latest %s release owned by %s prepared successfully", repo, owner))

	return nil
}
