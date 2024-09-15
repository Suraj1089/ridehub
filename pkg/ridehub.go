package ridehub

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v64/github"
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
	GitService   *github.GitService
}

func GetRiderClient() *RiderClient {
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_SECRET_TOKEN"))
	return &RiderClient{Issues: client.Issues, PullRequests: client.PullRequests, Users: client.Users, Teams: client.Teams}
}

func (rider *RiderClient) RiderPullRequests(ctx context.Context, owner string, repoName string, options *github.PullRequestListOptions) ([]RiderPull, error) {
	pulls, _, err := rider.PullRequests.List(ctx, owner, repoName, options)
	if err != nil {
		return nil, err
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
		return nil, err
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

func (rider *RiderClient) AddLabelToPull(ctx context.Context, owner string, repoName string, pullNumber int, labelsToAdd []string) ([]string, error) {
	_, _, err := rider.PullRequests.Get(ctx, owner, repoName, pullNumber)
	if err != nil {
		return nil, err
	}
	labels, _, err := rider.Issues.AddLabelsToIssue(ctx, owner, repoName, pullNumber, labelsToAdd)
	if err != nil {
		return nil, err
	}
	var allLabelsInPull []string
	for _, val := range labels {
		allLabelsInPull = append(allLabelsInPull, *val.Name)
	}
	return allLabelsInPull, nil
}

func (rider *RiderClient) GetPullLabels(ctx context.Context, owner string, repoName string, pullNumber int) ([]string, error) {
	pull, _, err := rider.PullRequests.Get(ctx, owner, repoName, pullNumber)
	if err != nil {
		return nil, err
	}
	var labels []string
	for _, val := range pull.Labels {
		labels = append(labels, *val.Name)
	}
	return labels, nil
}

func (rider *RiderClient) RemoveLabelFromPull(ctx context.Context, owner string, repoName string, pullNumber int, label string) ([]string, error) {
	_, err := rider.Issues.RemoveLabelForIssue(ctx, owner, repoName, pullNumber, label)
	if err != nil {
		return nil, err
	}
	labels, err := rider.GetPullLabels(ctx, owner, repoName, pullNumber)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func (rider *RiderClient) CreatePull(ctx context.Context, owner string, repoName string, pull *github.NewPullRequest) (*int, error) {
	featureBranch := pull.Base
	fmt.Println(*featureBranch)
	ref, _, err := rider.GitService.GetRef(ctx, owner, repoName, *featureBranch)
	fmt.Println("ref is ")
	if err != nil {
		fmt.Println("returning err")
		return nil, err
	}
	fmt.Printf("ref = %s", ref)
	if err != nil {
		featureBranchUrl, err := rider.CreateBranch(ctx, owner, repoName, *pull.Head, *featureBranch)
		if err != nil {
			return nil, err
		}
		fmt.Println(featureBranchUrl)
	}
	return nil, nil
	// pr, _, err := rider.PullRequests.Create(ctx, owner, repoName, pull)
	// if err != nil {
	// 	return nil, err
	// }
	// return pr.Number, nil
}

func (rider *RiderClient) CreateBranch(ctx context.Context, owner string, repoName string, baseBranch string, featureBranch string) (*string, error) {
	ref := "refs/heads/" + baseBranch
	baseRef, _, err := rider.GitService.GetRef(ctx, owner, repoName, ref)
	if err != nil {
		return nil, nil
	}
	reference := &github.Reference{
		Ref:    github.String("refs/heads/" + baseBranch),
		Object: &github.GitObject{SHA: baseRef.Object.SHA},
	}
	newRef, _, err := rider.GitService.CreateRef(ctx, owner, repoName, reference)

	if err != nil {
		return nil, err
	}
	return newRef.URL, nil
}

type RefService struct {
	GitService *github.GitService
}

func (b *RefService) GetRef(ctx context.Context, owner string, repo string, baseRefName string, featureRefName string) (*github.Reference, error) {
	ref := "refs/heads/" + baseRefName
	baseRef, _, err := b.GitService.GetRef(ctx, owner, repo, ref)
	if err != nil {
		return nil, err
	}
	return baseRef, nil
}

func (b *RefService) CreateRef(ctx context.Context, owner string, repo string, baseRefName string, featureRefName string) (*github.Reference, error) {
	baseRef, err := b.GetRef(ctx, owner, repo, baseRefName, featureRefName)
	if err != nil {
		return nil, err
	}
	reference := &github.Reference{
		Ref:    github.String("refs/heads/" + baseRefName),
		Object: &github.GitObject{SHA: baseRef.Object.SHA},
	}
	newRef, _, err := b.GitService.CreateRef(ctx, owner, repo, reference)
	if err != nil {
		return nil, err
	}
	return newRef, nil
}
