/*
 *
 * Copyright 2024 gotofuenv authors.
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

package tofuretriever

import (
	"runtime"
	"strings"

	"github.com/dvaumoron/gotofuenv/pkg/github"
)

type TofuRetriever struct {
	authorizationHeader string
	githubReleaseUrl    string
}

func MakeTofuRetriever(githubReleaseUrl string, githubToken string) TofuRetriever {
	return TofuRetriever{
		authorizationHeader: github.BuildAuthorizationHeader(githubToken),
		githubReleaseUrl:    githubReleaseUrl,
	}
}

func (v TofuRetriever) DownloadAssetUrl(version string) (string, error) {
	tag := version
	// assume that opentofu tags start with a 'v'
	// and version in asset name does not
	if tag[0] == 'v' {
		version = version[1:]
	} else {
		tag = "v" + version
	}
	return github.DownloadAssetUrl(tag, buildAssetName(version), v.githubReleaseUrl, v.authorizationHeader)
}

func (v TofuRetriever) LatestRelease() (string, error) {
	return github.LatestRelease(v.githubReleaseUrl, v.authorizationHeader)
}

func (v TofuRetriever) ListReleases() ([]string, error) {
	return github.ListReleases(v.githubReleaseUrl, v.authorizationHeader)
}

func buildAssetName(version string) string {
	var nameBuilder strings.Builder
	nameBuilder.WriteString("tofu_")
	nameBuilder.WriteString(version)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(runtime.GOOS)
	nameBuilder.WriteByte('_')
	nameBuilder.WriteString(runtime.GOARCH)
	nameBuilder.WriteString(".zip")
	return nameBuilder.String()
}
