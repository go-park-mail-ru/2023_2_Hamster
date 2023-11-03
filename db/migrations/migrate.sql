CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "user" (
    user_id        UUID         DEFAULT uuid_generate_v4() PRIMARY KEY,
    username       TEXT         UNIQUE                     NOT NULL,
    full_name      TEXT                                    NOT NULL,
    password_hash  TEXT                                    NOT NULL,
	planned_budget MONEY        DEFAULT 0.00               NOT NULL,
    avatar_url     UUID,
    CHECK(LENGTH(username) <= 30),
    CHECK(LENGTH(full_name) <= 60)
);


CREATE TABLE IF NOT EXISTS account (
    account_id          UUID         DEFAULT uuid_generate_v4() PRIMARY KEY,
    balance             MONEY        DEFAULT 0.00,
    "description"       TEXT,
    bank_name           TEXT,
    currency            TEXT         DEFAULT 'RUB',
    created_at          TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at          TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CHECK(LENGTH("description") <= 255),
    CHECK(LENGTH(bank_name) <= 60),
    CHECK(LENGTH(currency) <= 3)
);


CREATE TABLE IF NOT EXISTS user_account (
    user_account_id UUID DEFAULT    uuid_generate_v4()     PRIMARY KEY,
    user_id         UUID REFERENCES "user"(user_id)     CONSTRAINT fk_user      NOT NULL,
    account_id      UUID REFERENCES account(account_id) CONSTRAINT fk_account   NOT NULL
);


CREATE TABLE IF NOT EXISTS category (
    category_id   UUID          DEFAULT uuid_generate_v4()   PRIMARY KEY,
    user_id       UUID          REFERENCES "user"(user_id)   CONSTRAINT fk_user_category  NOT NULL,
    "name"        TEXT                                                                    NOT NULL,
    CHECK(LENGTH("name") <= 60)
);

CREATE TABLE IF NOT EXISTS "transaction" (
    transaction_id          UUID        DEFAULT    uuid_generate_v4()       PRIMARY KEY,
    user_id                 UUID        REFERENCES "user"(user_id),
    category_id             UUID        REFERENCES "category"(category_id), 
    account_id              UUID        REFERENCES "account"(account_id),
    is_income               BOOLEAN,
    total                   MONEY       DEFAULT 0.00,
    "date"                  DATE        DEFAULT CURRENT_DATE                NOT NULL,
    payer_name              TEXT        DEFAULT ''                          NOT NULL,
    "description"           TEXT,
    created_at              TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP           NOT NULL,
    updated_at              TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP           NOT NULL,
    CHECK(LENGTH(payer_name) <= 50),
    CHECK(LENGTH("description") <= 255)
);


CREATE TABLE IF NOT EXISTS deposit (
    deposit_id    UUID          DEFAULT uuid_generate_v4()      PRIMARY KEY,
    account_id    UUID          REFERENCES account(account_id),
    total_amount  MONEY         DEFAULT 0.00,
    interest_rate DECIMAL(5, 2) DEFAULT 0.00,
    start_at      DATE          DEFAULT CURRENT_DATE            NOT NULL,
    end_at        DATE,
    created_at    TIMESTAMPTZ   DEFAULT CURRENT_TIMESTAMP       NOT NULL,
    updated_at    TIMESTAMPTZ   DEFAULT CURRENT_TIMESTAMP       NOT NULL
);

CREATE TABLE IF NOT EXISTS credit (
    credit_id           UUID        DEFAULT uuid_generate_v4()     PRIMARY KEY,
    account_id          UUID        REFERENCES account(account_id) NOT NULL,
    total_amount        MONEY       DEFAULT 0.00                   NOT NULL,
    "status"            TEXT        DEFAULT 'active'               NOT NULL,
    monthly_payment     MONEY       DEFAULT 0.00,
    start_at            DATE        DEFAULT CURRENT_DATE,
    calculation_details TEXT        DEFAULT '',
    end_at              DATE,
    payments_received   MONEY       DEFAULT 0.00                   NOT NULL,
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP      NOT NULL,
    updated_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP      NOT NULL,
    CHECK(LENGTH(calculation_details) <= 20),
    CHECK(LENGTH("status") <= 60)
);


CREATE TABLE IF NOT EXISTS investment (
    investment_id   UUID         DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id         UUID         REFERENCES "user"(user_id) NOT NULL,
    asset_type      TEXT                                    NOT NULL,
    asset_name      TEXT                                    NOT NULL,
    purchase_price  MONEY        DEFAULT 0.00               NOT NULL,
    quantity        NUMERIC                                 NOT NULL,
    purchase_at     DATE         DEFAULT CURRENT_DATE       NOT NULL,
    created_at      TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    updated_at      TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    CHECK(LENGTH(asset_name) <= 70),
    CHECK(LENGTH(asset_type) <= 40)
);


CREATE TABLE IF NOT EXISTS goal (
    goal_id       UUID            DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id       UUID            REFERENCES "user"(user_id) NOT NULL,
    "name"        TEXT,
    "description" TEXT            DEFAULT '',
    "target"      MONEY           DEFAULT 0.00               NOT NULL,
    "date"        DATE,
    created_at    TIMESTAMPTZ     DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    updated_at    TIMESTAMPTZ     DEFAULT CURRENT_TIMESTAMP  NOT NULL,
    CHECK(LENGTH("description") <= 255),
    CHECK(LENGTH("name") <= 50)
);

CREATE OR REPLACE FUNCTION public.moddatetime()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER modify_updated_at
    BEFORE UPDATE
    ON account
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

CREATE OR REPLACE TRIGGER modify_updated_at
    BEFORE UPDATE
    ON "transaction"
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

CREATE OR REPLACE TRIGGER modify_updated_at
    BEFORE UPDATE
    ON deposit
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

CREATE OR REPLACE TRIGGER modify_updated_at
    BEFORE UPDATE
    ON credit
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

CREATE OR REPLACE TRIGGER modify_updated_at
    BEFORE UPDATE
    ON investment
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);

CREATE OR REPLACE TRIGGER modify_updated_at
    BEFORE UPDATE
    ON goal
    FOR EACH ROW
EXECUTE PROCEDURE public.moddatetime(updated_at);