CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE walls (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    filename TEXT NOT NULL,
    upload_time TIMESTAMP DEFAULT now(),
    album_date DATE DEFAULT CURRENT_DATE,
    memo TEXT,
    user_id INTEGER REFERENCES users(id),
    wall_id INTEGER REFERENCES walls(id)
);

CREATE TABLE albums (
    id SERIAL PRIMARY KEY,
    album_date DATE NOT NULL,
    wall_id INTEGER REFERENCES walls(id),
    memo TEXT,
    UNIQUE (album_date, wall_id)
);


