CREATE TABLE IF NOT EXISTS short_links (
    id BIGSERIAL PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    original_url TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_alias ON short_links(alias)
