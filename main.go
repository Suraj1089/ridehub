package main

import (
	"context"
	"fmt"

	ridehub "github.com/Suraj1089/ridehub/pkg"
	"github.com/google/go-github/v64/github"
	"github.com/joho/godotenv"
)


func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	rider := ridehub.GetRiderClient()
	newPR := &github.NewPullRequest{
		Title:               github.String("Update file Name"),
		Head:                github.String("main"),
		Base:                github.String("ridehub"),
		Body:                github.String("Body of pr"),
		MaintainerCanModify: github.Bool(true),
	}
	pr, err := rider.CreatePull(ctx, "Suraj1089", "SPPU-Result-Convertor", newPR)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pr)
}
