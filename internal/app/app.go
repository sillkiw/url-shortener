package app

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jxskiss/base62"
)

type App struct {
	storage *sql.DB
}

func New(storage *sql.DB) App {
	return App{storage}
}

func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/shorten", a.shorten)
	mux.HandleFunc("/s/", a.redirect)

	return mux
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "assets/index.html")
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/s/")
	if slug == "" || strings.Contains(slug, "/") {
		http.NotFound(w, r)
		return
	}
	url, err := a.getOriginalUrl(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, 302)

}

func (a *App) shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}
	rawUrl := r.FormValue("url")
	if rawUrl == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	u, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		http.Error(w, "URL is invalid", http.StatusBadRequest)
		return
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		http.Error(w, "Schema is not http or https", http.StatusBadRequest)
		return
	}

	salt := strconv.FormatInt(time.Now().UnixNano(), 36)
	hash := sha256.Sum256([]byte(u.Hostname()[:1] + u.RequestURI()[:1] + salt))
	shortHash := base62.StdEncoding.EncodeToString(hash[:8])
	short, err := a.storeURL(rawUrl, shortHash)
	if err != nil {
		log.Printf("failed to store in db %s", err.Error())
		http.Error(w, "Failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>slug : %s", "http://localhost:4000/s/"+short)

}

func (a *App) storeURL(url string, short string) (string, error) {
	res := a.storage.QueryRow("INSERT INTO short_links (slug,original_url) VALUES ($1,$2) ON CONFLICT (original_url) DO UPDATE SET original_url = short_links.original_url RETURNING slug",
		short, url)
	err := res.Scan(&short)
	if err != nil {
		return "", fmt.Errorf("failed to store data %s", err.Error())
	}

	return short, nil

}

func (a *App) getOriginalUrl(slug string) (origUrl string, err error) {
	res := a.storage.QueryRow("SELECT original_url FROM short_links WHERE slug = $1", slug)
	err = res.Scan(&origUrl)
	if err != nil {
		return "", fmt.Errorf("failed to get url by slug %s", err.Error())
	}
	return
}
