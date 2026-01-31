package validation_test

import (
	"errors"
	"testing"

	valid "github.com/sillkiw/url-shorten/internal/lib/validation"
)

func TestValidatorURL(t *testing.T) {
	v := valid.Validator{MyHost: "my.app", MaxURLLen: 4096, MinAliasLen: 3, MaxAliasLen: 10}
	type testCase struct {
		name     string
		url      string
		expected error
	}

	tests := []testCase{
		{name: "empty", url: "", expected: valid.ErrURLRequired},
		{name: "space", url: "      ", expected: valid.ErrURLRequired},
		{name: "non-http scheme", url: "postgres://localhost:5432/urldb", expected: valid.ErrInvalidURL},
		{name: "non-http scheme", url: "ftp://example.com/file", expected: valid.ErrInvalidURL},
		{name: "space in hostname", url: "http://ex mple.com", expected: valid.ErrInvalidURL},
		{name: "user info", url: "http://user:pass@example.com", expected: valid.ErrUserInfo},
		{name: "upper case", url: "HTTP://EXAMPLE.COM", expected: nil},
		{name: "two dots", url: "http://example..com", expected: valid.ErrInvalidURL},
		{name: "no host", url: "http://:80", expected: valid.ErrInvalidURL},
		{name: "ok", url: "https://github.com/defer-panic/dfrp.cc/blob/master/internal/shorten/shorten_test.go", expected: nil},
		{name: "myhost", url: "https://my.app/s/abs", expected: valid.ErrOurOwnURL},
		{name: "myhost with dot", url: "https://my.app./s/abc", expected: valid.ErrOurOwnURL},
		{name: "myhost with port", url: "https://my.app:8000/s/abs", expected: valid.ErrOurOwnURL},
		{name: "myhost with uppercase scheme", url: "HTTPS://my.app/s/abs", expected: valid.ErrOurOwnURL},
		{name: "myhost with uppercase", url: "HTTPS://my.aPp/s/abs", expected: valid.ErrOurOwnURL},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := v.URL(tc.url)

			switch {
			case tc.expected == nil:
				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}
			default:
				if err == nil {
					t.Fatalf("expected error %v, got nil", tc.expected)
				}
				if !errors.Is(err, tc.expected) {
					t.Fatalf("expected error %v, got %v", tc.expected, err)
				}
			}
		})
	}
}

func TestValidatorAlias(t *testing.T) {
	v := valid.Validator{MyHost: "my.app", MaxURLLen: 4096, MinAliasLen: 3, MaxAliasLen: 10}
	type testCase struct {
		name     string
		alias    string
		expected error
	}

	tests := []testCase{
		{name: "alias too long", alias: "asdlfshdfgashdgfadsgfjhadsgf", expected: valid.ErrAliasTooLong},
		{name: "alias too short", alias: "a", expected: valid.ErrAliasTooShort},
		{name: "alias is not match regex", alias: "+%$#!@#!@", expected: valid.ErrInvalidAlias},
		{name: "alias is ok", alias: "abc", expected: nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := v.Alias(tc.alias)

			switch {
			case tc.expected == nil:
				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}
			default:
				if err == nil {
					t.Fatalf("expected error %v, got nil", tc.expected)
				}
				if !errors.Is(err, tc.expected) {
					t.Fatalf("expected error %v, got %v", tc.expected, err)
				}
			}
		})
	}
}
