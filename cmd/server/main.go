package main

import (
	"fmt"
	"net/http"

	tr "github.com/dethancosta/timeruler/internal"
	"github.com/gorilla/mux"
)

func main() {
	sched, err := tr.BuildFromFile("schedule.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	s := Server{
		Owner: "",
		Addr: "",
		AOFPath: "",
		Schedule: sched,
	}

	router := mux.NewRouter()
	router.Handle("/get", http.HandlerFunc(s.GetSchedule))
	router.Handle("/build", http.HandlerFunc(s.BuildSchedule))
	router.Handle("/current_task", http.HandlerFunc(s.GetCurrentTask))
	router.Handle("/change_current", http.HandlerFunc(s.ChangeCurrentTask))

	fmt.Printf("Running on %s\n", DefaultPort)
	http.ListenAndServe("localhost:" + DefaultPort, router)
}

