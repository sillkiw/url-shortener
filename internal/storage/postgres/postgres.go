package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/sillkiw/url-shorten/internal/storage"
)

const uniqueViolationCode = "23505"

type Storage struct {
	db *sql.DB
}

func New(postgresUrl string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"
	const q = `
		INSERT INTO short_links(original_url, alias) 
		VALUES ($1, $2)  
		RETURNING id
	`

	var id int64
	err := s.db.QueryRow(q, urlToSave, alias).Scan(&id)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == uniqueViolationCode {
			switch pqErr.Constraint {
			case "short_links_original_url_key":
				return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
			case "short_links_alias_key":
				return 0, fmt.Errorf("%s: %w", op, storage.ErrAliasExists)
			}
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (url string, err error) {
	const op = "storage.postgres.GetURL"
	const q = `
		SELECT original_url 
		FROM short_links
		WHERE alias = $1
	`
	var originalURL string
	err = s.db.QueryRow(q, alias).Scan(&originalURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrAliasNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return originalURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"
	const q = `
		DELETE FROM short_links
		WHERE alias = $1
	`
	res, err := s.db.Exec(q, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: rows affected: %w", op, err)
	}
	if n == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrAliasNotFound)
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
