package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Aritra640/Short-Kings/server/urlshort"
	"github.com/gorilla/mux"
)

//build a crud api to get websites and their short forms
//store them some where
//do the url shortening

type entry struct {
	Id      string `json:"id"`
	Webpath string `json:"webpath"`
	Weburl  string `json:"weburl"`
}

var entries []entry

func main() {

	r := mux.NewRouter()

	entries = append(entries, entry{
		Id:      "1",
		Webpath: "https://www.this_is_test.com",
		Weburl:  "https://www.google.com",
	})

	r.HandleFunc("/urlShort", getEntries).Methods("GET")
	r.HandleFunc("/urlShort/{id}", getEntriesByID).Methods("GET")
	r.HandleFunc("/urlShort", createEntry).Methods("POST")
	r.HandleFunc("/urlShort/{id}", updateEntry).Methods("PUT")
	r.HandleFunc("/urlShort/{id}", deleteEntry).Methods("DELETE")

	log.Println("Starting port at localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	fallback := defaultMux()

	serveURL(fallback)
}

func getEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func getEntriesByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range entries {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func deleteEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range entries {
		if item.Id == params["id"] {
			entries = append(entries[:index], entries[index+1:]...)
		}
	}

	json.NewEncoder(w).Encode(entries)
}

func createEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newEntry entry
	_ = json.NewDecoder(r.Body).Decode(&newEntry)

	entries = append(entries, newEntry)
	json.NewEncoder(w).Encode(entries)
}

func updateEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range entries {
		if item.Id == params["id"] {
			entries = append(entries[:index], entries[index+1:]...)

			var newEntry entry
			_ = json.NewDecoder(r.Body).Decode(&newEntry)
			newEntry.Id = params["id"]

			entries = append(entries, newEntry)
			json.NewEncoder(w).Encode(entries)

			return
		}
	}
}

func defaultMux() *http.ServeMux { //use gorilla and replace it later
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("default condition triggered")
}

func serveURL(fallback *http.ServeMux) {
	pathsToUrls := map[string]string{}

	for _, item := range entries {
		pathsToUrls[item.Webpath] = item.Weburl
	}

	mapHandler := urlshort.MapHandler(pathsToUrls, fallback)

	//log.Println("Starting Server at :8080\n")
	http.ListenAndServe(":8080", mapHandler)
}
