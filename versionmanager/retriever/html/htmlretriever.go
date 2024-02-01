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

package htmlretriever

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/tofuutils/tenv/config"
	"github.com/tofuutils/tenv/pkg/apimsg"
	"github.com/tofuutils/tenv/pkg/htmlquery"
	versionfinder "github.com/tofuutils/tenv/versionmanager/semantic/finder"
)

const (
	InstallMode       = "install_mode"
	InstallModeDirect = "direct"
	ListMode          = "list_mode"
	ListModeHTML      = "html"
	ListURL           = "list_url"
	URL               = "url"
)

func BuildAssetURLs(baseAssetURL string, assetNames []string) ([]string, error) {
	assetURLs := make([]string, 0, len(assetNames))
	for _, assetName := range assetNames {
		assetURL, err := url.JoinPath(baseAssetURL, assetName)
		if err != nil {
			return nil, err
		}

		assetURLs = append(assetURLs, assetURL)
	}

	return assetURLs, nil
}

func ListReleases(baseURL string, remoteConf map[string]string, verbose bool) ([]string, error) {
	if verbose {
		fmt.Println(apimsg.MsgFetchAllReleases, baseURL) //nolint
	}

	selector := config.MapGetDefault(remoteConf, "selector", "a")
	extractor := htmlquery.SelectionExtractor(config.MapGetDefault(remoteConf, "part", "href"))
	versionExtractor := func(s *goquery.Selection) string {
		return versionfinder.Find(extractor(s))
	}

	return htmlquery.Request(baseURL, selector, versionExtractor)
}
