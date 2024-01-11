package tofu

import (
	"context"
	"fmt"
	"github.com/google/go-github/v58/github"
	"github.com/opentofuutils/tenv/pkg/utils/archive"
	"net/http"
	"os"
	"path"
	"runtime"
)

func InstallSpecificVersion(client *github.Client, owner, repo, version string) error {
	installPath := "./"
	entries, err := os.ReadDir(installPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() && version == entry.Name() {
			return nil
		}
	}

	release, _, err := client.Repositories.GetReleaseByTag(context.Background(), owner, repo, fmt.Sprintf("v%s", version))

	expectedAsset := fmt.Sprintf("tofu_%s_%s_%s.zip", version, runtime.GOOS, runtime.GOARCH)
	fmt.Println(expectedAsset)

	var assetUrl string
	for _, asset := range release.Assets {
		if *asset.Name == expectedAsset {
			assetUrl = *asset.BrowserDownloadURL
			break
		}

	}

	response, err := http.Get(assetUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	targetPath := path.Join("./", version)
	fmt.Println(targetPath)
	err = archive.ExtractZipToDir(response.Body, targetPath)
	return nil
}
