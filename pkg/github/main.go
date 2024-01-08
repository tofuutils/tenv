package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadLatestRelease(owner, repo, destPath string) error {
	// Use GitHub API to get the latest release information
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch release information, status code: %d", resp.StatusCode)
	}

	// Parse the response JSON to get the download URL for the latest release
	var releaseInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&releaseInfo); err != nil {
		return err
	}

	zipball_url := releaseInfo["tarball_url"].(string)

	// Download the latest release zip file
	if err := DownloadFile(zipball_url, destPath); err != nil {
		return err
	}

	return nil
}

func DownloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
