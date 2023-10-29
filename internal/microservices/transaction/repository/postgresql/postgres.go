package postgresql

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgtype/pgxtype"
)

const (
	transactionCreate               = "INSERT INTO transaction (user_id, category_id, account_id, total, is_income, date, payer, description) VALUES ($1, $2, $3) RETURNING id;"
	transactionGetFeed              = "SELECT * FROM transaction WHERE user_id = $1"
	transactionUpdateBalanceAccount = `UPDATE accounts
										SET balance = CASE
											WHEN is_income = true THEN balance + $2
											WHEN is_income = false THEN balance - $2
											ELSE balance
											END
										WHERE account_id = $1;`
)

type transactionRep struct {
	db     pgxtype.Querier
	logger logger.CustomLogger
}

func NewRepository(db pgxtype.Querier, l logger.CustomLogger) *transactionRep {
	return &transactionRep{
		db:     db,
		logger: l,
	}
}

func (r *transactionRep) GetFeed(ctx context.Context, user_id uuid.UUID) ([]models.Transaction, error) { // need test
	var transactions []models.Transaction
	rows, err := r.db.Query(ctx, transactionGetFeed, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account models.Transaction
		if err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.CategoryID,
			&account.AccountID,
			&account.Total,
			&account.IsIncome,
			&account.Date,
			&account.Payer,
			&account.Description,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("[repo] %w: %v", &models.NoSuchTransactionError{UserID: user_id}, err)
	}

	return transactions, nil
}

func (r *transactionRep) CreateTransaction(ctx context.Context, transaction *models.Transaction) (uuid.UUID, error) { // need test
	row := r.db.QueryRow(ctx, transactionGetFeed,
		transaction.UserID,
		transaction.CategoryID,
		transaction.AccountID,
		transaction.Total,
		transaction.IsIncome,
		transaction.Date,
		transaction.Payer,
		transaction.Description)

	var id uuid.UUID

	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("[repo] %w", err)
	}

	_, err = r.db.Exec(ctx, transactionUpdateBalanceAccount, transaction.AccountID, transaction.Total)

	if err != nil {
		return id, fmt.Errorf("[repo] %v", err)
	}

	return id, nil
}
