CREATE TABLE "User" (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    login VARCHAR(50) UNIQUE
);

CREATE TABLE Account (
    account_id SERIAL PRIMARY KEY,
    balance MONEY DEFAULT '$0.00',
    description TEXT,
    bank VARCHAR(255) DEFAULT ''
);

CREATE TABLE UserAccount (
    user_account_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES "User"(user_id),
    account_id INTEGER REFERENCES Account(account_id)
);

CREATE TABLE Category (
    category_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES "User"(user_id),
    name VARCHAR(255)
);

CREATE TABLE Transaction (
    transaction_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES "User"(user_id),
    category_id INTEGER REFERENCES Category(category_id),
    account_id INTEGER REFERENCES Account(account_id),
    is_income BOOLEAN,
    total MONEY DEFAULT '$0.00',
    date DATE,
    payer VARCHAR(255) DEFAULT '',
    description TEXT
);

CREATE TABLE Deposit (
    deposit_id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES Account(account_id),
    total MONEY DEFAULT '$0.00',
    date_start DATE,
    deposit_term INTEGER,
    interest_rate DECIMAL(5, 2) DEFAULT 0.00,
    interest_calculation VARCHAR(255) DEFAULT ''
);

CREATE TABLE Credit (
    credit_id SERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES Account(account_id),
    total MONEY DEFAULT '$0.00',
    date_start DATE,
    summary TEXT DEFAULT '',
    date_end DATE,
    credit_calculation VARCHAR(255) DEFAULT '',
    payments MONEY DEFAULT '$0.00'
);

CREATE TABLE Investment (
    investment_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES "User"(user_id),
    name VARCHAR(255),
    total MONEY DEFAULT '$0.00',
    date_start DATE,
    price MONEY DEFAULT '$0.00',
    percentage DECIMAL(5, 2) DEFAULT 0.00
);

CREATE TABLE Debt (
    debt_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES "User"(user_id),
    total MONEY DEFAULT '$0.00',
    date DATE,
    status VARCHAR(50) DEFAULT '',
    description TEXT DEFAULT '',
    creditor VARCHAR(255) DEFAULT ''
);

CREATE TABLE Goal (
    goal_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES "User"(user_id),
    name VARCHAR(255),
    description TEXT DEFAULT '',
    total MONEY DEFAULT '$0.00',
    date DATE
);
