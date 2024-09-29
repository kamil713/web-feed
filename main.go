package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

// Handler for  the main page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

// Function to fetch and parse the RSS feed
func fetchRSSFeed(feedURL string) (*RSS, error) {
	// Fetch the feed data
	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, fmt.Errorf("could not fetch RSS feed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %v", err)
	}

	// Parse the XML into our RSS struct
	var rss RSS
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		return nil, fmt.Errorf("could not parse RSS feed: %v", err)
	}

	return &rss, nil
}

// Fetch RSS Handler: Fetches and parses the RSS feed, then returns the items
func fetchRSSHandler(w http.ResponseWriter, r *http.Request) {
	feedURL := r.URL.Query().Get("feed_url")
	if feedURL == "" {
		http.Error(w, "Feed URL is required", http.StatusBadRequest)
		return
	}

	// Fetch and parse the RSS feed
	rssFeed, err := fetchRSSFeed(feedURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch RSS feed: %v", err), http.StatusInternalServerError)
		return
	}

	// Render the parsed feed items in the HTML template
	tmpl.ExecuteTemplate(w, "rssItems", rssFeed.Channel.Items)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/fetch-rss", fetchRSSHandler)

	// Serve static files (e.g., CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
