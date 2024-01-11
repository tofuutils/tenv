package common

import (
	"context"
	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

func ListRemote(owner, repo string) ([]string, error) {
	tc := oauth2.NewClient(ctx, nil)
	githubClient := github.NewClient(tc)

	// Get releases for the specified repository
	gitHubReleases, _, err := githubClient.Repositories.ListReleases(context.Background(), owner, repo, nil)
	if err != nil {
		return nil, err
	}

	var releases []string
	for _, release := range gitHubReleases {
		releases = append(releases, *release.TagName)
	}

	return releases, nil
}

func List() {
}

func Use() {
}

func Install() {
}

func Uninstall() {
}

func VersionName() {
}

func Init() {
}

func Pin() {
}
