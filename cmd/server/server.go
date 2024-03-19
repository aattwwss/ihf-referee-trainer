package main

import (
	"log"
	"net/http"
)

func main() {
	// Define the directory containing your HTML, CSS, and JS files
	dir := "./public"

	// Create a file server to serve static files from the directory
	fileServer := http.FileServer(http.Dir(dir + "/static"))

	// Handle requests to /static/ using the file server
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Handle the root URL ("/") by serving an HTML file (e.g., index.html)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, dir+"/index.html")
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		log.Printf("Form submitted: %s", r.Form)
	})

	// Set up and start the HTTP server on port 8080
	port := "8080"
	log.Printf("Server is listening on :%s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
