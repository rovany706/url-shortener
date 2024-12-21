package main

import (
	"io"
	"net/http"

	"github.com/rovany706/url-shortener/internal/app"
)

const baseURL = "http://localhost:8080/"

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(mainHook))

	return http.ListenAndServe(":8080", mux)
}

func mainHook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		makeShortURLHandler(w, r)
	case http.MethodGet:
		redirectHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortId := r.URL.Path[1:] // убирает слеш (возможна ошибка, лучше переписать через rune)
	if fullURL, ok := app.GetFullURL(shortId); ok {
		w.Header().Add("Location", fullURL)
		http.Redirect(w, r, fullURL, http.StatusTemporaryRedirect)
	} else {
		http.NotFound(w, r)
	}
}

func makeShortURLHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	shortId, err := app.GetShortId(string(body))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(baseURL + shortId))
}
