package main

import (
	"fmt"

	"github.com/dethancosta/timecop/internal"
)

type App struct {
	Schedule *internal.Schedule
	User     string
	AuthKey  string
}

func main() {
	schedule, err := internal.BuildFromFile()
	if err != nil {
		panic(err)
	}

	fmt.Println(schedule.String())
}