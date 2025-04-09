CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS posts (
    "id" UUID PRIMARY KEY,
    "user_id" UUID NOT NULL,
    "title" VARCHAR(20) NOT NULL,
    "content" TEXT NOT NULL,
    "tags" VARCHAR(20)[] NOT NULL,
    "created_at" TIMESTAMPTZ,
    "updated_at" TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);