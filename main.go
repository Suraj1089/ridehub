package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v64/github"
	"github.com/joho/godotenv"
)

// ""
type RiderPull struct {
	ID     int64    `json:"id,omitempty"`
	State  string   `json:"state,omitempty"`
	Title  string   `json:"title,omitempty"`
	Labels []string `json:"labels,omitempty"`
	URL    string   `json:"url,omitempty"`
}

type RiderIssue struct {
	ID               int64                    `json:"id,omitempty"`
	State            string                   `json:"state,omitempty"`
	Title            string                   `json:"title,omitempty"`
	Labels           []*string                `json:"labels,omitempty"`
	URL              string                   `json:"url,omitempty"`
	PullRequestLinks *github.PullRequestLinks `json:"pull_request,omitempty"`
}

type RiderClient struct {
	Issues       *github.IssuesService
	PullRequests *github.PullRequestsService
	Users        *github.UsersService
	Teams        *github.TeamsService
}

func GetRiderClient() *RiderClient {
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_SECRET_TOKEN"))
	return &RiderClient{Issues: client.Issues, PullRequests: client.PullRequests, Users: client.Users, Teams: client.Teams}
}

func (rider *RiderClient) RiderPullRequests(ctx context.Context, owner string, repoName string, options *github.PullRequestListOptions) ([]RiderPull, error) {
	pulls, _, err := rider.PullRequests.List(ctx, owner, repoName, options)
	if err != nil {
		panic(err)
	}
	var riderPulls []RiderPull
	// output format id : Title : Url : state : labels
	for _, val := range pulls {
		var labels []string
		for _, label := range val.Labels {
			labels = append(labels, *label.Name)
		}
		riderPull := RiderPull{
			ID:     *val.ID,
			State:  *val.State,
			Title:  *val.Title,
			Labels: labels,
			URL:    *val.URL,
		}
		riderPulls = append(riderPulls, riderPull)
	}
	return riderPulls, err
}

func (rider *RiderClient) RiderIssues(ctx context.Context, owner string, repoName string, options *github.IssueListOptions) ([]RiderIssue, error) {
	issues, _, err := rider.Issues.List(ctx, false, options)
	if err != nil {
		fmt.Println(err)
	}
	var riderIssues []RiderIssue
	// output format id : Title : Url : state : labels
	for _, val := range issues {
		var labels []*string
		for _, label := range val.Labels {
			labels = append(labels, (label.Name))
		}
		riderIssue := RiderIssue{
			ID:               *val.ID,
			State:            *val.State,
			Title:            *val.Title,
			Labels:           labels,
			URL:              *val.URL,
			PullRequestLinks: val.PullRequestLinks,
		}
		riderIssues = append(riderIssues, riderIssue)
	}
	return riderIssues, err

}

func main() {
	fmt.Println("Welcome to ridehub")
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	rider := GetRiderClient()
	pullOptions := &github.PullRequestListOptions{State: "open"}
	pulls, err := rider.RiderPullRequests(ctx, "Suraj1089", "SPPU-Result-Convertor", pullOptions)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pulls)
	fmt.Println("Issues")
	opts := &github.IssueListOptions{Filter: "created"}
	issues, err := rider.RiderIssues(ctx, "Suraj1089", "SPPU-Result-Convertor", opts)
	fmt.Println(issues)

}
