-- migrate:up

-- SEQUENCES
CREATE SEQUENCE public.albums_id_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE public.images_id_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE public.users_id_seq  START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE public.walls_id_seq  START WITH 1 INCREMENT BY 1;

-- TABLES
CREATE TABLE public.walls (
  id          integer PRIMARY KEY DEFAULT nextval('public.walls_id_seq'),
  name        text NOT NULL,
  created_at  timestamp without time zone DEFAULT now()
);

CREATE TABLE public.users (
  id            integer PRIMARY KEY DEFAULT nextval('public.users_id_seq'),
  username      text NOT NULL UNIQUE,
  email         text NOT NULL UNIQUE,
  password_hash text NOT NULL,
  created_at    timestamp without time zone DEFAULT now(),
  wall_id       integer REFERENCES public.walls(id)
);

CREATE TABLE public.albums (
  id         integer PRIMARY KEY DEFAULT nextval('public.albums_id_seq'),
  album_date date NOT NULL,
  wall_id    integer REFERENCES public.walls(id),
  memo       text,
  CONSTRAINT albums_album_date_wall_id_key UNIQUE (album_date, wall_id)
);

CREATE TABLE public.images (
  id           integer PRIMARY KEY DEFAULT nextval('public.images_id_seq'),
  filename     text NOT NULL,                 -- store R2 object key here
  upload_time  timestamp without time zone DEFAULT now(),
  memo         text,
  user_id      integer REFERENCES public.users(id),
  wall_id      integer REFERENCES public.walls(id),
  album_date   date DEFAULT CURRENT_DATE
);

-- (Optional) helpful indexes
CREATE INDEX IF NOT EXISTS idx_images_wall_date ON public.images (wall_id, album_date);
CREATE INDEX IF NOT EXISTS idx_images_wall_time ON public.images (wall_id, upload_time DESC);

-- migrate:down

DROP INDEX IF EXISTS idx_images_wall_time;
DROP INDEX IF EXISTS idx_images_wall_date;

DROP TABLE IF EXISTS public.images;
DROP TABLE IF EXISTS public.albums;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.walls;

DROP SEQUENCE IF EXISTS public.walls_id_seq;
DROP SEQUENCE IF EXISTS public.users_id_seq;
DROP SEQUENCE IF EXISTS public.images_id_seq;
DROP SEQUENCE IF EXISTS public.albums_id_seq;
