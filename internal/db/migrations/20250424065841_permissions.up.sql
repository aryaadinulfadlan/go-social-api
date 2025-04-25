CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS permissions (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(50) NOT NULL,
    "description" VARCHAR(100) NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO permissions (id, name, description)
VALUES
  ('50d71cf1-19c3-49d3-b670-e5e0c31b079f', 'post:create', 'Create a new post.'),
  ('8dd90a7a-a1b9-4ae7-be90-4051071daa1b', 'post:detail', 'Get detail of a post.'),
  ('b277ca4b-8be4-4119-95e3-26b2089021c5', 'post:update', 'Update or modify an existing post.'),
  ('ea2ab559-3a6b-4320-8d32-41c11b092a44', 'post:delete', 'Remove or delete a post.'),
  ('3f1c5226-7e68-4034-98d0-8dfa19c8d4e6', 'user:detail', 'Get detail of a user account.'),
  ('463d9784-e62b-4d51-a98e-d5ecd04d0947', 'user:follow', 'Follow or unfollow a user account.'),
  ('a3e14d41-62b5-48d8-9a1d-c95117a78b25', 'user:followers', 'Get followers of a user account.'),
  ('83e2c6f1-f425-45cc-b360-a7f6eb9f57e7', 'user:following', 'Get following of a user account.'),
  ('59dc7703-b90b-472a-8722-5785406e836b', 'user:feed', 'Get the feed of a user account.'),
  ('fd408311-636b-44fe-8f93-d847acbc3ccf', 'user:delete', 'Remove or delete a user account.'),
  ('86329004-3735-41ae-93a3-8c0fe965845f', 'comment:create', 'Create a new comment for a post.');