package main

import (
	"fmt"
	"net/http"
	"time"

	tc "github.com/dethancosta/timecop/internal"
	"github.com/gorilla/mux"
)

func main() {
	sched, err := tc.BuildFromFile("schedule.csv")
	now := time.Now()
	// TODO delete this line
	fmt.Println("Here: " + now.Format(time.UnixDate))
	_, offset := now.Zone()
	fmt.Printf("Offset: %d\n", offset/3600)
	for i := range sched.Tasks {
		sched.Tasks[i].StartTime = sched.Tasks[i].StartTime.AddDate(now.Year(), int(now.Month()-1), now.Day()-1)
		sched.Tasks[i].EndTime = sched.Tasks[i].EndTime.AddDate(now.Year(), int(now.Month()-1), now.Day()-1)
	}
	fmt.Println(sched.Tasks.String())
	err = sched.UpdateCurrentTask()
	if err != nil {
		panic(err)
	}
	// TODO delete
	fmt.Printf("Current idx: %d\n", sched.CurrentID)
	fmt.Println(sched.String())

	if err != nil {
		panic(err)
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

	fmt.Printf("Running on %s\n", DefaultPort)
	http.ListenAndServe("localhost:" + DefaultPort, router)
}

