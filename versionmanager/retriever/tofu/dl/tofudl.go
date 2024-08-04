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

package tofudlmirroring

import (
	"strings"
	"text/template"

	"github.com/tofuutils/tenv/v3/pkg/apimsg"
)

type artifactDesc struct {
	Artifact string
	Version  string
}

type URLBuilder struct {
	t *template.Template
	v string
}

func MakeURLBuilder(templateURL string, version string) (URLBuilder, error) {
	t, err := template.New("").Parse(templateURL)

	return URLBuilder{t: t, v: version}, err
}

func (b URLBuilder) Build(artifactName string) (string, error) {
	var builder strings.Builder
	err := b.t.Execute(&builder, artifactDesc{
		Artifact: artifactName,
		Version:  b.v,
	})
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func ExtractReleases(value any) ([]string, error) {
	object, _ := value.(map[string]any)
	versions, ok := object["versions"].([]any)
	if !ok {
		return nil, apimsg.ErrReturn
	}

	releases := make([]string, 0, len(object))
	for _, versionDesc := range versions {
		castedVersionDesc, _ := versionDesc.(map[string]any)
		versionID := castedVersionDesc["id"]
		version, ok := versionID.(string)
		if !ok {
			return nil, apimsg.ErrReturn
		}

		releases = append(releases, version)
	}

	return releases, nil
}
