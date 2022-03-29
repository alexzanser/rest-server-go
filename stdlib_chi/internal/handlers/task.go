package task

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"

	"stdlib_chi/internal/taskstore"
	"github.com/go-chi/chi/v5"
)

type TaskServer struct {
	http.Server
	store *taskstore.TaskStore
}

func NewTaskServer() *TaskServer {
	store := taskstore.New()
	return &TaskServer{store: store}
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *TaskServer) CreateTaskHandler(w http.ResponseWriter, req *http.Request) {
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

func (ts *TaskServer) GetAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all task at %s\n", req.URL.Path)

	tasks := ts.store.GetAllTasks()
	renderJSON(w, tasks)

}

func (ts *TaskServer) DeleteAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete all tasks at %s\n", req.URL.Path)
	ts.store.DeleteAllTasks()
}

func (ts *TaskServer) DeleteTaskHandler(w http.ResponseWriter, req *http.Request) {
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

func (ts *TaskServer) GetTaskHandler(w http.ResponseWriter, req *http.Request) {
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

func (ts *TaskServer) TagHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling tasks by tag at %s\n", req.URL.Path)

	tag := chi.URLParam(req, "tagname")
	if tag == "" {
		http.Error(w, "expect /tag/<tagname> in tag handler", http.StatusBadRequest)
	}

	tasks := ts.store.GetTasksByTag(tag)
	renderJSON(w, tasks)
}

func (ts *TaskServer) DueHandler(w http.ResponseWriter, req *http.Request) {
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
