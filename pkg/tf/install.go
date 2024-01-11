package tf

import (
	"fmt"
	"github.com/opentofuutils/tenv/pkg/utils/archive"
	"net/http"
	"os"
	"path"
	"runtime"
)

func InstallSpecificVersion(version string) error {
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

	url := fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip", version, version, runtime.GOOS, runtime.GOARCH)
	fmt.Println(url)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	targetPath := path.Join("./", version)
	fmt.Println(targetPath)
	err = archive.ExtractZipToDir(response.Body, targetPath)
	return nil
}
