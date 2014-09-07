package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

func main() {
	log.Print("Starting...")
	//http.HandleFunc("/", HomeHandler) // homepage
	serveSingle("/", "./index.html")
	http.HandleFunc("/upload.go", upload)
	// Mandatory root-based resources
	serveSingle("/sitemap.xml", "./sitemap.xml")
	serveSingle("/favicon.ico", "./favicon.ico")
	serveSingle("/robots.txt", "./robots.txt")

	// Normal resources
	http.Handle("/static", http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":9999", nil)
}

func upload(w http.ResponseWriter, r *http.Request) {
	// why do I need to parse the query manually?
	uri := strings.SplitN(r.RequestURI, "?", 2) // gimme two parts. part 2 is the query.
	if len(uri) != 2 {
		log.Print("upload.go accessed with no query string")
		return
	}
	v, err := url.ParseQuery(uri[1])
	if err != nil {
		log.Print(err)
	}
	t := v["time"]
	p := v["placeid"]
	log.Printf("%s::%s", t, p)
	w.Write([]byte(fmt.Sprintf("submitted: %+v %+v", p, t)))
}
