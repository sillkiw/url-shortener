package redirect

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	apierrors "github.com/sillkiw/url-shorten/internal/lib/api/errors"
	resp "github.com/sillkiw/url-shorten/internal/lib/api/response"
	"github.com/sillkiw/url-shorten/internal/lib/validation"
)

type URLGetter interface {
	GetURL(alias string) (url string, err error)
}

func New(log *slog.Logger, urlGetter URLGetter, v validation.Validator) http.HandlerFunc {
	baseLog := log
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"
		log := baseLog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		alias = strings.TrimSpace(alias)
		if alias == "" {
			log.Info("alias is empty")
			resp.Write(w, r, http.StatusBadRequest, resp.Error("alias is required"))
			return
		}
		if err := v.Alias(alias); err != nil {
			log.Info("alias is invalid", slog.Any("err", err))
			resp.Write(w, r, http.StatusBadRequest, resp.Error("invalid request"))
			return
		}

		redirectURL, err := urlGetter.GetURL(alias)
		if err != nil {
			log.Info("url not found",
				slog.String("alias", alias),
				slog.Any("err", err),
			)
			code, msg := apierrors.Storage(err)
			resp.Write(w, r, code, resp.Error(msg))
			return
		}
		log.Info("redirecting", slog.String("url", redirectURL))
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}
