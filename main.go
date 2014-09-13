package main

// TODO: db configuration, port configuration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

var db *sql.DB

func main() {
	log.Print("Starting...")
	var err error
	db, err = sql.Open("mysql", "/test")
	if err != nil {
		log.Fatal(err)
	}

	// Mandatory root-based resources
	serveSingle("/sitemap.xml", "./sitemap.xml")
	serveSingle("/favicon.ico", "./favicon.ico")
	serveSingle("/robots.txt", "./robots.txt")

	// submit a new stat or get the stats for a list of place ids
	http.HandleFunc("/submit", submissionHandler)
	http.HandleFunc("/stats", statsHandler)

	log.Println("listening on :9999")
	http.ListenAndServe(":9999", nil)
}

type Place struct {
	ID string

	Average, Median, Percentile_90th, Slowest, Fastest float32
}

func returnErr(err error, w http.ResponseWriter) {
	log.Println("error retrieving data: ", err)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error": "unable to fetch data"}`))
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

	questionMarks := make([]string, len(placeIds))
	for i := 0; i < len(placeIds); i++ {
		questionMarks[i] = "?"
	}

	queryString := "select place_id, timer_ms from timer where place_id in (" + strings.Join(questionMarks, ",") + ")"

	asInterface := make([]interface{}, len(placeIds))

	for i, v := range placeIds {
		asInterface[i] = interface{}(v)
	}

	rows, err := db.Query(queryString, asInterface...)

	defer rows.Close()
	if err != nil {
		returnErr(err, w)
		return
	}

	collection := make(map[string][]int)

	for rows.Next() {
		var place_id string
		var time int

		if err := rows.Scan(&place_id, &time); err != nil {
			returnErr(err, w)
			return
		}
		collection[place_id] = append(collection[place_id], time)
	}

	response := make([]Place, 0)
	for _, placeId := range placeIds {
		if _, ok := collection[placeId]; ok {
			sort.Ints(collection[placeId])
			response = append(response, Place{
				ID:              placeId,
				Average:         average(collection[placeId]),
				Median:          median(collection[placeId]),
				Percentile_90th: percentile(collection[placeId], 90),
				Slowest:         max(collection[placeId]),
				Fastest:         min(collection[placeId]),
			})
		}
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("error marhsalling to json: ", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}

func average(i []int) float32 {
	sum := 0
	for _, v := range i {
		sum += v
	}
	if len(i) == 0 {
		return float32(0)
	}
	return (float32(sum) / float32(len(i))) / float32(1000)
}
func median(i []int) float32 {
	return float32(i[len(i)/2]) / float32(1000)
}
func percentile(i []int, p int) float32 {
	if len(i) == 0 {
		return float32(0)
	}
	percentileIndex := int(math.Ceil(float64(len(i)) * float64(p) / float64(100)))
	if percentileIndex >= len(i) {
		percentileIndex = len(i) - 1
	} else if percentileIndex < 0 {
		percentileIndex = 0
	}
	value := i[percentileIndex]
	return float32(value) / float32(1000)
}
func max(i []int) float32 {
	if len(i) == 0 {
		return float32(0)
	}
	return float32(i[len(i)-1]) / float32(1000)
}
func min(i []int) float32 {
	if len(i) == 0 {
		return float32(0)
	}
	return float32(i[0]) / float32(1000)
}

func submissionHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		log.Println("error parsing stats url: ", err)
	}
	v, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Println("error parsing stats url: ", err)
	}

	var placeID string
	if id, ok := v["place_id"]; ok {
		placeID = id[0]
	}
	var timeMS string
	if time, ok := v["time"]; ok {
		timeMS = time[0]
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	_, err = db.Exec("insert into timer set place_id=?, timer_ms=?, ip=INET_ATON(?)", placeID, timeMS, ip)
	if err != nil {
		log.Println("error inserting stat: ", err)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "insert error"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": "created"}`))
}
