package main

import (
	"fmt"
	"net/http"

	tc "github.com/dethancosta/timecop/internal"
	"github.com/gorilla/mux"
)

func main() {
	sched, err := tc.BuildFromFile("schedule.csv")
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
	router.Handle("/current_task", http.HandlerFunc(s.GetCurrentTask))
	router.Handle("/change_current", http.HandlerFunc(s.ChangeCurrentTask))

	fmt.Printf("Running on %s\n", DefaultPort)
	http.ListenAndServe("localhost:" + DefaultPort, router)
}

