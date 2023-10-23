CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "user" (
    user_id    UUID   DEFAULT uuid_generate_v4() PRIMARY KEY,
    full_name  VARCHAR(255) NOT NULL,
    username   VARCHAR(50)  UNIQUE NOT NULL
);


CREATE TABLE IF NOT EXISTS account (
    account_id UUID      DEFAULT uuid_generate_v4() PRIMARY KEY,
    account_balance    MONEY DEFAULT '$0.00',
    account_description TEXT,
    bank_name VARCHAR(255) DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS user_account (
    user_account_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES "user"(user_id),
    account_id UUID REFERENCES account(account_id)
);


CREATE TABLE IF NOT EXISTS category (
    category_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES user(user_id),
    category_name VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS "transaction" (
    transaction_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES "user"(user_id),
    category_id UUID REFERENCES "category"(category_id),
    account_id UUID REFERENCES "account"(account_id),
    is_income BOOLEAN,
    total_money MONEY DEFAULT '$0.00',
    transaction_date DATE,
    payer_name VARCHAR(40) DEFAULT '',
    transaction_description TEXT,
    created_at_timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at_timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS deposit (
    deposit_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    account_id UUID REFERENCES account(account_id),
    total_money MONEY DEFAULT '$0.00',
    start_date DATE DEFAULT CURRENT_DATE,
    end_date DATE,
    interest_rate DECIMAL(5, 2) DEFAULT 0.00,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS credit (
    credit_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    account_id UUID REFERENCES account(account_id) NOT NULL,
    total_amount MONEY DEFAULT '$0.00' NOT NULL,
    monthly_payment MONEY DEFAULT '$0.00',
    start_date DATE  DEFAULT CURRENT_DATE,
    end_date DATE,
    calculation_details VARCHAR(30) DEFAULT '',
    payments_received MONEY DEFAULT '$0.00',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS investment (
    investment_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES "user"(user_id) NOT NULL,
    asset_type VARCHAR(255) NOT NULL,
    asset_name VARCHAR(255) NOT NULL,
    purchase_price MONEY DEFAULT '$0.00' NOT NULL,
    quantity NUMERIC NOT NULL,
    purchase_date DATE DEFAULT CURRENT_DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS goal (
    goal_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id UUID REFERENCES "user"(user_id) NOT NULL,
    goal_name VARCHAR(255),
    goal_description TEXT DEFAULT '',
    goal_target MONEY DEFAULT '$0.00',
    goal_date DATE
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
