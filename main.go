package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alexzanser/rest-server-go.git/server"
	"github.com/go-chi/chi/v5"
)


func main() {
	server := server.NewTaskServer()
	r := chi.NewRouter()
	
	r.Get("/task/", server.GetAllTasksHandler)
	r.Post("/task/", server.CreateTaskHandler)
	r.Delete("/task/{id}", server.DeleteTaskHandler)
	r.Get("/task/{id}", server.GetTaskHandler)
	r.Get("/tag/{tagname}", server.TagHandler)
	r.Get("/due/{yy}/{mm}/{dd}", server.DueHandler)

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), r))

}

// func main() {
// 	mux := http.NewServeMux()
// 	server := server.NewTaskServer()
	
// 	mux.HandleFunc("/task/", server.TaskHandler)
// 	mux.HandleFunc("/tag/", server.TagHandler)
// 	mux.HandleFunc("/due/", server.DueHandler)

// log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux))
// }