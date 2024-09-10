package main

import (
	"context"
	"fmt"
	"strings"

	ridehub "github.com/Suraj1089/ridehub/pkg"
	"github.com/joho/godotenv"
)


func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	rider := ridehub.GetRiderClient()

	pulls, err := rider.GetPullLabels(ctx, "Suraj1089", "SPPU-RESULT-CONVERTER", 24)
	if err != nil {
		fmt.Println(strings.Split(err.Error(), ":")[2])
	}
	fmt.Println(pulls)
}
