package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	// This is the main entry point for the application.
	// You can initialize your application here, set up routes, etc.
	// For example, you might want to start a web server or connect to a database.

	// Example: Start a web server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}
