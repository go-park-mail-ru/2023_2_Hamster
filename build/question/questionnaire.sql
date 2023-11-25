CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS User (
    user_id         UUID PRIMARY KEY,
    answer_id UUID REFERENCES Answer(id),
);

CREATE TABLE IF NOT EXISTS Answer (
    id UUID PRIMARY KEY,
    name TEXT CHECK (LENGTH(name) <= 100)
    rating INT NOT NULL,
);