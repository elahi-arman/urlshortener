package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"

	"github.com/elahi-arman/urlshortener"
)

// PostLinkHandler retrieves data about the given link
func PostLinkHandler(w http.ResponseWriter, r *http.Request) {
	address := fmt.Sprintf("%s:%s", "127.0.0.1", "6379")
	decoder := json.NewDecoder(r.Body)

	var l urlshortener.Link
	err := decoder.Decode(&l)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	l.Scope = "global"
	l.User = "arman"

	c, err := redis.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Redis Dial Error %e", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer c.Close()

	l.DateCreated = time.Now().Unix()
	l.DateModified = time.Now().Unix()
	l.Visits = 0

	fmt.Printf("%d %+v\n", time.Now().Unix(), l)

	if l.Commit(c) {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// GetLinkHandler retrieves data about the given link
func GetLinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := fmt.Sprintf("%s:%s", "127.0.0.1", "6379")

	c, err := redis.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Redis Dial Error %e", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer c.Close()

	l, err := urlshortener.GetLink(c, vars["scope"], "arman", vars["name"])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Received an error while retrieving value from Redis")
	}

	// res, err := json.Marshal(l)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not unmarshal link internally")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v", l)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/link/{scope}/{name}", GetLinkHandler).Methods("GET")
	r.HandleFunc("/link/{name}", PostLinkHandler).Methods("POST")
	http.ListenAndServe(":3001", r)
}
