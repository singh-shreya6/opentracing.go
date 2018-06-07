package main

import (
	"fmt"
	"net/http"
	"time"
	"io/ioutil"
	"strings"
	"context"
	"go-opentracing"
	"log"
)

func main() {
	go_opentracing.Init("server-1")
	s:=&http.Server{
		Addr:    ":8082",
		Handler: go_opentracing.HttpMiddleware("server-1", http.HandlerFunc(handle)),
	}
	log.Fatal(s.ListenAndServe())

}


func handle(w http.ResponseWriter, r *http.Request) {
	t:=time.Now()
	name:=r.URL.Path[1:]
	time.Sleep(250 * time.Millisecond)
	makeRequest(r.Context(), "http://localhost:8081",name)
	fmt.Fprintf(w, "Hello %s! time is %s", name,t.Format("Mon Jan _2 15:04:05 2006"))
}


func makeRequest(ctx context.Context, serverAddr string, name string){
	span, ctx:=go_opentracing.Introduce_span(ctx,"server-1_to_server-2")
	defer span.Finish()

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8081", strings.NewReader(name))
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return
	}

	go_opentracing.Serialise(ctx,req)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() error: %v\n", err)
		return
	}
	fmt.Printf("\n%s", string(data))
}


