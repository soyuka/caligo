package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"unicode/utf8"

	"github.com/soyuka/caligo/config"
	"github.com/soyuka/caligo/id_generator"
	"github.com/soyuka/caligo/storage"
)

// GET ?http://link
func CreateLink(config config.Config) func(http.ResponseWriter, *http.Request, string) {
	return func(w http.ResponseWriter, r *http.Request, inputUrl string) {
		id, err := id_generator.GetId(config.IdAlphabet, config.IdLength)

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

		http.Redirect(w, r, fmt.Sprintf("%s/%s", config.ShortenerHostname, id), 302)
	}
}

/// Redirects short link to url
func Redirect(config config.Config) func(http.ResponseWriter, *http.Request, string) {
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
	decoded, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAACh0lEQVQ4jXWTX0hTURzHT1IPkXP3/K6gZAWRkDNjhY3MvWkUBYJPFfYWEb1HoAUpFuUUFNKXXoLewqAQKUPP2ZRB7prpRi8jH5xE/kmZm7W267273x70TrfyB+flcL6f3+98+f4YKyinM8g5Fw+I5DSpMsVJ6JxEVOHyeXHxh+rC93mlKKKZk1wlVeJ/h5MwOckexl4c+Efs5PIuJ5ElVUItlWhq+oL2jjk8ap9DQ+NUHkjhYpixwP6cuISPXuUkTVIlvN4QIpEN7C7LsjAy8hOVlRM703Dh25YPOTiJRVIlztdNIpHYBADE0xZmlk3MLJuIpy0AQDT6Gye2IZzkZklJoJIpJO7Z1NHRVQDA7IqJbk3Hw/fzaH0XhU/LYHrJBAAMDi7t8kT2ME4yRKrE0WPjMIws1jMWuqd03Ox8CZe7Fi53La7cboVPy2A1lYVpWqhyBW1AhJEqUvbfASC8YuLpZBqnPd4cwOWuxf03XxFaNAAA12/M2lP8YZz8m6RK1F2YzI3/LJTBmfqGPEDb0LccoKUlnPOBcZIxUiXKyv3Y2DAQT1vwaTru9A+jxuNF9VkPrrX1o0vTsZLKwrIsnPN8sgFRxkm8tk0ZGIgBACZ/GOjSdDwJJtE5HkeXpiP4fat7OJzcyQPJXqYootm+qDgSgKatAwAWklmImAERM7CQzAIA0mkTly5/3uquyozTOXGcMdZRxLmYtSGHKwLo7ZvH2pqeC5KuZ+H3r6Hx4k4inXysNZfEQ8pHN6ny1+64qqUSVa4gTtUEUVbuL4zyK8bYvrxdcDhkPSexsNcibZtmKCQfM9ZRtMc+Dh5UFP8tIvmWq3KOk0xwEglOMqLQWJ/D4T9ZqPgLEisPAet87nEAAAAASUVORK5CYII=")
	w.Header().Set("Content-Type", "image/x-icon")
	w.Write(decoded)
}
