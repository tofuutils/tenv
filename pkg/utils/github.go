package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func DownloadLatestRelease(owner, repo, destFolder string) error {
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

	assets, ok := releaseInfo["assets"].([]interface{})
	if !ok || len(assets) == 0 {
		return fmt.Errorf("no assets found for the latest release")
	}

	latestAsset := assets[0].(map[string]interface{})
	latestAssetName := latestAsset["name"].(string)
	//latestAssetURL := latestAsset["browser_download_url"].(string)

	// Download the latest release zip file
	downloadPath := fmt.Sprintf("%s/%s", destFolder, latestAssetName)
	//if err := downloadFile(latestAssetURL, downloadPath); err != nil {
	//	return err
	//}

	// Unzip the downloaded file to the destination folder
	if err := unzipFile(downloadPath, destFolder); err != nil {
		return err
	}

	return nil
}

//func downloadFile(url, dest string) error {
//	resp, err := http.Get(url)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
//	}
//
//	out, err := os.Create(dest)
//	if err != nil {
//		return err
//	}
//	defer out.Close()
//
//	_, err = io.Copy(out, resp.Body)
//	return err
//}

func unzipFile(src, dest string) error {
	cmd := exec.Command("unzip", "-o", src, "-d", dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
