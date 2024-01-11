package common

import (
	"context"
	"github.com/google/go-github/v58/github"
)

func ListRemote(client *github.Client, owner, repo string) ([]string, error) {
	ctx := context.Background()

	gitHubReleases, _, err := client.Repositories.ListReleases(ctx, owner, repo, nil)
	if err != nil {
		return nil, err
	}

	var releases []string
	for _, release := range gitHubReleases {
		releases = append(releases, *release.TagName)
	}

	return releases, nil
}
