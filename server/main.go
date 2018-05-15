package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/elahi-arman/urlshortener"
)

// PostLinkHandler retrieves data about the given link
func PostLinkHandler(w http.ResponseWriter, r *http.Request) {
	var l urlshortener.Link

	address := fmt.Sprintf("%s:%s", "127.0.0.1", "6379")
	log.Debug("74530::Received a POST, decoding body")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&l)
	if err != nil {
		log.Error(fmt.Sprintf("74531::Could not decode post body %e", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if l.Scope == "" {
		l.Scope = "global"
	}

	l.User = "arman"

	log.Debug("74532::Connecting to Redis")
	c, err := redis.Dial("tcp", address)
	if err != nil {
		log.Error(fmt.Sprintf("74533::Could not connect to redis %e", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer c.Close()

	l.DateCreated = time.Now().Unix()
	l.DateModified = time.Now().Unix()
	l.Visits = 0

	if l.Commit(c) {
		log.Info("74534::Committed link 201")
		w.WriteHeader(http.StatusCreated)
	} else {
		log.Warn("74535::Could not commit link 400")
		w.WriteHeader(http.StatusBadRequest)
	}
}

// GetLinkHandler retrieves data about the given link
func GetLinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := fmt.Sprintf("%s:%s", "127.0.0.1", "6379")

	log.Debug("15460::Instantiating Redis Connection")
	c, err := redis.Dial("tcp", address)
	if err != nil {
		log.Error(fmt.Sprintf("15461::Couldn't connect to Redis %e", err))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Couldn't contact data store")
		return
	}
	defer c.Close()

	log.Debug(fmt.Sprintf("15462::Searching redis for %s:arman:%s", vars["scope"], vars["name"]))
	l, err := urlshortener.GetLink(c, vars["scope"], "arman", vars["name"])

	if err != nil {
		log.Error(fmt.Sprintf("15463::Couldn't retrieve value from Redis %e", err))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Received an error while retrieving value from data store")
		return
	}

	log.Debug(fmt.Sprintf("15464::Received link, updating visit count"))
	l.Visit(c)
	log.Info(fmt.Sprintf("15465::Sending 301 %#v", l))
	http.Redirect(w, r, l.Link, 301)
}

// GetGlobalLinkHandler retrieves data about the given link from the global scope
func GetGlobalLinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := fmt.Sprintf("%s:%s", "127.0.0.1", "6379")

	log.Debug("53170::Instantiating Redis Connection")
	c, err := redis.Dial("tcp", address)
	if err != nil {
		log.Error(fmt.Sprintf("53171::Couldn't connect to Redis %e", err))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Couldn't contact data store")
		return
	}
	defer c.Close()

	log.Debug(fmt.Sprintf("53172::Searching redis for global:arman:%s", vars["name"]))
	l, err := urlshortener.GetLink(c, "global", "arman", vars["name"])

	if err != nil {
		log.Error(fmt.Sprintf("53173::Couldn't retrieve value from Redis %e", err))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Received an error while retrieving value from data store")
		return
	}

	log.Debug(fmt.Sprintf("53174::Received link, updating visit count"))
	l.Visit(c)
	log.Info(fmt.Sprintf("53175::Sending 301 %#v", l))
	http.Redirect(w, r, l.Link, 301)
}

func main() {
	Formatter := new(log.TextFormatter)
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	log.SetLevel(log.DebugLevel)
	r := mux.NewRouter()
	r.HandleFunc("/link/{scope}/{name}", GetLinkHandler).Methods("GET")
	r.HandleFunc("/link/{name}", GetGlobalLinkHandler).Methods("GET")
	r.HandleFunc("/link", PostLinkHandler).Methods("POST")
	http.ListenAndServe(":3001", r)
}
