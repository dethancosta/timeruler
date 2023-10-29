package main

import (
	"fmt"
	"log"
	"net/http"

	tc "github.com/dethancosta/timecop/internal"
)

const (
	DefaultPort = "6576"
)

type Server struct {
	Owner    string // TODO replace with actual credentials for auth
	Addr     string
	AOFPath  string // Filepath for append-only log file
	Schedule *tc.Schedule
}

func (s *Server) GetSchedule(w http.ResponseWriter, r *http.Request) {
	// TODO authenticate
	w.Write([]byte(s.Schedule.String()))
}

func (s *Server) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}

func (s *Server) GetCurrentTask(w http.ResponseWriter, r *http.Request) {
	// TODO implement
	err := tc.SendJson(s.Schedule.CurrentTask, w)

	if err != nil {
		// TODO append to AOF
		log.Println(fmt.Errorf("GetCurrentTask: %w", err))
		w.WriteHeader(http.StatusInternalServerError) // TODO send body with message
	}
}

func (s *Server) RemoveTask(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}

func (s *Server) ChangeCurrentTask(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}

func (s *Server) PlanSchedule(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}
