package save

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/sillkiw/url-shorten/internal/lib/aliasgen"
	apierrors "github.com/sillkiw/url-shorten/internal/lib/api/errors"
	resp "github.com/sillkiw/url-shorten/internal/lib/api/response"
	"github.com/sillkiw/url-shorten/internal/lib/validation"
	"github.com/sillkiw/url-shorten/internal/storage"
)

type Request struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Body
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (id int64, err error)
}

func New(log *slog.Logger, urlSaver URLSaver, v validation.Validator) http.HandlerFunc {
	baseLog := log
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log := baseLog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.Any("err", err))
			resp.Write(w, r, http.StatusBadRequest, resp.Error("invalid json"))
			return
		}

		log.Debug("request body decode", slog.Any("request", req))

		// validate url
		req.URL = strings.TrimSpace(req.URL)
		if err := v.URL(req.URL); err != nil {
			log.Info("url is invalid",
				slog.String("url", req.URL),
				slog.Any("err", err),
			)
			msg := apierrors.URLValidation(err, v.MaxURLLen)
			resp.Write(w, r, http.StatusBadRequest, resp.Error(msg))
			return
		}

		// validate alias
		req.Alias = strings.TrimSpace(req.Alias)
		genByServer := false
		if req.Alias != "" {
			if err := v.Alias(req.Alias); err != nil {
				log.Info("alias is invalid",
					slog.String("alias", req.Alias),
					slog.Any("err", err),
				)
				msg := apierrors.AliasValidation(err, v.MinAliasLen, v.MaxAliasLen)
				resp.Write(w, r, http.StatusBadRequest, resp.Error(msg))
				return
			}
		} else {
			genByServer = true
			req.Alias = aliasgen.GenRandomString(v.DefaultAliasLen)
		}

		// save url
		id, err := urlSaver.SaveURL(req.URL, req.Alias)
		if err != nil {
			// if alias already exist, generate another one
			if genByServer && errors.Is(err, storage.ErrAliasExists) {
				id, req.Alias, err = saveWithGeneratedAlias(urlSaver, req.URL, v.DefaultAliasLen, 10)
				if err == nil {
					log.Info("url added", slog.Int64("id", id))
					resp.Write(w, r, http.StatusOK,
						Response{
							Body:  resp.OK(),
							Alias: req.Alias,
						},
					)
					return
				}
			}

			log.Error("failed to add url",
				slog.String("url", req.URL),
				slog.Any("err", err),
			)
			code, msg := apierrors.Storage(err)
			resp.Write(w, r, code, resp.Error(msg))
			return
		}

		log.Info("url added",
			slog.Int64("id", id),
			slog.String("url", req.URL),
		)
		resp.Write(w, r, http.StatusCreated,
			Response{
				Body:  resp.OK(),
				Alias: req.Alias,
			},
		)

	}
}

// saveWithGeneratedAlias try @attempts times to save url with genereated alias
var ErrAliasGenFailed = errors.New("failed to generate unique alias")

func saveWithGeneratedAlias(urlSaver URLSaver, urlToSave string, aliasLen int, attempts int) (id int64, alias string, err error) {
	const op = "SaveWithGeneratedAlias"

	for i := 0; i < attempts; i++ {
		alias = aliasgen.GenRandomString(aliasLen)
		id, err = urlSaver.SaveURL(urlToSave, alias)
		if err == nil {
			return id, alias, nil
		}
		if errors.Is(err, storage.ErrAliasExists) {
			continue
		}
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}
	return 0, "", fmt.Errorf("%s: %w", op, ErrAliasGenFailed)
}
