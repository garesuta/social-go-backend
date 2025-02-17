CREATE TABLE roles (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  level INT NOT NULL
);

INSERT INTO roles (name, description, level) VALUES
  ('user', 'A user can create posts and comments', 1),
  ('moderator', 'A moderator can update other users posts', 2),
  ('admin', 'An admin can update and delete other users posts', 3);