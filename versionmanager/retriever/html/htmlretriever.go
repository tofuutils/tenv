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
	"context"
	"net/url"

	"github.com/PuerkitoBio/goquery"

	"github.com/tofuutils/tenv/v3/config"
	"github.com/tofuutils/tenv/v3/pkg/download"
	"github.com/tofuutils/tenv/v3/pkg/htmlquery"
	versionfinder "github.com/tofuutils/tenv/v3/versionmanager/semantic/finder"
)

func BuildAssetURLs(baseAssetURL string, assetNames ...string) ([]string, error) {
	joinTransformer := func(assetName string) (string, error) {
		return url.JoinPath(baseAssetURL, assetName) //nolint
	}

	return download.ApplyUrlTranformer(joinTransformer, assetNames...)
}

func ListReleases(ctx context.Context, baseURL string, remoteConf map[string]string, ro []download.RequestOption) ([]string, error) {
	selector := config.MapGetDefault(remoteConf, "selector", "a")
	extractor := htmlquery.SelectionExtractor(config.MapGetDefault(remoteConf, "part", "href"))
	versionExtractor := func(s *goquery.Selection) string {
		return versionfinder.Find(extractor(s))
	}

	return htmlquery.Request(ctx, baseURL, selector, versionExtractor, ro...)
}
