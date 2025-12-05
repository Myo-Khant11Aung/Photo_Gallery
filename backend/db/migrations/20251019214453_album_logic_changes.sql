-- migrate:up
ALTER TABLE public.albums
    ADD COLUMN name text NOT NULL DEFAULT 'Untitled Album',
    DROP CONSTRAINT albums_album_date_wall_id_key,
    ADD CONSTRAINT albums_name_wall_id_key UNIQUE (name,wall_id);
ALTER TABLE public.images 
    ADD COLUMN album_id INTEGER REFERENCES public.albums(id);
-- migrate:down
ALTER TABLE public.images
DROP COLUMN album_id;

ALTER TABLE public.albums
DROP COLUMN name,
DROP CONSTRAINT unique_album_name_per_wall,
ADD CONSTRAINT albums_album_date_wall_id_key UNIQUE (album_date, wall_id);
