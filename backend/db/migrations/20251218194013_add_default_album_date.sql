-- migrate:up
ALTER TABLE albums
ALTER COLUMN album_date SET DEFAULT CURRENT_DATE;

-- migrate:down
ALTER TABLE albums
ALTER COLUMN album_date DROP DEFAULT;