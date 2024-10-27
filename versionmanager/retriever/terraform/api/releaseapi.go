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

package releaseapi

import "github.com/tofuutils/tenv/v3/pkg/apimsg"

func ExtractAssetURLs(searchedOs string, searchedArch string, value any) (string, string, string, string, error) {
	object, _ := value.(map[string]any)
	builds, ok := object["builds"].([]any)
	shaFileName, ok2 := object["shasums"].(string)
	shaSigFileName, ok3 := object["shasums_signature"].(string)
	if !ok || !ok2 || !ok3 {
		return "", "", "", "", apimsg.ErrReturn
	}

	for _, build := range builds {
		object, _ = build.(map[string]any)
		osStr, ok := object["os"].(string)
		archStr, ok2 := object["arch"].(string)
		downloadURL, ok3 := object["url"].(string)
		fileName, ok4 := object["filename"].(string)
		if !ok || !ok2 || !ok3 || !ok4 {
			return "", "", "", "", apimsg.ErrReturn
		}

		if osStr != searchedOs || archStr != searchedArch {
			continue
		}

		return fileName, downloadURL, shaFileName, shaSigFileName, nil
	}

	return "", "", "", "", apimsg.ErrAsset
}

func ExtractReleases(value any) ([]string, error) {
	object, _ := value.(map[string]any)
	versions, ok := object["versions"].(map[string]any)
	if !ok {
		return nil, apimsg.ErrReturn
	}

	releases := make([]string, 0, len(object))
	for version := range versions {
		releases = append(releases, version)
	}

	return releases, nil
}
