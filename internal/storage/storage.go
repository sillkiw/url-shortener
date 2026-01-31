package storage

import "errors"

var (
	ErrAliasNotFound = errors.New("alias not found")
	ErrURLExists     = errors.New("url exists")
	ErrAliasExists   = errors.New("alias exists")
)
