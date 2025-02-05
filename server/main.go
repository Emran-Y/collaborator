package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"terminal_collab/server/handlers"
	"terminal_collab/server/helper_function"
	"terminal_collab/server/types"
)

// main is the entry point of your Go application.
func main() {
	addr := flag.String("addr", ":8080", "Server's network address")
	flag.Parse()

	// Create a new Clients instance from the types package
	clients := types.NewClients()

	mux := http.NewServeMux()
	// Pass 'clients' into the handler via closure
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleConn(w, r, clients)
	})

	// Goroutine to handle the concurrency loop (add/delete/read requests)
	go helper_function.HandleClients(clients)

	// Goroutine to handle incoming messages
	go handlers.HandleMsg(clients)

	// Goroutine to handle document syncing
	go handlers.HandleSync(clients)

	// Start the server
	log.Printf("Starting server on %s", *addr)
	server := &http.Server{
		Addr:         *addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
