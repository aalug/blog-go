CREATE TABLE "users"
(
    "id"                  BIGSERIAL PRIMARY KEY,
    "username"            VARCHAR        NOT NULL,
    "email"               VARCHAR UNIQUE NOT NULL,
    "hashed_password"     VARCHAR        NOT NULL,
    "password_changed_at" TIMESTAMPTZ    NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at"          TIMESTAMPTZ    NOT NULL DEFAULT (now())
);

CREATE INDEX idx_users_email ON users ("email");

CREATE TABLE "categories"
(
    "id"         BIGSERIAL PRIMARY KEY,
    "name"       VARCHAR     NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE INDEX idx_categories_name ON categories ("name");

CREATE TABLE "posts"
(
    "id"          BIGSERIAL PRIMARY KEY,
    "title"       VARCHAR     NOT NULL,
    "description" VARCHAR     NOT NULL,
    "content"     TEXT        NOT NULL,
    "author_id"   INTEGER REFERENCES users ("id"),
    "category_id" INTEGER REFERENCES categories ("id"),
    "image"       VARCHAR     NOT NULL,
    "created_at"  TIMESTAMPTZ NOT NULL DEFAULT (now()),
    "updated_at"  TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE INDEX idx_posts_title ON posts ("title");
CREATE INDEX idx_posts_created_at ON posts ("created_at");

CREATE TABLE "tags"
(
    "id"   SERIAL PRIMARY KEY,
    "name" VARCHAR(50) NOT NULL
);

CREATE INDEX idx_tags_name ON tags ("name");

CREATE TABLE "post_tags"
(
    "post_id" BIGINT REFERENCES posts ("id"),
    "tag_id"  INTEGER REFERENCES tags ("id"),
    PRIMARY KEY ("post_id", "tag_id")
);

CREATE TABLE "comments"
(
    "id"         BIGSERIAL PRIMARY KEY,
    "content"    TEXT        NOT NULL,
    "user_id"    INTEGER REFERENCES users ("id"),
    "post_id"    INTEGER REFERENCES posts ("id"),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);

CREATE INDEX idx_comments_created_at ON comments ("created_at");