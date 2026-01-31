package delete

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

type AliasDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, aliasDeleter AliasDeleter, v validation.Validator) http.HandlerFunc {
	baseLog := log
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New"

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

		if err := aliasDeleter.DeleteURL(alias); err != nil {
			log.Info("failed to delete alias",
				slog.String("alias", alias),
				slog.Any("err", err),
			)
			code, msg := apierrors.Storage(err)
			resp.Write(w, r, code, resp.Error(msg))
			return
		}
		log.Info("alias was deleted", slog.String("alias", alias))
		w.WriteHeader(http.StatusNoContent)

	}
}
