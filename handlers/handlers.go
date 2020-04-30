package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/soyuka/caligo/id_generator"
	"github.com/soyuka/caligo/storage"
)

/// Creates a link on form POST / with url field
func CreateLink(config storage.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := id_generator.GetId(config.IdLength)

		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}

		formUrl := r.PostFormValue("url")

		if formUrl == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))
			return
		}

		if utf8.RuneCountInString(formUrl) > 2000 {
			w.WriteHeader(http.StatusRequestURITooLong)
			w.Write([]byte(http.StatusText(http.StatusRequestURITooLong)))
			return
		}

		parsedUrl, err := url.Parse(formUrl)

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
func Redirect(config storage.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		url, err := storage.Read(config.Etcd, vars["key"])

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

/// Serve a html file
func HtmlFile(index []byte) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(index)
	}
}
