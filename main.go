package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	http.HandleFunc("/submit", submissionHandler)
	http.HandleFunc("/stats", statsHandler)

	log.Println("listening on :9999")
	http.ListenAndServe(":9999", nil)
}

type Place struct {
	ID string

	Average, Median, Percentile_90th, Slowest, Fastest float32
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		log.Println("error parsing stats url: ", err)
	}
	v, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Println("error parsing stats url: ", err)
	}

	placeIds := make([]string, 0)
	if placesUrlArray, ok := v["place_id[]"]; ok {
		placeIds = placesUrlArray
	}

	response := make([]Place, 0)
	for _, placeId := range placeIds {
		response = append(response, Place{
			ID:              placeId,
			Average:         30.5,
			Median:          25.2,
			Percentile_90th: 45.9,
			Slowest:         120.4,
			Fastest:         12.1,
		})
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("error marhsalling to json: ", err)
	}
	w.Write([]byte(jsonResponse))
}

func submissionHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	log.Println(r.RequestURI)
	log.Println(string(body))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": "created"}`))
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
