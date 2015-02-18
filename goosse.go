package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type GooSSE struct {
	clients  map[chan string]bool
	messages chan string
}

func (sse *GooSSE) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Stream is not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")

	msgChan := make(chan string)
	sse.clients[msgChan] = true

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		delete(sse.clients, msgChan)
		log.Println("Client disconnected.")
	}()

	go sse.ping(w, flusher)

	// Send retry interval
	fmt.Fprintf(w, "retry: 1000\n")
	flusher.Flush()

	// Await for messages
	for {
		msg := <-msgChan
		fmt.Fprintf(w, "data: %s\n\n", msg)
		flusher.Flush()
	}

	log.Println("Finished SSE request")
}

func (sse *GooSSE) ping(w http.ResponseWriter, f http.Flusher) {
	for {
		fmt.Fprintf(w, ":\n\n")
		f.Flush()
		time.Sleep(10 * time.Second)
	}
}

func (sse *GooSSE) Listen() {
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
	goosse := &GooSSE{
		clients:  make(map[chan string]bool),
		messages: make(chan string),
	}

	goosse.Listen()

	// Goroutine to generate events
	go func() {
		for {
			goosse.messages <- fmt.Sprintf(`{"time": %q}`, time.Now())
			time.Sleep(30 * time.Second)
		}
	}()

	http.Handle("/", http.HandlerFunc(HomePage))
	http.Handle("/events", goosse)

	http.ListenAndServe(":3000", nil)
}
