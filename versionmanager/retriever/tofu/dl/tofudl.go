package tofudlmirroring

import (
	"strings"
	"text/template"

	"github.com/tofuutils/tenv/v2/pkg/apimsg"
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
		versionId, _ := castedVersionDesc["id"]
		version, ok := versionId.(string)
		if !ok {
			return nil, apimsg.ErrReturn
		}

		releases = append(releases, version)
	}

	return releases, nil
}
