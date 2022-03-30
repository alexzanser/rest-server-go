package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	task "stdlib_chi_postgres/internal/handlers"

	"github.com/go-chi/chi"
)

func main() {
	server := task.NewTaskServer()
	//Context will be used in other processes
	_, cancel := context.WithCancel(context.Background())

	r := chi.NewRouter()
	server.Handler = r
	server.Addr = ":" + os.Getenv("SERVERPORT")
	log.Println(server.Addr)
	r.Post("/task/", server.CreateTaskHandler)
	r.Get("/task/", server.GetAllTasksHandler)
	r.Get("/task/{id}", server.GetTaskHandler)
	r.Delete("/task/{id}", server.DeleteTaskHandler)
	r.Get("/tag/{tagname}", server.TagHandler)
	r.Get("/due/{yy}/{mm}/{dd}", server.DueHandler)

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	//Gracefull shutdown
	stopAppCh := make(chan struct{})
	<- gracefullShutDown(server, stopAppCh, cancel)
}

func gracefullShutDown(server *task.TaskServer,stopAppCh chan struct{}, cancel context.CancelFunc) <- chan struct{}{
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
