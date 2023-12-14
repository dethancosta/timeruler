package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const DefaultPort = 6576

var Address = "127.0.0.1"

func main() {

	s := Server{
		Owner: "",
		Addr:  "",
	}

	// TODO ensure port value is valid
	var port int
	var standalone bool
	flag.IntVar(&port, "p", DefaultPort, "The port that the server will run on")
	flag.BoolVar(&standalone, "sa", false, "Whether or not the server is run locally (StandAlone)")
	flag.Parse()
	portStr := strconv.Itoa(port)

	if standalone {
		err := SetPid(Address, portStr)
		if err != nil {
			panic(err)
		}
		//log.SetOutput(io.Discard)
	}

	router := mux.NewRouter()
	// TODO update API to use
	// - GET /schedule instead of /get, POST /schedule instead of /build, PUT /schedule instead of /update
	// - GET /current instead of /current, POST /current instead of /change_current
	router.Handle("/get", http.HandlerFunc(s.GetSchedule))
	router.Handle("/build", http.HandlerFunc(s.BuildSchedule))
	router.Handle("/current", http.HandlerFunc(s.GetCurrentTask))
	router.Handle("/change_current", http.HandlerFunc(s.ChangeCurrentTask))
	router.Handle("/update", http.HandlerFunc(s.UpdateTasks))

	log.Printf("Running on %s\n", portStr)
	err := http.ListenAndServe(Address+":"+portStr, router)
	if err != nil {
		panic(err)
	}
	log.Printf("Stopped running on %s\n", portStr)
}
