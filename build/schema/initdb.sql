CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS Users (
    id             UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username       VARCHAR(20)              NOT NULL,
    login          VARCHAR(20)       UNIQUE NOT NULL,
    password_hash  VARCHAR(256)             NOT NULL,
	planned_budget numeric(10, 2),
    avatar_url     UUID
);

CREATE TABLE IF NOT EXISTS Accounts (
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    balance numeric(10, 2),
    accumulation BOOLEAN,
    balance_enabled BOOLEAN,
    mean_payment VARCHAR(30)
);

CREATE TABLE IF NOT EXISTS UserAccount (
    user_id    UUID REFERENCES Users(id),
    account_id UUID REFERENCES Accounts(id),
    PRIMARY KEY (user_id, account_id)
);

CREATE TABLE IF NOT EXISTS category (
    id              UUID          DEFAULT uuid_generate_v4()   PRIMARY KEY,
    user_id         UUID          REFERENCES Users(id)    CONSTRAINT fk_user_category       NOT NULL,
    parent_tag      UUID          REFERENCES category(id),
    "name"          VARCHAR(30)                                                             NOT NULL,
    --image_id        INT           DEFAULT 0                                                 NOT NULL,
    show_income     BOOLEAN,
    show_outcome    BOOLEAN,
    regular         BOOLEAN                                                                 NOT NULL
);


CREATE TABLE IF NOT EXISTS Transaction (
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

CREATE TABLE IF NOT EXISTS TransactionCategory (
    transaction_id UUID REFERENCES Transaction(id),
    category_id UUID REFERENCES Category(id),
    PRIMARY KEY (transaction_id, category_id)
);

-- CREATE TABLE IF NOT EXISTS goal (
--     id            UUID            DEFAULT uuid_generate_v4() PRIMARY KEY,
--     user_id       UUID            REFERENCES Users(id)                                             NOT NULL,
--     "name"        TEXT                                       CHECK(LENGTH("name") <= 50)           NOT NULL,
--     "description" TEXT            DEFAULT ''                 CHECK(LENGTH("description") <= 255),
--     "target"      NUMERIC(10,2)                                                                    NOT NULL,
--     "date"        DATE,
--     "state"       TEXT            DEFAULT ''                 CHECK(LENGTH("state") <= 20),
--     created_at    TIMESTAMPTZ     DEFAULT CURRENT_TIMESTAMP                                        NOT NULL,
--     updated_at    TIMESTAMPTZ     DEFAULT CURRENT_TIMESTAMP                                        NOT NULL
-- );

--========================================================================

--CREATE OR REPLACE FUNCTION public.moddatetime()
--    RETURNS TRIGGER AS $$
--BEGIN
--    NEW.updated_at = NOW();
--    RETURN NEW;
--END;
--$$ LANGUAGE plpgsql;
--
--CREATE OR REPLACE TRIGGER modify_updated_at
--    BEFORE UPDATE
--    ON goal
--    FOR EACH ROW
--EXECUTE PROCEDURE public.moddatetime(updated_at);

--========================================================================

CREATE OR REPLACE FUNCTION add_default_categories_accounts_transactions()
RETURNS TRIGGER AS $$
DECLARE
    categoryID UUID;
    transaction_idI UUID;
    transaction_idO UUID;
    accountCashID UUID;
    accountCardID UUID;
BEGIN
    INSERT INTO category (user_id, parent_tag, "name", show_income, show_outcome, regular)
    VALUES  (NEW.id, NULL, 'Дети',                     false, true,  false),
            (NEW.id, NULL, 'Забота о себе',            false, true,  false),
            (NEW.id, NULL, 'Зарплата',                 true,  false, true),
            (NEW.id, NULL, 'Здоровье и фитнес',         false, true,  false),
            (NEW.id, NULL, 'Кафе и рестораны',         false, true,  false),
            (NEW.id, NULL, 'Машина',                   false, true,  false),
            (NEW.id, NULL, 'Образование',              false, true,  false),
            (NEW.id, NULL, 'Отдых и развлечения',      false, true,  false),
            (NEW.id, NULL, 'Подарки',                  false, true,  false),
            (NEW.id, NULL, 'Покупки: одежда, техника', false, true,  false),
            (NEW.id, NULL, 'Проезд',                   false, true,  false),
            (NEW.id, NULL, 'Подписки',                 false, true,  true),
            (NEW.id, NULL, 'Продукты',                 false, true, false);
    
    SELECT id INTO categoryID FROM category WHERE name = 'Продукты' AND user_id = NEW.id;

    INSERT INTO accounts(balance, mean_payment, accumulation, balance_enabled)
    VALUES (0, 'Карта', false, true) RETURNING id INTO accountCardID;
           
    INSERT INTO accounts(balance, mean_payment, accumulation, balance_enabled)
    VALUES (0, 'Наличка', false, true) RETURNING id INTO accountCashID;

    INSERT INTO userAccount(user_id, account_id)
    VALUES (NEW.id, accountCardID);

    INSERT INTO userAccount(user_id, account_id)
    VALUES (NEW.id, accountCashID);

    INSERT INTO transaction(user_id, account_income, account_outcome, income, outcome, payer, description)
    VALUES (NEW.id, accountCardID,
                    accountCardID, 0, 100, 'Тестовая транзакция1', 'Все хомячьки приветствуют вас, и просят удалить эти транзации)') RETURNING id INTO transaction_idI;

    INSERT INTO transaction(user_id, account_income, account_outcome, income, outcome, payer, description)
    VALUES (NEW.id, accountCardID,
                    accountCardID, 100, 0, 'Тестовая транзакция2', 'Все хомячьки приветствуют вас, и просят удалить эти транзации)') RETURNING id INTO transaction_idO;
            
    INSERT INTO TransactionCategory(transaction_id, category_id)
    VALUES (transaction_idI, categoryID),
            (transaction_idO, categoryID);

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_user_created
    AFTER INSERT ON users
    FOR EACH ROW
    EXECUTE FUNCTION add_default_categories_accounts_transactions();
--=============================================================================

ALTER TABLE Users
ALTER COLUMN planned_budget SET DEFAULT 0.0;

--=============================================================================

INSERT INTO "users"(login, username, password_hash, planned_budget)
VALUES ('kossmatof','komarov', '$argon2id$v=19$m=65536,t=1,p=4$m8qhM3XLae+RCTGirBFEww$Znu5RBnxlam2xRoVtwBzbdSrN4/sRCm1IMOVX4N2uxw', 10000);
