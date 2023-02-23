CREATE TABLE IF NOT EXISTS mangas (
	id bigserial PRIMARY KEY,
	title text NOT NULL,
	studio text NOT NULL,
	year integer NOT NULL,
	chapters integer NOT NULL,
	rating NUMERIC(5, 2) NOT NULL,
	version integer NOT NULL DEFAULT 1
);