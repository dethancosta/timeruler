package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const DefaultPort = 6576

func main() {

	s := Server{
		Owner: "",
		Addr:  "",
	}

	// TODO ensure port value is valid
	var port int
	flag.IntVar(&port,"p", DefaultPort, "The port that the server will run on")
	flag.Parse()
	portStr := strconv.Itoa(port)

	router := mux.NewRouter()
	// TODO update API to use
	// - GET /schedule instead of /get, POST /schedule instead of /build, PUT /schedule instead of /update
	// - GET /current instead of /current, POST /current instead of /change_current
	router.Handle("/get", http.HandlerFunc(s.GetSchedule))
	router.Handle("/build", http.HandlerFunc(s.BuildSchedule))
	router.Handle("/current", http.HandlerFunc(s.GetCurrentTask))
	router.Handle("/change_current", http.HandlerFunc(s.ChangeCurrentTask))
	router.Handle("/update", http.HandlerFunc(s.UpdateTasks))

	fmt.Printf("Running on %s\n", portStr)
	http.ListenAndServe("localhost:"+portStr, router)
}
