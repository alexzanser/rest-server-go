package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexzanser/rest-server-go.git/server"
	"github.com/go-chi/chi/v5"
)


func main() {
	//Context will be used in other processes
	_, cancel := context.WithCancel(context.Background())

	server := server.NewTaskServer()
	r := chi.NewRouter()

	r.Post("/task/", server.CreateTaskHandler)
	r.Get("/task/", server.GetAllTasksHandler)
	r.Get("/task/{id}", server.GetTaskHandler)
	r.Delete("/task/{id}", server.DeleteTaskHandler)
	r.Get("/tag/{tagname}", server.TagHandler)
	r.Get("/due/{yy}/{mm}/{dd}", server.DueHandler)

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), r))


	//Gracefull shutdown
	sigquit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)
	stopAppCh := make(chan struct{})

	go func() {
		log.Println("Captured signal: ", <- sigquit)
		log.Println("Gracefully shutting down server ...")

		cancel()

		if err := server.Shutdown(context.Background()); err != nil {
			log.Println("Can't shutdown main server: ", err.Error())
		}
		stopAppCh <- struct{}{}
	}()

	<- stopAppCh
}

// func main() {
// 	mux := http.NewServeMux()
// 	server := server.NewTaskServer()
	
// 	mux.HandleFunc("/task/", server.TaskHandler)
// 	mux.HandleFunc("/tag/", server.TagHandler)
// 	mux.HandleFunc("/due/", server.DueHandler)

// log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux))
// }