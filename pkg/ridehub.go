package ridehub

import (
	"fmt"

	"github.com/google/go-github/v64/github"
)

const Token string = "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"

func Ridehub() *github.Client {
	client := github.NewClient(nil).WithAuthToken(Token)
	return client
}

func ListPr(username string) {
	client := github.NewClient(nil).WithAuthToken(Token)
	fmt.Println(client.PullRequests)
}