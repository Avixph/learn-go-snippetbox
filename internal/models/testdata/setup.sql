-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- SELECT
--   uuid_generate_v4();
CREATE TABLE snippets (
  id uuid DEFAULT uuid_generate_v4() NOT NULL,
  title VARCHAR(120) NOT NULL,
  content TEXT NOT NULL,
  created_on TIMESTAMP NOT NULL,
  expires_on TIMESTAMP NOT NULL,
  PRIMARY KEY (id)
);

CREATE INDEX idx_snippets_created_on ON snippets(created_on);

CREATE TABLE users (
  id uuid DEFAULT uuid_generate_v4() NOT NULL,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  hashed_password CHAR(60) NOT NULL,
  created_on TIMESTAMP NOT NULL,
  PRIMARY KEY (id)
);

ALTER TABLE
  users
ADD
  CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO
  users (id, name, email, hashed_password, created_on)
VALUES
  (
    '6ba7b811-9dad-11d1-80b4-00c04fd430c8',
    'Nom Falso',
    'falso@example.com',
    '$2a$12$D2ndhbqWL99PVZPZDNX5nuWLqVU3pMvdyuBaJxhTnn5UlFw6Bu4Bq',
    '2023-01-23 13:25:37.403671'
  );