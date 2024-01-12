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

package terraformretriever

import "github.com/dvaumoron/gotofuenv/config"

type TerraformRetriever struct {
	conf *config.Config
}

func MakeTerraformRetriever(conf *config.Config) TerraformRetriever {
	return TerraformRetriever{conf: conf}
}

func (v TerraformRetriever) DownloadAssetUrl(version string) (string, error) {
	// TODO call hashicorp release api
	return "", nil
}

func (v TerraformRetriever) LatestRelease() (string, error) {
	// TODO call hashicorp release api
	return "", nil
}

func (v TerraformRetriever) ListReleases() ([]string, error) {
	// TODO call hashicorp release api
	return nil, nil
}
