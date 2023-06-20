ALTER TABLE categories
    ADD CONSTRAINT unique_name UNIQUE (name);