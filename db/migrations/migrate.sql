CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "user" (
    user_id        UUID         DEFAULT uuid_generate_v4() PRIMARY KEY,
    login          VARCHAR(255)                                   NOT NULL,
    username       VARCHAR(50)                             UNIQUE NOT NULL,
    password_hash  VARCHAR(256)                                   NOT NULL,
	planned_budget MONEY        DEFAULT 0.00,
    avatar_url     UUID
);


CREATE TABLE IF NOT EXISTS account (
    account_id          UUID        DEFAULT uuid_generate_v4() PRIMARY KEY,
    account_balance     MONEY       DEFAULT 0.00,
    account_description TEXT,
    bank_name           VARCHAR(30),
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS user_account (
    user_account_id UUID DEFAULT uuid_generate_v4()     PRIMARY KEY,
    user_id         UUID REFERENCES "user"(user_id)     CONSTRAINT fk_user      NOT NULL,
    account_id      UUID REFERENCES account(account_id) CONSTRAINT fk_account   NOT NULL
);


CREATE TABLE IF NOT EXISTS category (
    category_id   UUID          DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id       UUID          REFERENCES "user"(user_id)   CONSTRAINT fk_user_category  NOT NULL,
    "name"        VARCHAR(20)                                                             NOT NULL
);

CREATE TABLE IF NOT EXISTS "transaction" (
    transaction_id          UUID DEFAULT uuid_generate_v4()         PRIMARY KEY,
    user_id                 UUID REFERENCES "user"(user_id)         CONSTRAINT fk_user_transaction     NOT NULL,
    category_id             UUID REFERENCES "category"(category_id) CONSTRAINT fk_category_transaction NOT NULL,
    account_id              UUID REFERENCES "account"(account_id)   CONSTRAINT fk_account_transaction  NOT NULL,
    is_income               BOOLEAN                                                                    NOT NULL,
    total_money             MONEY       DEFAULT 0.00                                                   NOT NULL,
    "date"                  DATE        DEFAULT CURRENT_DATE                                           NOT NULL,
    receiver_name           VARCHAR(40) DEFAULT '',
    "description"           TEXT,
    created_at              TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP                                      NOT NULL,
    updated_at              TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP                                      NOT NULL
);


CREATE TABLE IF NOT EXISTS deposit (
    deposit_id      UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    account_id      UUID REFERENCES account(account_id) CONSTRAINT fk_account_deposit NOT NULL,
    total_money     MONEY         DEFAULT 0.00                                        NOT NULL,
    start_date      DATE          DEFAULT CURRENT_DATE                                NOT NULL,
    end_date        DATE                                                              NOT NULL,
    interest_rate   DECIMAL(5, 2) DEFAULT 0.00                                        NOT NULL,
    created_at      TIMESTAMPTZ   DEFAULT CURRENT_TIMESTAMP                           NOT NULL,
    updated_at      TIMESTAMPTZ   DEFAULT CURRENT_TIMESTAMP                           NOT NULL
);

CREATE TABLE IF NOT EXISTS credit (
    credit_id           UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    account_id          UUID REFERENCES account(account_id) CONSTRAINT fk_account_credit NOT NULL,
    total_amount        MONEY       DEFAULT 0.00                                         NOT NULL,
    monthly_payment     MONEY       DEFAULT 0.00                                         NOT NULL,
    start_date          DATE        DEFAULT CURRENT_DATE                                 NOT NULL,
    end_date            DATE,
    calculation_details VARCHAR(30) DEFAULT 'annuity'                                    NOT NULL,
    payments_received   MONEY       DEFAULT 0.00,
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP                            NOT NULL,
    updated_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP                            NOT NULL
);


CREATE TABLE IF NOT EXISTS investment (
    investment_id   UUID            DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id         UUID            REFERENCES "user"(user_id) CONSTRAINT fk_user_investment NOT NULL,
    asset_type      VARCHAR(255)                                                             NOT NULL,
    asset_name      VARCHAR(255)                                                             NOT NULL,
    purchase_price  MONEY           DEFAULT 0.00                                             NOT NULL,
    quantity        NUMERIC                                                                  NOT NULL,
    purchase_date   DATE            DEFAULT CURRENT_DATE                                     NOT NULL,
    created_at      TIMESTAMPTZ     DEFAULT CURRENT_TIMESTAMP                                NOT NULL,
    updated_at      TIMESTAMPTZ     DEFAULT CURRENT_TIMESTAMP                                NOT NULL
);


CREATE TABLE IF NOT EXISTS goal (
    goal_id          UUID           DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id          UUID           REFERENCES "user"(user_id) CONSTRAINT fk_user_goal NOT NULL,
    "name"           VARCHAR(255)                                                      NOT NULL,
    "description"    TEXT           DEFAULT '',
    amount           MONEY          DEFAULT 0.00                                       NOT NULL,
    start_date        DATE                                                              NOT NULL,
    created_at       TIMESTAMPTZ    DEFAULT CURRENT_TIMESTAMP                          NOT NULL,
    updated_at       TIMESTAMPTZ    DEFAULT CURRENT_TIMESTAMP                          NOT NULL
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