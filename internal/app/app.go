package app

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jxskiss/base62"
	"github.com/lib/pq"
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

	mux.Handle("/assets/",
		http.StripPrefix("/assets/",
			http.FileServer(http.Dir("assets")),
		),
	)

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
	u, err := url.Parse(rawUrl)
	if err != nil {
		http.Error(w, "URL is invalid", http.StatusBadRequest)
		return
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		http.Error(w, "Schema is not http or https", http.StatusBadRequest)
		return
	}

	slug, err := a.storeURL(u, rawUrl)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}
	short := "http://" + r.Host + "/s/" + slug
	tmpl, _ := template.ParseFiles("assets/newurl.html")
	tmpl.Execute(w, short)
}

func (a *App) storeURL(u *url.URL, rawUrl string) (string, error) {
	const q = `
		INSERT INTO short_links (slug,original_url) 
		VALUES ($1,$2) 
		ON CONFLICT (original_url) DO UPDATE 
		SET original_url = short_links.original_url 
		RETURNING slug
	`

	for attempts := 0; attempts < 10; attempts++ {
		slug := a.hash(u)
		var out string
		err := a.storage.QueryRow(q, slug, rawUrl).Scan(&out)
		if err == nil {
			return out, nil
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			if pqErr.Constraint == "short_links_slug_key" {
				continue
			}
		}
		return "", fmt.Errorf("storeURL: insert into short_links: %w", err)
	}

	return "", fmt.Errorf("storeURL: generate unique slug after several attempts")
}

func (a *App) getOriginalUrl(slug string) (origUrl string, err error) {
	res := a.storage.QueryRow("SELECT original_url FROM short_links WHERE slug = $1", slug)
	err = res.Scan(&origUrl)
	if err != nil {
		return "", fmt.Errorf("failed to get url by slug %s", err.Error())
	}
	return
}

func (a *App) hash(url *url.URL) string {
	salt := strconv.FormatInt(time.Now().UnixNano(), 36)
	hash := sha256.Sum256([]byte(url.Hostname()[:1] + url.RequestURI()[:1] + salt))
	slug := base62.StdEncoding.EncodeToString(hash[:8])
	return slug[:8]
}
