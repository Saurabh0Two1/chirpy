package main

import (
	"net/http"
)

func HealthCheckHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("ContentType", "text/plain; charset=utf-8")
	resp.WriteHeader(200)
	resp.Write([]byte("OK \n"))
}

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// )

// // plainTextHandler writes raw text data directly using Write()
// func plainTextHandler(w http.ResponseWriter, r *http.Request) {
// 	// 1. Explicitly convert string to []byte to use w.Write()
// 	message := []byte("Hello, World!")
// 	w.Write(message)
// }

// // alternativeWriteHandler shows wrapper methods that implicitely call w.Write()
// func alternativeWriteHandler(w http.ResponseWriter, r *http.Request) {
// 	// w implements the io.Writer interface, so helper functions work perfectly
// 	io.WriteString(w, "Writing a string directly!\n")
// 	fmt.Fprintf(w, "Formatted data: %s\n", "Hello")
// }

// // jsonHandler sets custom headers before calling Write()
// func jsonHandler(w http.ResponseWriter, r *http.Request) {
// 	data := map[string]string{"status": "success", "message": "Data retrieved"}

// 	// Convert your map or struct into raw bytes
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// ALWAYS set headers BEFORE calling w.Write()
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	// Send the payload JSON bytes
// 	w.Write(jsonData)
// }

// func main() {
// 	http.HandleFunc("/text", plainTextHandler)
// 	http.HandleFunc("/alt", alternativeWriteHandler)
// 	http.HandleFunc("/json", jsonHandler)

// 	fmt.Println("Server starting on :8080...")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
