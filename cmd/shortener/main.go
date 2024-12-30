package main

import (
	"io"
	"net/http"

	"github.com/rovany706/url-shortener/internal/app"
)

const BaseURL = "http://localhost:8080/"

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	app := app.URLShortenerApp{}
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(MainHook(&app)))

	return http.ListenAndServe(":8080", mux)
}

func MainHook(app *app.URLShortenerApp) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			MakeShortURLHandler(app, w, r)
		case http.MethodGet:
			RedirectHandler(app, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func RedirectHandler(app *app.URLShortenerApp, w http.ResponseWriter, r *http.Request) {
	shortID := r.URL.Path[1:] // убирает слеш (возможна ошибка, лучше переписать через rune)
	if fullURL, ok := app.GetFullURL(shortID); ok {
		http.Redirect(w, r, fullURL, http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	}
}

func MakeShortURLHandler(app *app.URLShortenerApp, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	shortID, err := app.GetShortID(string(body))

	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(BaseURL + shortID))
}
