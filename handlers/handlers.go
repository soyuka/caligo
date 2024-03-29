package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
    "strconv"
	"log"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/matoous/go-nanoid"
)

const cookieName = "caligo"

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func makeStatusError(code int) StatusError {
	return StatusError{
		Code: code,
		Err:  errors.New(http.StatusText(code)),
	}
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

type Handler struct {
	Env *Env
	Handler func(e *Env, w http.ResponseWriter, r *http.Request) error
}

// Adapted from https://blog.questionable.services/article/http-handler-error-handling-revisited/
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.Handler(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

/// Single endpoint /
/// When there's a query we're using it's value to save the link in the database
/// If there's a code (eg: hostname.com/EnYQkRXzK30d) we redirect to the given value
func GetIndex(env *Env, w http.ResponseWriter, r *http.Request) error {

	url := strings.Replace(r.URL.RawQuery, "?", "", 1)
	if (strings.HasPrefix(url, "u=")) {
		q := r.URL.Query()
		url = q.Get("u")
	}

	if url != "" {
		cookie := &http.Cookie{Name: cookieName, SameSite: http.SameSiteStrictMode, Secure: true, HttpOnly: true}
		http.SetCookie(w, cookie)
		return CreateLink(env, w, r, url)
	}

	key := strings.Replace(r.URL.Path, "/", "", 1)

	if key != "" {
		_, err := r.Cookie(cookieName)

		if err == nil {
			cookie := &http.Cookie{Name: cookieName, MaxAge: -1, SameSite: http.SameSiteStrictMode, Secure: true, HttpOnly: true}
			http.SetCookie(w, cookie)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(http.StatusText(http.StatusCreated)))
			return nil
		}

		return Redirect(env, w, r, key)
	}

	return Index(env, w, r)
}

// GET ?http://link creates the link and redirect to the link
func CreateLink(env *Env, w http.ResponseWriter, r *http.Request, inputUrl string) error {
	if utf8.RuneCountInString(inputUrl) > 2000 {
		return makeStatusError(http.StatusRequestURITooLong)
	}

	parsedUrl, err := url.Parse(inputUrl)

	if err != nil {
		return StatusError{http.StatusBadRequest, err}
	}

	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "https"
	}

	id, err := gonanoid.Generate(env.Config.IdAlphabet, env.Config.IdLength)

	if err != nil {
		return StatusError{http.StatusInternalServerError, err}
	}

	err = env.Transport.Put(id, parsedUrl.String())

	if err != nil {
		return StatusError{http.StatusInternalServerError, err}
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", env.Config.ShortenerHostname, id), 302)
	return nil
}

/// Redirects short link to url
func Redirect(env *Env, w http.ResponseWriter, r *http.Request, key string) error {
	url, err := env.Transport.Get(key)

	if err != nil {
		return StatusError{http.StatusInternalServerError, err}
	}

	if url == "" {
		return makeStatusError(http.StatusNotFound)
	}

	http.Redirect(w, r, string(url), 301)
	return nil
}

/// Favicon just for fun
func Favicon(env *Env, w http.ResponseWriter, r *http.Request) error {
	decoded, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAACh0lEQVQ4jXWTX0hTURzHT1IPkXP3/K6gZAWRkDNjhY3MvWkUBYJPFfYWEb1HoAUpFuUUFNKXXoLewqAQKUPP2ZRB7prpRi8jH5xE/kmZm7W267273x70TrfyB+flcL6f3+98+f4YKyinM8g5Fw+I5DSpMsVJ6JxEVOHyeXHxh+rC93mlKKKZk1wlVeJ/h5MwOckexl4c+Efs5PIuJ5ElVUItlWhq+oL2jjk8ap9DQ+NUHkjhYpixwP6cuISPXuUkTVIlvN4QIpEN7C7LsjAy8hOVlRM703Dh25YPOTiJRVIlztdNIpHYBADE0xZmlk3MLJuIpy0AQDT6Gye2IZzkZklJoJIpJO7Z1NHRVQDA7IqJbk3Hw/fzaH0XhU/LYHrJBAAMDi7t8kT2ME4yRKrE0WPjMIws1jMWuqd03Ox8CZe7Fi53La7cboVPy2A1lYVpWqhyBW1AhJEqUvbfASC8YuLpZBqnPd4cwOWuxf03XxFaNAAA12/M2lP8YZz8m6RK1F2YzI3/LJTBmfqGPEDb0LccoKUlnPOBcZIxUiXKyv3Y2DAQT1vwaTru9A+jxuNF9VkPrrX1o0vTsZLKwrIsnPN8sgFRxkm8tk0ZGIgBACZ/GOjSdDwJJtE5HkeXpiP4fat7OJzcyQPJXqYootm+qDgSgKatAwAWklmImAERM7CQzAIA0mkTly5/3uquyozTOXGcMdZRxLmYtSGHKwLo7ZvH2pqeC5KuZ+H3r6Hx4k4inXysNZfEQ8pHN6ny1+64qqUSVa4gTtUEUVbuL4zyK8bYvrxdcDhkPSexsNcibZtmKCQfM9ZRtMc+Dh5UFP8tIvmWq3KOk0xwEglOMqLQWJ/D4T9ZqPgLEisPAet87nEAAAAASUVORK5CYII=")
	w.Header().Set("Content-Type", "image/x-icon")
	w.Write(decoded)
	return nil
}

/// Favicon just for fun
func Index(env *Env, w http.ResponseWriter, r *http.Request) error {
	count, _ := env.Transport.Count()

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	index := `<!DOCTYPE html>
<html lang="en">
<head>
  <title>Caligo</title>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <meta name="description" content="URL obfuscator" />
  <style>
	body {margin: 5% auto; background: #f2f2f2; color: #444444; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; font-size: 16px; line-height: 1.8; text-shadow: 0 1px 0 #ffffff; max-width: 73%;}
	code {background: white;}
	a {border-bottom: 1px solid #444444; color: #444444; text-decoration: none;}
	a:hover {border-bottom: 0;}
  </style>
</head>
<body>
  <h1>Caligo</h1>
  <h2>Obfuscate an URL</h2>
  <form method="GET" action="/">
	<input type="text" name="u" />
	<input type="submit" value="Obfuscate"/>
	<p><small>Data has no warranty and can be removed at any time.</small></p>
  </form>
  <h2>API</h2>
  <p>Open <code>`+env.Config.ShortenerHostname+`?URL</code> in your browser. Copy the redirected URL from the address bar (CTRL+L, CTRL+C).</p>
  <p><a href="https://github.com/soyuka/caligo">Code on github</a></p>
  <h2>Statistics</h2>
  <p>`+strconv.FormatInt(count, 10)+` links obfuscated</p>
</body>
</html>`
	w.Write([]byte(index))
	return nil
}
