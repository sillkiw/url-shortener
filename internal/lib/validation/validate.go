package validation

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type Validator struct {
	MyHost          string
	MaxURLLen       int
	MinAliasLen     int
	MaxAliasLen     int
	DefaultAliasLen int
}

func New(myHost string, maxURLLen, minAliasLen, maxAliasLen, defaultAliasLen int) Validator {
	return Validator{MyHost: myHost,
		MaxURLLen:       maxURLLen,
		MinAliasLen:     minAliasLen,
		MaxAliasLen:     maxAliasLen,
		DefaultAliasLen: defaultAliasLen,
	}
}

// URL validation errors
var (
	ErrURLRequired = errors.New("url is required")
	ErrInvalidURL  = errors.New("invalid url")
	ErrURLToLong   = errors.New("url is longer than available")
	ErrOurOwnURL   = errors.New("can't shorten our own url")
	ErrUserInfo    = errors.New("user info in url")
)

// Alias validation errors
var (
	ErrInvalidAlias  = errors.New("invalid alias")
	ErrAliasTooLong  = errors.New("alias is longer than available")
	ErrAliasTooShort = errors.New("alias is shorter than available")
)

// Alias regex
var aliasRe = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func (v Validator) URL(urlToCheck string) error {
	const op = "validation.Validator.URL"

	if urlToCheck == "" {
		return fmt.Errorf("%s: %w", op, ErrURLRequired)
	}

	if len(urlToCheck) > v.MaxURLLen {
		return fmt.Errorf("%s: %w", op, ErrURLToLong)
	}

	u, err := url.ParseRequestURI(urlToCheck)
	if err != nil {
		return fmt.Errorf("%s: %w", op, ErrInvalidURL)
	}

	if u.User != nil {
		return fmt.Errorf("%s: %w", op, ErrUserInfo)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%s: %w", op, ErrInvalidURL)
	}

	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("%s: %w", op, ErrInvalidURL)
	}

	host = strings.TrimSuffix(host, ".")
	if strings.Contains(host, "..") {
		return fmt.Errorf("%s: %w", op, ErrInvalidURL)
	}

	if v.MyHost != "" {
		myHost := strings.TrimSuffix(strings.TrimSpace(v.MyHost), ".")
		if strings.EqualFold(host, myHost) {
			return fmt.Errorf("%s: %w", op, ErrOurOwnURL)
		}
	}

	return nil
}

func (v Validator) Alias(alias string) error {
	const op = "validation.Alias.URL"

	if len(alias) < v.MinAliasLen {
		return fmt.Errorf("%s: %w", op, ErrAliasTooShort)
	}

	if len(alias) > v.MaxAliasLen {
		return fmt.Errorf("%s: %w", op, ErrAliasTooLong)
	}

	if match := aliasRe.MatchString(alias); !match {
		return fmt.Errorf("%s: %w", op, ErrInvalidAlias)
	}

	return nil
}
