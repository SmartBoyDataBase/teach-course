package main

import (
	"log"
	"net/http"
	"sbdb-teach-course/handler"
)

func main() {
	//http.HandleFunc goes here
	http.HandleFunc("/ping", handler.PingPongHandler)
	http.HandleFunc("/teach-course", handler.Handler)
	http.HandleFunc("/teach-courses", handler.AllHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
