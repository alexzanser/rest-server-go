package server

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexzanser/rest-server-go.git/taskstore"
	"github.com/go-chi/chi/v5"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store: store}
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) CreateTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling create task at %s\n", req.URL.Path)

	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	type RespondeId struct {
		Id int `json:"id"`
	}

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	var rt RequestTask

	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	renderJSON(w, RespondeId{Id: id})
}

func (ts *taskServer) GetAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all task at %s\n", req.URL.Path)

	tasks := ts.store.GetAllTasks()
	renderJSON(w, tasks)

}

func (ts *taskServer) DeleteAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete all tasks at %s\n", req.URL.Path)
	ts.store.DeleteAllTasks()
}

func (ts *taskServer) DeleteTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete task at %s\n", req.URL.Path)

	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (ts *taskServer) GetTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get task at %s\n", req.URL.Path)

	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	renderJSON(w, task)
}

func (ts *taskServer) TaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/task/" {
		if req.Method == http.MethodPost {
			ts.CreateTaskHandler(w, req)
		} else if req.Method == http.MethodGet {
			ts.GetAllTasksHandler(w, req)
		} else if req.Method == http.MethodDelete {
			ts.DeleteAllTasksHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /task/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")

		if len(pathParts) < 2 {
			http.Error(w, "expect /task/<id> in task handler", http.StatusBadRequest)
			return
		}

		if req.Method == http.MethodDelete {
			ts.DeleteTaskHandler(w, req)
		} else if req.Method == http.MethodGet {
			ts.GetTaskHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expected method GET or DELETE at /task/<id> got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (ts *taskServer) TagHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling tasks by tag at %s\n", req.URL.Path)

	tag := chi.URLParam(req, "tagname")
	if tag == "" {
		http.Error(w, "expect /tag/<tagname> in tag handler", http.StatusBadRequest)
	}

	tasks := ts.store.GetTasksByTag(tag)
	renderJSON(w, tasks)
}

func (ts *taskServer) DueHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling tasks by due at %s\n", req.URL.Path)

	badRequestError := func() {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", req.URL.Path), http.StatusBadRequest)
	}

	year, err := strconv.Atoi(chi.URLParam(req, "yy"))
	if err != nil {
		badRequestError()
	}

	month, err := strconv.Atoi(chi.URLParam(req, "mm"))
	if err != nil {
		badRequestError()
	}

	day, err := strconv.Atoi(chi.URLParam(req, "dd"))
	if err != nil {
		badRequestError()
	}

	tasks := ts.store.GetTasksByDueDate(year, time.Month(month), day)
	renderJSON(w, tasks)
}
