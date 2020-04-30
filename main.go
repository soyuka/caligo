package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

	dialTimeout, _ := time.ParseDuration(os.Getenv("ETCD_DIAL_TIMEOUT"))

	if dialTimeout == 0 {
		dialTimeout = 5 * time.Second
	}

	shortenerHostname := os.Getenv("CALIGO_HOSTNAME")

	if shortenerHostname == "" {
		shortenerHostname = "http://localhost:8080"
	}

	log.Println("ETCD endpoints", endpoints)
	log.Println("Hostname", shortenerHostname)

	idLength, err := strconv.ParseInt(os.Getenv("CALIGO_ID_LENGTH"), 10, 32)

	if idLength == 0 || err != nil {
		idLength = 12
	}

	return storage.Config{
		ShortenerHostname: shortenerHostname,
		IdLength:          int(idLength),
		Etcd: clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: dialTimeout,
		},
	}
}

func main() {
	config := getConfig()
	createLinkHandler := handlers.CreateLink(config)
	redirectHandler := handlers.Redirect(config)
	cookieName := "created"

	http.HandleFunc("/favicon.ico", handlers.Favicon)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := strings.Replace(r.URL.RawQuery, "?", "", 1)

		if url != "" {
			cookie := &http.Cookie{Name: cookieName}
			http.SetCookie(w, cookie)
			createLinkHandler(w, r, url)
			return
		}

		key := strings.Replace(r.URL.Path, "/", "", 1)

		if key != "" {
			_, err := r.Cookie(cookieName)

			if err == nil {
				cookie := &http.Cookie{Name: cookieName, MaxAge: -1}
				http.SetCookie(w, cookie)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(http.StatusText(http.StatusCreated)))
				return
			}

			redirectHandler(w, r, key)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	})

	port := os.Getenv("CALIGO_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listen", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
