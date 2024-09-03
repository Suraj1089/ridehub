package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v64/github"
	"github.com/spf13/viper"
)

type Config struct {
	Token string `mapstructure:"Token"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

// ""
type RiderPull struct {
	ID     int64    `json:"id,omitempty"`
	State  string   `json:"state,omitempty"`
	Title  string   `json:"title,omitempty"`
	Labels []string `json:"labels,omitempty"`
	URL    string   `json:"url,omitempty"`
}

func RideHubClient() *github.Client {
	config, _ := LoadConfig(".")
	return github.NewClient(nil).WithAuthToken(config.Token)
}

func ListPullRequests(ctx context.Context, owner string, repoName string, options *github.PullRequestListOptions) ([]RiderPull, error) {
	rider := RideHubClient()
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

func main() {
	fmt.Println("Welcome to ridehub")
	ctx := context.Background()
	pullOptions := &github.PullRequestListOptions{State: "open"}
	pulls, err := ListPullRequests(ctx, "Suraj1089", "SPPU-Result-Convertor", pullOptions)
	if err != nil {
		fmt.Println(err)
	}
	for _, pull := range pulls {
		fmt.Println(pull)
	}
}
