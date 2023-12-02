package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	tr "github.com/dethancosta/timeruler/internal"
)

const (
	DefaultPort = "6576"
)

type Server struct {
	Owner    string // TODO replace with actual credentials for auth
	Addr     string // Address of the server
	Schedule *tr.Schedule
}

type TaskModel struct {
	Description string `json:"Description"`
	Tag         string `json:"Tag"`
	Until       string `json:"Until"`
}

func (s *Server) GetSchedule(w http.ResponseWriter, r *http.Request) {
	// TODO test
	// TODO authenticate
	if s.Schedule == nil {
		http.Error(w, "No schedule has been built yet.", http.StatusNotFound)
		return
	}
	_, idx := s.Schedule.Tasks.GetTaskAtTime(time.Now())
	if idx != s.Schedule.CurrentID {
		err := s.Schedule.UpdateCurrentTask()
		if err != nil {
			log.Printf("GetSchedule: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body := make(map[string]string)
	body["Schedule"] = s.Schedule.String()
	msg, err := json.Marshal(body)
	if err != nil {
		log.Printf("GetSchedule: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(msg)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	msg, err := json.Marshal(map[string]struct {
		Description string `json:"Description"`
		Tag         string `json:"Tag"`
		Until       string `json:"Until"`
	}{
		"Task": {
			Description: current.Description,
			Tag:         current.Tag,
			Until:       current.EndTime.Format(time.TimeOnly),
		},
	})
	if err != nil {
		log.Printf("GetCurrentTask: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(msg)
}

func (s *Server) ChangeCurrentTask(w http.ResponseWriter, r *http.Request) {
	// TODO test
	// TODO authenticate
	var taskModel TaskModel
	err := json.NewDecoder(r.Body).Decode(&taskModel)
	if err != nil {
		log.Printf("ChangeCurrentTask: %s", err.Error())
		http.Error(w, "Invalid HTTP Body", http.StatusBadRequest)
		return
	}

	// TODO validate time
	end, err := time.Parse(time.TimeOnly, taskModel.Until)
	now := time.Now()
	end = time.Date(now.Year(), now.Month(), now.Day(), end.Hour(), end.Minute(), 0, 0, time.Local)
	if err != nil {
		log.Printf("ChangeCurrentTask: %s", err)
		http.Error(w, fmt.Sprintf("Please give the time in the following format: %s", time.TimeOnly), http.StatusBadRequest)
		return
	}
	err = s.Schedule.ChangeCurrentTaskUntil(taskModel.Description, taskModel.Tag, end)
	if err != nil {
		log.Printf("ChangeCurrentTask: %s", err.Error())
		if errors.Is(err, tr.InvalidTimeError{}) {
			http.Error(w, "Please give a valid time for the task to finish.", http.StatusBadRequest)
		} else {
			http.Error(w, "Encountered an internal server error.", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) UpdateTasks(w http.ResponseWriter, r *http.Request) {

	if s.Schedule == nil {
		http.Error(w, "No schedule has been built yet.", http.StatusBadRequest)
		return
	}

	var tasks []tr.Task

	err := json.NewDecoder(r.Body).Decode(&tasks)
	if err != nil {
		log.Printf("UpdateTasks: %s", err.Error())
		http.Error(w, "Invalid HTTP Body", http.StatusBadRequest)
		return
	}

	now := time.Now()
	for _, t := range tasks {
		if t.StartTime.Before(now) {
			http.Error(w, "A task cannot start before the current time", http.StatusBadRequest)
			return
		}
	}

	err = s.Schedule.UpdateTimeBlock(tasks...)
	// TODO check type of error and return appropriate response
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) BuildSchedule(w http.ResponseWriter, r *http.Request) {
	// TODO authenticate
	// TODO test
	if s.Schedule != nil {
		http.Error(w, "Today's schedule has already been built.", http.StatusBadRequest)
		return
	}
	err := r.ParseMultipartForm(8 << 20) // max file size 8 MB
	if err != nil {
		log.Printf("BuildSchedule: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buildFile, h, err := r.FormFile("buildFile")
	if err != nil {
		log.Printf("BuildSchedule: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tmpfile, err := os.CreateTemp("./", h.Filename)
	defer func() {
		tmpfile.Close()
		err = os.Remove(tmpfile.Name())
		if err != nil {
			log.Printf("BuildSchedule: %s", err.Error())
		}
	}()
	if err != nil {
		log.Printf("BuildSchedule: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, buildFile)
	if err != nil {
		log.Printf("BuildSchedule: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.Schedule, err = tr.BuildFromFile(tmpfile.Name())
	if err != nil {
		log.Printf("BuildSchedule: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) PlanSchedule(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}
