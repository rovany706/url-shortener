package handlers

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"go.uber.org/zap"

	"github.com/rovany706/url-shortener/internal/app"
	"github.com/rovany706/url-shortener/internal/auth"
	"github.com/rovany706/url-shortener/internal/config"
	"github.com/rovany706/url-shortener/internal/repository"
	"github.com/rovany706/url-shortener/internal/service"
)

type exampleRepository struct {
	repository.Repository
}

type exampleURLShortener struct {
	app.URLShortener
}

type exampleTokenManager struct {
	auth.TokenManager
}

type exampleDeleteService struct {
	service.DeleteService
}

func ExamplePingHandler() {
	repository := new(exampleRepository)
	logger := zap.NewNop()
	handler := PingHandler(repository, logger)

	// Example of registering handler:
	http.HandleFunc("/ping", handler)

	// Example of sending request:
	resp, err := http.Get("http://service:8080/ping")
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()
}

func ExampleRedirectHandlers_RedirectHandler() {
	app := new(exampleURLShortener)
	redirectHandlers := NewRedirectHandlers(app)
	handler := redirectHandlers.RedirectHandler()

	// Example of registering handler:
	http.HandleFunc("/{id}", handler)

	// Example of sending request:
	resp, err := http.Get("http://service:8080/488575e6")
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()
}

func ExampleShortenURLHandlers_MakeShortURLHandler() {
	app := new(exampleURLShortener)
	repository := new(exampleRepository)
	tokenManager := new(exampleTokenManager)
	logger := zap.NewNop()
	appConfig := config.NewConfig()

	shortenHandlers := NewShortenURLHandlers(app, tokenManager, repository, appConfig, logger)
	handler := shortenHandlers.MakeShortURLHandler()

	// Example of registering handler:
	http.HandleFunc("/", handler)

	// Example of sending request:
	requestBody := "http://example.com"
	resp, err := http.Post("http://service:8080/", "text/plain", strings.NewReader(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()
}

func ExampleShortenURLHandlers_MakeShortURLHandlerJSON() {
	app := new(exampleURLShortener)
	repository := new(exampleRepository)
	tokenManager := new(exampleTokenManager)
	logger := zap.NewNop()
	appConfig := config.NewConfig()

	shortenHandlers := NewShortenURLHandlers(app, tokenManager, repository, appConfig, logger)
	handler := shortenHandlers.MakeShortURLHandlerJSON()

	// Example of registering handler:
	http.HandleFunc("/api/shorten", handler)

	// Example of sending request:
	requestBody := `{"url": "https://practicum.yandex.ru"}`
	resp, err := http.Post("http://service:8080/api/shorten", "application/json", strings.NewReader(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()
}

func ExampleShortenURLHandlers_MakeShortURLBatchHandler() {
	app := new(exampleURLShortener)
	repository := new(exampleRepository)
	tokenManager := new(exampleTokenManager)
	logger := zap.NewNop()
	appConfig := config.NewConfig()

	shortenHandlers := NewShortenURLHandlers(app, tokenManager, repository, appConfig, logger)
	handler := shortenHandlers.MakeShortURLBatchHandler()

	// Example of registering handler:
	http.HandleFunc("/api/shorten/batch", handler)

	// Example of sending request:
	requestBody := `
[
  {
    "correlation_id": "d9200816-793a-469c-bf04-976754db63ca",
    "original_url": "https://www.uuidgenerator.net/guid1"
  },
  {
    "correlation_id": "d9200816-793a-469c-bf04-976754db63ca",
    "original_url": "https://www.google.com/"
  },
  {
    "correlation_id": "d9200816-793a-469c-bf04-976754db63ca",
    "original_url": "http://example.com/1"
  }
]
	`
	resp, err := http.Post("http://service:8080/api/shorten/batch", "application/json", strings.NewReader(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()
}

func ExampleUserHandlers_GetUserURLsHandler() {
	deleteService := new(exampleDeleteService)
	repository := new(exampleRepository)
	tokenManager := new(exampleTokenManager)
	logger := zap.NewNop()
	appConfig := config.NewConfig()

	shortenHandlers := NewUserHandlers(deleteService, tokenManager, repository, appConfig, logger)
	handler := shortenHandlers.GetUserURLsHandler()

	// Example of registering handler:
	http.HandleFunc("/api/user/urls", handler)

	// Example of sending request:
	jar, _ := cookiejar.New(nil)
	cookie := &http.Cookie{
		Name:   "token",
		Value:  "<jwt-token>",
		Path:   "/",
		Domain: "service:8080",
	}

	u, _ := url.Parse("http://service:8080")
	jar.SetCookies(u, []*http.Cookie{cookie})
	client := &http.Client{
		Jar: jar,
	}

	resp, err := client.Get("http://service:8080/api/user/urls")
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()
}

func ExampleUserHandlers_DeleteUserURLsHandler() {
	deleteService := new(exampleDeleteService)
	repository := new(exampleRepository)
	tokenManager := new(exampleTokenManager)
	logger := zap.NewNop()
	appConfig := config.NewConfig()

	shortenHandlers := NewUserHandlers(deleteService, tokenManager, repository, appConfig, logger)
	handler := shortenHandlers.DeleteUserURLsHandler()

	// Example of registering handler:
	http.HandleFunc("/api/user/urls", handler)

	// Example of sending request:
	jar, _ := cookiejar.New(nil)
	cookie := &http.Cookie{
		Name:   "token",
		Value:  "<jwt-token>",
		Path:   "/",
		Domain: "service:8080",
	}

	u, err := url.Parse("http://service:8080")
	if err != nil {
		log.Fatal(err)
	}

	jar.SetCookies(u, []*http.Cookie{cookie})
	client := &http.Client{
		Jar: jar,
	}
	requestBody := `["67b00967", "595c3cce"]`

	request, _ := http.NewRequest(http.MethodDelete, "http://service:8080/api/user/urls", strings.NewReader(requestBody))
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()
}
