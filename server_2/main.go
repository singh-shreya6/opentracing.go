package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"log"
	"go-opentracing"
	"time"
)

func handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	serverSpan,r:= go_opentracing.Deserialize(r,"server-2")
	defer serverSpan.Finish()

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := string(reqBody)
	time.Sleep(500 * time.Millisecond)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello %s!", name)
}


func main() {
	go_opentracing.Init("server-2")
	s := &http.Server{
		Addr:    ":8081",
		Handler: go_opentracing.HttpMiddleware("server-2", http.HandlerFunc(handle)),
	}
	log.Fatal(s.ListenAndServe())
}