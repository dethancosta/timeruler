package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	tc "github.com/dethancosta/timecop/internal"
	"github.com/gorilla/mux"
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
	// TODO authenticate
	taskId := mux.Vars(r)["taskId"]
	if taskId == "" {
		http.Error(w, "No Task Id given.", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(taskId)
	if err != nil {
		log.Printf("RemoveTask: %s", err.Error())
		http.Error(w, "Task Id must be a valid number", http.StatusBadRequest)
		return
	}

	err = s.Schedule.RemoveTask(id)
	if err != nil {
		log.Printf("RemoveTask: %s", err.Error())
		if errors.Is(err, tc.IndexOutOfBoundsError{}) {
			http.Error(w, "Please give a valid index", http.StatusBadRequest)
		} else if errors.Is(err, tc.InvalidScheduleError{}) {
			http.Error(w, "Operation not allowed on this schedule", http.StatusBadRequest)
		} else {
			http.Error(w, "Encountered an internal server error.", http.StatusInternalServerError)
		}
		return
	}
	w.Write([]byte(s.Schedule.String()))
}

func (s *Server) ChangeCurrentTask(w http.ResponseWriter, r *http.Request) {
	// TODO implement
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
	w.Write([]byte(s.Schedule.String()))
}

func (s *Server) PlanSchedule(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}
