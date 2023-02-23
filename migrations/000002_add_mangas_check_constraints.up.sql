ALTER TABLE mangas ADD CONSTRAINT mangas_rating_check CHECK (rating >= 0);
ALTER TABLE mangas ADD CONSTRAINT mangas_chapters_check CHECK (chapters >= 0);
ALTER TABLE mangas ADD CONSTRAINT mangas_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));