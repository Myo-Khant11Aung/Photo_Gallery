-- migrate:up
ALTER TABLE users ADD COLUMN wall_id INTEGER REFERENCES walls(id);

-- migrate:down

