package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alexzanser/rest-server-go.git/server"
)

func main() {
	mux := http.NewServeMux()
	server := server.NewTaskServer()
	
	mux.HandleFunc("/task/", server.TaskHandler)
	mux.HandleFunc("/tag/", server.TagHandler)
	mux.HandleFunc("/due/", server.DueHandler)

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux))
}