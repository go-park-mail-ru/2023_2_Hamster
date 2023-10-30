package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgtype/pgxtype"
)

const (
	transactionCreate  = "INSERT INTO transaction (user_id, category_id, account_id, total, is_income, date, payer, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;"
	transactionGetFeed = `SELECT * 
						FROM (
							SELECT * 
							FROM transaction 
							WHERE user_id = $1
							LIMIT $2 
							OFFSET $3
						) AS subquery
						ORDER BY date DESC;`
	transactionUpdateBalanceAccount = `UPDATE accounts
										SET balance = CASE
											WHEN $3 = true THEN balance + $2
											WHEN $3 = false THEN balance - $2
											ELSE balance
											END
										WHERE id = $1;`
	transactionUpdate = "UPDATE transaction set category_id=$2, account_id=$3, total=$4, is_income=$5, date=$6, payer=$7, description=$8 WHERE id = $1;"
	transactionGet    = "SELECT total, is_income, account_id FROM transaction WHERE id = $1"
	transactionCheck  = "SELECT EXISTS(SELECT id FROM transaction WHERE id = $1 AND user_id = $2);"
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

func (r *transactionRep) GetFeed(ctx context.Context, user_id uuid.UUID, page int, pageSize int) ([]models.Transaction, bool, error) { // need test
	var transactions []models.Transaction
	offset := (page - 1) * pageSize

	rows, err := r.db.Query(ctx, transactionGetFeed, user_id, pageSize, offset)
	if err != nil {
		return nil, false, err
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
			return nil, false, err
		}
		transactions = append(transactions, account)
	}

	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	if len(transactions) == 0 {
		return nil, false, fmt.Errorf("[repo] %w: %v", &models.NoSuchTransactionError{UserID: user_id}, err)
	}

	return transactions, len(transactions) < pageSize, nil
}

func (r *transactionRep) CreateTransaction(ctx context.Context, transaction *models.Transaction) (uuid.UUID, error) { // need test
	row := r.db.QueryRow(ctx, transactionCreate,
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
		return id, fmt.Errorf("[repo] failed create transaction %w", err)
	}

	_, err = r.db.Exec(ctx, transactionUpdateBalanceAccount, transaction.AccountID, transaction.Total, transaction.IsIncome)

	if err != nil {
		return id, fmt.Errorf("[repo] failed update account %v", err)
	}

	return id, nil
}

func (r *transactionRep) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	row := r.db.QueryRow(ctx, transactionGet, transaction.ID)

	var total float64
	var isIncome bool
	var accountID uuid.UUID
	err := row.Scan(&total, &isIncome, &accountID)
	fmt.Println(total, isIncome, accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("[repo] %w: %v", &models.NoSuchTransactionError{UserID: transaction.ID}, err)
	} else if err != nil {
		return fmt.Errorf("[repo] failed request db %s, %w", transactionGet, err)
	}

	_, err = r.db.Exec(ctx, transactionUpdateBalanceAccount, accountID, total, !isIncome)
	if err != nil {
		return fmt.Errorf("[repo] failed update transaction %s, %w", transactionUpdateBalanceAccount, err)
	}

	_, err = r.db.Exec(ctx, transactionUpdate,
		transaction.ID,
		transaction.CategoryID,
		transaction.AccountID,
		transaction.Total,
		transaction.IsIncome,
		transaction.Date,
		transaction.Payer,
		transaction.Description)

	if err != nil {
		return fmt.Errorf("[repo] failed update transaction  %s, %w", transactionUpdate, err)
	}

	_, err = r.db.Exec(ctx, transactionUpdateBalanceAccount, transaction.AccountID, transaction.Total, transaction.IsIncome)

	if err != nil {
		return fmt.Errorf("[repo] failed update account %s, %w", transactionUpdateBalanceAccount, err)
	}
	return nil
}

func (r *transactionRep) CheckTransaciont(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) error {
	var exists bool
	err := r.db.QueryRow(ctx, transactionCheck, transactionID, userID).Scan(&exists)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("[repo] %w: %v", &models.NoSuchTransactionError{UserID: transactionID}, err)
	}

	if err != nil {
		return fmt.Errorf("[repo] %s: %v", transactionCheck, err)
	}

	return nil
}
