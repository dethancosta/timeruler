package main

import (
	"fmt"
	"net/http"

	"github.com/dethancosta/timecop/internal"
	"github.com/gorilla/mux"
)

func main() {
	schedule, err := internal.BuildFromFile("schedule.csv")
	if err != nil {
		panic(err)
	}

	s := Server{
		Owner: "",
		Addr: "",
		AOFPath: "",
		Schedule: schedule,
	}

	router := mux.NewRouter()
	router.Handle("/get", http.HandlerFunc(s.GetSchedule))

	fmt.Printf("Running on %s\n", DefaultPort)
	http.ListenAndServe("localhost:" + DefaultPort, router)
}

