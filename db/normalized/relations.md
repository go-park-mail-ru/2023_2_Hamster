# Описание Таблиц
## User
- Хранит информацию о пользователях.
- {id} -> {username, login, password_hash, planned_budget, avatar_url}
- {login}->{id, username, password_hash, planned_budget, avatar_url}
## Account
- Хранит информацию о банковских счетах пользователя.
- {id} -> {user_id, balance, mean_payment}
- {user_id} -> {id, balance, mean_payment}
## Investment
- Хранит информацию о инвестициях пользователя.
- id -> {name, total, date start, price, percentage}
## Category
- Хранит информацию о категориях транзакций.
- {id} -> {user_id, name}
## Transaction
- Хранит информацию о транзакциях пользователя.
- {id} -> {user_id, category_id, account_id, total, is_income, date, payer, description}
## Goal
- Хранит информацию о финансовых целях пользователя.
- {id} -> {user_id, name, description, total, date}
## UserAccount
- Служит для связи таблиц Users и Accounts.
- {id} -> {user_id, account_id}
## Deposit
- Хранит информацию о вкладах
- {id} -> {account_id, total_amount, start_at, end_at, interest_rate}
## Credit
- Хрнит информаци о кредитах на аккаунте на аккаунте
- {id} -> {account_id, amount, date_start, date_end, status, credit_type, monthly_payment}
## Debt
- {id} -> {user_id, total, date, status, description, creditor}




```mermaid
erDiagram
    user {
        id              uuid PK
        login           string
        username        string
        password_hash   string
        planned_budget  money
        avatar_url      uuid
    }

    account {
        id           uuid PK
        user_id      uuid FK
        bank_name    string
        balance      money
        description  text
        mean_payment string
    }
    
    category {
        id      uuid  PK
        user_id uuid  FK
        name    string
    }

    transaction {
        id          uuid  PK
        user_id     uuid  FK
        category_id uuid  FK
        account_id  uuid  FK
        is_income   boolean
        total       money
        date        date
        payer       string
        description string
    }
    
    user_account {
        id         uuid  PK
        user_id    uuid  FK
        account_id uuid  FK
    }

    goal {
        id          uuid  PK
        user_id     uuid  FK
        name        string
        description string
        amount      money
        start_date  date
    }
    investment {
        id uuid  PK
        user_id uuid  FK
        asset_type string
        asset_name string
        purchase_price money 
        purchase_date date
    }

    credit {
        id              uuid  PK
        account_id      uuid  FK        total_amount    money
        date_start      date
        date_end        date
        status          string
        credit_type     string
        monthly_payment money
    }
    
    deposit {
        id            uuid  PK
        account_id    uuid  FK
        total         money
        date_start    date
        date_end      date
        interest_rate numeric
    }

    account ||--o{ deposit : has
    account ||--o{ credit : has
    user ||--o{ investment : has
    user ||--o{ user_account : has
    account ||--o{ user_account : has
    user ||--o{ category : has
    user ||--o{ transaction : has
    account ||--o{ transaction : has
    category ||--o{ transaction : has
    user ||--o{ goal : has
    
```
