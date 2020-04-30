package main

import (
	"log"
	"net/http"

	"strings"

	"github.com/soyuka/caligo/config"
	"github.com/soyuka/caligo/handlers"
)

func main() {
	config := config.GetConfig()
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

	log.Println("Listen", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
