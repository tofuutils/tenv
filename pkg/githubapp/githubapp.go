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

package githubapp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation/v2"

	"github.com/tofuutils/tenv/v4/config/envname"
)

const githubAPIURL = "https://api.github.com/"

// InstallationToken generates a GitHub App installation access token.
// installationIDStr is optional; when empty the first installation returned
// by the API is used.
func InstallationToken(ctx context.Context, appIDStr, installationIDStr, pem, pemFile string) (string, error) {
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid %s: %w", envname.TenvGithubAppID, err)
	}

	var keyData []byte
	switch {
	case pem != "":
		keyData = []byte(pem)
	case pemFile != "":
		keyData, err = os.ReadFile(pemFile)
		if err != nil {
			return "", fmt.Errorf("reading %s: %w", envname.TenvGithubAppPEMFile, err)
		}
	default:
		return "", fmt.Errorf("one of %s or %s must be set", envname.TenvGithubAppPEM, envname.TenvGithubAppPEMFile)
	}

	installationID, err := resolveInstallationID(ctx, appID, installationIDStr, keyData, githubAPIURL)
	if err != nil {
		return "", err
	}

	itr, err := ghinstallation.New(http.DefaultTransport, appID, installationID, keyData)
	if err != nil {
		return "", fmt.Errorf("creating GitHub App transport: %w", err)
	}

	token, err := itr.Token(ctx)
	if err != nil {
		return "", fmt.Errorf("getting GitHub App installation token: %w", err)
	}

	return token, nil
}

// apiURL is parameterized so tests can point it at an httptest server.
func resolveInstallationID(ctx context.Context, appID int64, installationIDStr string, keyData []byte, apiURL string) (int64, error) { //nolint:unparam
	if installationIDStr != "" {
		id, err := strconv.ParseInt(installationIDStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid %s: %w", envname.TenvGithubAppInstallationID, err)
		}

		return id, nil
	}

	atr, err := ghinstallation.NewAppsTransport(http.DefaultTransport, appID, keyData)
	if err != nil {
		return 0, fmt.Errorf("creating GitHub App JWT transport: %w", err)
	}

	installationsURL, err := url.JoinPath(apiURL, "app/installations")
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, installationsURL, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28") //nolint

	resp, err := (&http.Client{Transport: atr}).Do(req)
	if err != nil {
		return 0, fmt.Errorf("listing GitHub App installations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GitHub API returned %d when listing installations", resp.StatusCode)
	}

	var installations []struct {
		ID int64 `json:"id"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&installations); err != nil {
		return 0, fmt.Errorf("decoding GitHub App installations: %w", err)
	}

	if len(installations) == 0 {
		return 0, errors.New("GitHub App has no installations")
	}

	return installations[0].ID, nil
}
