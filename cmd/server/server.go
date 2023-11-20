package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

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
	// TODO test
	// TODO authenticate
	w.Write([]byte(s.Schedule.String()))
}

func (s *Server) GetCurrentTask(w http.ResponseWriter, r *http.Request) {
	// TODO test
	// TODO authenticate
	current, idx := s.Schedule.Tasks.GetTaskAtTime(time.Now())
	if current == nil {
		http.Error(w, "No current task found.", http.StatusNotFound)
		return
	}
	if idx != s.Schedule.CurrentID {
		err := s.Schedule.UpdateCurrentTask()
		if err != nil {
			log.Printf("GetCurrentTask: %s", err.Error())
			http.Error(w, "Encountered an internal server error.", http.StatusInternalServerError)
			return
		}
	}
	msg, err := json.Marshal(struct{
		Description string `json:"Description"`
		Tag string `json:"Tag"`
		Until string `json:"Until"`
	}{
		Description: current.Description,
		Tag: current.Tag,
		Until: current.EndTime.Format(time.TimeOnly),
	})
	if err != nil {
		log.Printf("GetCurrentTask: %s", err.Error())
		http.Error(w, "Encountered an internal server error.", http.StatusInternalServerError)
		return
	}
	w.Write(msg)
}

func (s *Server) ChangeCurrentTask(w http.ResponseWriter, r *http.Request) {
	// TODO test
	// TODO authenticate
	var taskModel struct{
		Description string `json:"Description"`
		Tag string `json:"Tag"`
		Until string `json:"Until"`
	}
	err := json.NewDecoder(r.Body).Decode(&taskModel)
	if err != nil {
		log.Printf("ChangeCurrentTask: %s", err.Error())
		http.Error(w, "Invalid HTTP Body", http.StatusBadRequest)
		return
	}
	// TODO validate time 
	end, err := time.Parse(time.TimeOnly, taskModel.Until)
	if err != nil {
		log.Printf("ChangeCurrentTask: %s", err)
		http.Error(w, fmt.Sprintf("Please give the time in the following format: %s", time.TimeOnly), http.StatusBadRequest)
		return
	}
	err = s.Schedule.ChangeCurrentTaskUntil(taskModel.Description, taskModel.Tag, end)
	if err != nil {
		log.Printf("ChangeCurrentTask: %s", err.Error())
		if errors.Is(err, tc.InvalidTimeError{}) {
			http.Error(w, "Please give a valid time for the task to finish.", http.StatusBadRequest)
		} else {
			http.Error(w, "Encountered an internal server error.", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) PlanSchedule(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}
