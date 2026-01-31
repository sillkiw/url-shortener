package apierrors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sillkiw/url-shorten/internal/lib/validation"
	"github.com/sillkiw/url-shorten/internal/storage"
)

const (
	msgURLRequired   = "url is required"
	msgUserInfo      = "userinfo in url is not allowed"
	msgShortenMyHost = "cannot shorten our own url"
	msgInvalidURL    = "invalid url"
)

func URLValidation(err error, maxURLLen int) string {
	switch {
	case errors.Is(err, validation.ErrURLRequired):
		return msgURLRequired
	case errors.Is(err, validation.ErrUserInfo):
		return msgUserInfo
	case errors.Is(err, validation.ErrOurOwnURL):
		return msgShortenMyHost
	case errors.Is(err, validation.ErrURLToLong):
		return fmt.Sprintf("url length must be <= %d characters", maxURLLen)
	default:
		return msgInvalidURL
	}
}

const (
	msgInvalidAlias = "alias may contain only letters, digits, '_' and '-'"
)

func AliasValidation(err error, minLen, maxLen int) string {
	if errors.Is(err, validation.ErrAliasTooShort) || errors.Is(err, validation.ErrAliasTooLong) {
		return fmt.Sprintf("alias length must be between %d and %d characters", minLen, maxLen)
	}
	return msgInvalidAlias
}

const (
	msgURLExists   = "url already exists"
	msgAliasExists = "alias already exists"
	msgNotFound    = "alias not found"
	msgInternal    = "internal server error"
)

func Storage(err error) (int, string) {
	switch {
	case errors.Is(err, storage.ErrURLExists):
		return http.StatusConflict, msgURLExists
	case errors.Is(err, storage.ErrAliasExists):
		return http.StatusConflict, msgAliasExists
	case errors.Is(err, storage.ErrAliasNotFound):
		return http.StatusNotFound, msgNotFound
	default:
		return http.StatusInternalServerError, msgInternal
	}

}
