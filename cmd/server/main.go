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
	// TODO delete
	fmt.Printf("Current idx: %d\n", sched.CurrentID)
	fmt.Println(sched.String())

	s := Server{
		Owner: "",
		Addr: "",
		AOFPath: "",
		Schedule: sched,
	}

	router := mux.NewRouter()
	router.Handle("/get", http.HandlerFunc(s.GetSchedule))
	router.Handle("/current_task", http.HandlerFunc(s.GetCurrentTask))

	fmt.Printf("Running on %s\n", DefaultPort)
	http.ListenAndServe("localhost:" + DefaultPort, router)
}

