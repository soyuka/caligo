package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"unicode/utf8"

	"github.com/soyuka/caligo/id_generator"
	"github.com/soyuka/caligo/storage"
)

// GET ?http://link
func CreateLink(config storage.Config) func(http.ResponseWriter, *http.Request, string) {
	return func(w http.ResponseWriter, r *http.Request, inputUrl string) {
		id, err := id_generator.GetId(config.IdLength)

		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}

		if utf8.RuneCountInString(inputUrl) > 2000 {
			w.WriteHeader(http.StatusRequestURITooLong)
			w.Write([]byte(http.StatusText(http.StatusRequestURITooLong)))
			return
		}

		parsedUrl, err := url.Parse(inputUrl)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))
			return
		}

		if parsedUrl.Scheme == "" {
			parsedUrl.Scheme = "http"
		}

		err = storage.Write(config.Etcd, id, parsedUrl.String())

		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s/%s", config.ShortenerHostname, id)
	}
}

/// Redirects short link to url
func Redirect(config storage.Config) func(http.ResponseWriter, *http.Request, string) {
	return func(w http.ResponseWriter, r *http.Request, key string) {
		url, err := storage.Read(config.Etcd, key)

		if url == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusText(http.StatusNotFound)))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			log.Print(err)
			return
		}

		http.Redirect(w, r, string(url), 301)
	}
}

/// Black favicon just for fun
func Favicon(w http.ResponseWriter, r *http.Request) {
	decoded, _ := base64.StdEncoding.DecodeString("AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAQbvwAApLaQACnuAADuTwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACIgAAIiAAACETIAIxEgAAI0QgAkQyAAACJDIjQiAAACNEMiNEMgACMUQyI0QTIAIxRDIjRBMgAjFEIiJEEyACMTICICMTIAAiICACAiIAAAAAIAIAAAAAAAIAACAAAAACIAAAAiAAAAAAAAAAAAD//wAA//8AAOPHAADBgwAAwYMAAOAHAADAAwAAgAEAAIABAACAAQAAgkEAAMWjAAD9vwAA+98AAOfnAAD//wAA")
	w.Header().Set("Content-Type", "image/x-icon")
	w.Write(decoded)
}
