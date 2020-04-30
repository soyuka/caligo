package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.etcd.io/etcd/v3/clientv3"

	"github.com/soyuka/caligo/handlers"
	"github.com/soyuka/caligo/storage"
)

func getConfig() storage.Config {

	etcdUrl := os.Getenv("ETCD_URL")
	endpoints := []string{"localhost:2379"}

	if etcdUrl != "" {
		endpoints = strings.Split(etcdUrl, ",")
	}

	log.Println(endpoints)

	dialTimeout, _ := time.ParseDuration(os.Getenv("ETCD_DIAL_TIMEOUT"))

	if dialTimeout == 0 {
		dialTimeout = 5 * time.Second
	}

	shortenerHostname := os.Getenv("SHORTENER_HOSTNAME")

	if shortenerHostname == "" {
		shortenerHostname = "localhost:8080"
	}

	return storage.Config{
		ShortenerHostname: shortenerHostname,
		Etcd: clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: dialTimeout,
		},
	}
}

func main() {
	config := getConfig()
	index, err := ioutil.ReadFile("./index.html")

	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HtmlFile(index)).Methods("GET")
	r.HandleFunc("/favicon.ico", handlers.Favicon)
	r.HandleFunc("/{key}", handlers.Redirect(config)).Methods("GET")
	r.HandleFunc("/", handlers.CreateLink(config)).Methods("POST")
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
