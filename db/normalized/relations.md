# Описание Таблиц
## User
- Хранит информацию о пользователях.
- {id} -> {username, login, password_hash, salt, planned_budget, avatar_url}
- {login}->{username, password_hash, salt, planned_budget, avatar_url}
## Account
- Хранит информацию о банковских счетах пользователя.
- {id} -> {user_id, balance, mean_payment}
- {user_id} -> {id}
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





```mermaid
erDiagram
    User {
        uuid id PK
        string username
        string login
        string password_hash
        string salt
        numeric planned_budget
        uuid avatar_url
    }
    Account {
        uuid id PK
        string user_id FK
        numeric balance
        string mean_payment
    }
    Category {
        uuid id PK
        uuid user_id FK
        string name
    }
    Transaction {
        uuid id PK
        uuid user_id FK
        uuid category_id FK
        uuid account_id FK
        numeric total
        boolean is_income
        date date
        string payer
        string description
    }
    UserAccount {
        uuid id PK
        uuid user_id FK
        uuid account_id FK
    }
    Goal {
        uuid id PK
        uuid user_id FK
        string name
        string description
        numeric total
        date date
    }
    Investment {
        uuid id PK
        uuid user_id FK
        name string
        total numeric
        date_start date
        price numeric
        percentage numeric
    }
    User ||--o{ Investment : has
    User ||--o{ UserAccount : has
    Account ||--o{ UserAccount : has
    User ||--o{ Category : has
    User ||--o{ Transaction : has
    Account ||--o{ Transaction : has
    Category ||--o{ Transaction : has
    User ||--o{ Goal : has

    
    
```
