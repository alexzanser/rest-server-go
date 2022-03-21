package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	server "github.com/alexzanser/rest-server-go.git/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	server := server.NewTaskServer()
	//Context will be used in other processes
	_, cancel := context.WithCancel(context.Background())

	r := chi.NewRouter()
	
	r.Post("/task/", server.CreateTaskHandler)
	r.Get("/task/", server.GetAllTasksHandler)
	r.Get("/task/{id}", server.GetTaskHandler)
	r.Delete("/task/{id}", server.DeleteTaskHandler)
	r.Get("/tag/{tagname}", server.TagHandler)
	r.Get("/due/{yy}/{mm}/{dd}", server.DueHandler)

	go func() {
		log.Fatal(http.ListenAndServe("localhost:"+ os.Args[1], r))
	}()

	//Gracefull shutdown
	stopAppCh := make(chan struct{})
	<- gracefullShutDown(server, stopAppCh, cancel)
}

func gracefullShutDown(server *server.TaskServer,stopAppCh chan struct{}, cancel context.CancelFunc) <- chan struct{}{
	sigquit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)
	
	go func() <- chan struct{} {
		log.Println("Captured signal: ", <- sigquit)
		log.Println("Gracefully shutting down server ...")
		cancel()

		if err := server.Shutdown(context.Background()); err != nil {
			log.Println("Can't shutdown main server: ", err.Error())
		}
		stopAppCh <- struct{}{}
		return stopAppCh
	}()
	return stopAppCh
}
