package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type FeedItem struct {
	Title       string
	Link        string
	Description string
	PubDate     string
}

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

// Handler for  the main page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", homeHandler)

	// Serve static files (e.g., CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
