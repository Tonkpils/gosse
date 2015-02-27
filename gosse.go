package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type GoSSE struct {
	clients  map[chan string]bool
	messages chan string
}

func (sse *GoSSE) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Stream is not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	msgChan := make(chan string)
	sse.clients[msgChan] = true

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		delete(sse.clients, msgChan)
		log.Println("Client disconnected.")
	}()

	go sse.ping(w, f)

	// Send retry interval
	fmt.Fprintf(w, "retry: 1000\n")
	f.Flush()

	// Await for messages
	for {
		fmt.Fprintf(w, "data: %s\n\n", <-msgChan)
		f.Flush()
	}

	log.Println("Finished SSE request")
}

func (sse *GoSSE) ping(w http.ResponseWriter, f http.Flusher) {
	for {
		fmt.Fprintf(w, ":\n\n")
		f.Flush()
		time.Sleep(10 * time.Second)
	}
}

func (sse *GoSSE) Listen() {
	go func() {
		for {
			select {
			case msg := <-sse.messages:
				for client, _ := range sse.clients {
					client <- msg
				}
			}
		}
	}()
}

func HomePage(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, nil)
	log.Printf("Completed HTTP request at %s", req.URL.Path)
}

func main() {
	gosse := &GoSSE{
		clients:  make(map[chan string]bool),
		messages: make(chan string),
	}

	gosse.Listen()

	// Goroutine to generate events
	go func() {
		for {
			gosse.messages <- fmt.Sprintf(`{"time": %q}`, time.Now())
			time.Sleep(30 * time.Second)
		}
	}()

	http.Handle("/", http.HandlerFunc(HomePage))
	http.Handle("/events", gosse)

	http.ListenAndServe(":8080", nil)
}
