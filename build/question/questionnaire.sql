CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user (
    user_id         UUID PRIMARY KEY,
    question_id UUID REFERENCES question(id),
    rating INT NOT NULL,
);

CREATE TABLE IF NOT EXISTS question (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name TEXT CHECK (LENGTH(name) <= 100) UNIQUE NOT NULL,
    -- rating INT NOT NULL,
);



INSERT INTO question(name)
VALUES ('CSAT');

INSERT INTO question(name)
VALUES ('NPS1');

INSERT INTO question(name)
VALUES ('NPS2');
