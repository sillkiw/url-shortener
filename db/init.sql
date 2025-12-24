CREATE TABLE IF NOT EXISTS short_links (
    id BIGSERIAL PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    original_url TEXT NOT NULL UNIQUE
);

CREATE UNIQUE INDEX short_links_slug_uq ON short_links(slug)

