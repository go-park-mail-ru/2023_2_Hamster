CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE Users
(
    id             UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username       VARCHAR(20)              NOT NULL,
    login          VARCHAR(20)       UNIQUE NOT NULL,
    password_hash  VARCHAR(256)             NOT NULL,
	planned_budget numeric(10, 2),
    avatar_url     UUID
);

CREATE TABLE Accounts (
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id       UUID REFERENCES Users(id),
    balance numeric(10, 2),
    mean_payment TEXT
);

CREATE TABLE IF NOT EXISTS category (
    id              UUID          DEFAULT uuid_generate_v4()   PRIMARY KEY,
    user_id         UUID          REFERENCES Users(id)    CONSTRAINT fk_user_category  NOT NULL,
    parent_tag      UUID          REFERENCES category(id),
    "name"          VARCHAR(30)                                                             NOT NULL,
    show_income     BOOLEAN,
    show_outcome    BOOLEAN,
    regular         BOOLEAN                                                                 NOT NULL
);


CREATE TABLE Transaction (
	id           UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	user_id      UUID REFERENCES Users(id),
    account_income UUID REFERENCES Accounts(id),
    account_outcome UUID REFERENCES Accounts(id),
	income       numeric(10, 2),
    outcome      numeric(10, 2),
	date         timestamp DEFAULT now(),
	payer        VARCHAR(20),
	description  VARCHAR(100)
);

CREATE TABLE TransactionCategory (
    transaction_id UUID REFERENCES Transaction(id),
    category_id UUID REFERENCES Category(id),
    PRIMARY KEY (transaction_id, category_id)
);

CREATE OR REPLACE FUNCTION add_default_categoies()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO category (user_id, parent_tag, "name", show_income, show_outcome, regular)
    VALUES  (NEW.id, NULL, 'Дети',                    false, true,  false),
            (NEW.id, NULL, 'Забота о себе',           false, true,  false),
            (NEW.id, NULL, 'Зарплата',                true,  false, true),
            (NEW.id, NULL, 'Здровье и фитнес',        false, true,  false),
            (NEW.id, NULL, 'Кафе и рестораны',        false, true,  false),
            (NEW.id, NULL, 'Машина',                  false, true,  false),
            (NEW.id, NULL, 'Образование',             false, true,  false),
            (NEW.id, NULL, 'Отдых и развлечения',     false, true,  false),
            (NEW.id, NULL, 'Подарки',                 false, true,  false),
            (NEW.id, NULL, 'Покупки: одежа, техника', false, true,  false),
            (NEW.id, NULL, 'Проезд',                  false, true,  false),
            (NEW.id, NULL, 'Продукты',                false, true,  false),
            (NEW.id, NULL, 'Подписки',                false, true,  true);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER after_user_created
    AFTER INSERT ON users
    FOR EACH ROW
EXECUTE FUNCTION add_default_categoies();
--=============================================================================

ALTER TABLE Users
ALTER COLUMN planned_budget SET DEFAULT 0.0;

--=============================================================================

INSERT INTO "users"(login, username, password_hash, planned_budget)
VALUES ('kossmatof','komarov', '$argon2id$v=19$m=65536,t=1,p=4$m8qhM3XLae+RCTGirBFEww$Znu5RBnxlam2xRoVtwBzbdSrN4/sRCm1IMOVX4N2uxw', 10000);

INSERT INTO "users"(login, username, password_hash, planned_budget)
VALUES ('test','test1', '$argon2id$v=19$m=65536,t=1,p=4$m8qhM3XLae+RCTGirBFEww$Znu5RBnxlam2xRoVtwBzbdSrN4/sRCm1IMOVX4N2uxw', 10000);

INSERT INTO "accounts"(user_id, balance, mean_payment)
VALUES ((SELECT id FROM Users limit 1), 1000, 'Кошелек');

INSERT INTO "accounts"(user_id, balance, mean_payment)
VALUES ((SELECT id FROM Users limit 1), 1000, 'Наличка');
