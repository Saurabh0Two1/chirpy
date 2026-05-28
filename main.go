package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()

	httpServer := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))

	// Example to serve to a url different from the directory folder names
	// dir := http.Dir("./assets/")
	// mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(dir)))

	log.Printf("Serving on port: %s\n", port)
	httpServer.ListenAndServe()
}
