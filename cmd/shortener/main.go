package main

import (
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/config"
)

var appConfig *config.AppConfig

func main() {
	var err error
	appConfig, err = config.ParseArgs(os.Args[0], os.Args[1:])

	if err != nil {
		panic(err)
	}

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	app := app.URLShortenerApp{}
	r := MainRouter(&app)

	return http.ListenAndServe(appConfig.AppRunAddress, r)
}

func MainRouter(app *app.URLShortenerApp) chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", RedirectHandler(app))
		r.Post("/", MakeShortURLHandler(app))
	})

	return r
}

func RedirectHandler(app *app.URLShortenerApp) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		shortID := chi.URLParam(r, "id")
		if fullURL, ok := app.GetFullURL(shortID); ok {
			http.Redirect(w, r, fullURL, http.StatusTemporaryRedirect)
		} else {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
		}
	}
}

func MakeShortURLHandler(app *app.URLShortenerApp) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		w.Write([]byte(appConfig.BaseURL + "/" + shortID))
	}
}
