package github

import (
	"context"
	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

func GetClient() (*github.Client, error) {
	ctx := context.Background()

	// Set up authentication using a personal access token
	//ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, nil)

	client := github.NewClient(tc)

	return client, nil
}
