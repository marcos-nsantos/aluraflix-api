ALTER TABLE videos
    ADD COLUMN category_id INTEGER REFERENCES categories (id) ON DELETE CASCADE;