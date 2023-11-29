package main

import (
	"fmt"
	"net/http"

	//tr "github.com/dethancosta/timeruler/internal"
	"github.com/gorilla/mux"
)

func main() {

	s := Server{
		Owner: "",
		Addr: "",
		AOFPath: "",
		//Schedule: sched,
	}

	router := mux.NewRouter()
	// TODO update API to use 
	// - GET /schedule instead of /get, POST /schedule instead of /build
	// - GET /current instead of /current, POST /current instead of /change_current
	router.Handle("/get", http.HandlerFunc(s.GetSchedule))
	router.Handle("/build", http.HandlerFunc(s.BuildSchedule))
	router.Handle("/current", http.HandlerFunc(s.GetCurrentTask))
	router.Handle("/change_current", http.HandlerFunc(s.ChangeCurrentTask))

	fmt.Printf("Running on %s\n", DefaultPort)
	http.ListenAndServe("localhost:" + DefaultPort, router)
}

