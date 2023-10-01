CREATE TABLE Users
(
    id            SERIAL       PRIMARY KEY,
    username      VARCHAR(20)  UNIQUE      NOT NULL,
    password_hash VARCHAR(256)             NOT NULL,
    first_name    VARCHAR(20),
    last_name     VARCHAR(20),
	planned_budget MONEY DEFAULT 0,
    avatar_url    TEXT
);

CREATE TABLE Accounts (
    id UUID PRIMARY KEY,
    user_id INT,
    balance MONEY,
    mean_payment TEXT
);

--=============================================================================

INSERT INTO "users"(username, password_hash, first_name, last_name, planned_budget, avatar_url)
VALUES ('kosmatoff', 'hash','Дмитрий', 'Комаров', 1000000, 'image/img1.png');

INSERT INTO "accounts"(UserID, Balance, MeanPayment)
VALUES (1, 'Карта', 25000);

INSERT INTO "accounts"(UserID, Balance, MeanPayment)
VALUES (1, 'Наличные', 450);
