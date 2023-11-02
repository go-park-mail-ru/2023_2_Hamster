package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

const (
	transactionCreate  = "INSERT INTO transaction (user_id, account_income, account_outcome, income, outcome, date, payer, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;"
	transactionGetFeed = `SELECT * 
						FROM (
							SELECT * 
							FROM transaction 
							WHERE user_id = $1
							LIMIT $2 
							OFFSET $3
						) AS subquery
						ORDER BY date DESC;`

	transactionUpdate         = "UPDATE transaction set account_income=$2, account_outcome=$3, income=$4, outcome=$5, date=$6, payer=$7, description=$8 WHERE id = $1;"
	transactionGet            = "SELECT income, outcome, account_income, account_outcome FROM transaction WHERE id = $1;"
	TransactionGetUserByID    = "SELECT user_id FROM transaction WHERE id = $1;"
	transactionDelete         = "DELETE FROM transaction WHERE id = $1;"
	transactionGetCategory    = "SELECT category_id FROM TransactionCategory WHERE transaction_id = $1;"
	transactionCreateCategory = "INSERT INTO transactionCategory (transaction_id, category_id) VALUES ($1, $2);"
	transactionDeleteCategory = "DELETE FROM transactionCategory WHERE transaction_id = $1;"
	transacitonUpdateAccount  = "UPDATE accounts SET balance = balance - $1 WHERE id = $2;"
	transactionCheck          = "SELECT EXISTS( SELECT id FROM transaction WHERE id = $1);"
)

type transactionRep struct {
	db     postgresql.DbConn
	logger logger.Logger
}

func NewRepository(db postgresql.DbConn, l logger.Logger) *transactionRep {
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
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.AccountIncomeID,
			&transaction.AccountOutcomeID,
			&transaction.Income,
			&transaction.Outcome,
			&transaction.Date,
			&transaction.Payer,
			&transaction.Description,
		); err != nil {
			return nil, false, err
		}
		categories, err := r.getCategoriesForTransaction(ctx, transaction.ID)
		if err != nil {
			return nil, false, err
		}
		transaction.Categories = categories

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	if len(transactions) == 0 {
		return nil, false, fmt.Errorf("[repo] %w: %v", &models.NoSuchTransactionError{UserID: user_id}, err)
	}

	return transactions, len(transactions) < pageSize, nil
}

func (r *transactionRep) getCategoriesForTransaction(ctx context.Context, transactionID uuid.UUID) ([]uuid.UUID, error) {
	var categoryIDs []uuid.UUID

	rows, err := r.db.Query(ctx, transactionGetCategory, transactionID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var categoryID uuid.UUID
		if err := rows.Scan(&categoryID); err != nil {
			return nil, err
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categoryIDs, nil
}

func (r *transactionRep) CreateTransaction(ctx context.Context, transaction *models.Transaction) (uuid.UUID, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, transactionCreate,
		transaction.UserID,
		transaction.AccountIncomeID,
		transaction.AccountOutcomeID,
		transaction.Income,
		transaction.Outcome,
		transaction.Date,
		transaction.Payer,
		transaction.Description)
	var id uuid.UUID

	err = row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("[repo] failed create transaction: %w", err)
	}

	_, err = tx.Exec(ctx, transacitonUpdateAccount, -transaction.Income, transaction.AccountIncomeID)
	if err != nil {
		return id, fmt.Errorf("failed to update old AccountIncome balance: %w", err)
	}

	_, err = tx.Exec(ctx, transacitonUpdateAccount, transaction.Outcome, transaction.AccountOutcomeID)
	if err != nil {
		return id, fmt.Errorf("failed to update old AccountIncome balance: %w", err)
	}

	for _, categoryID := range transaction.Categories {
		_, err = tx.Exec(ctx, transactionCreateCategory, id, categoryID)
		if err != nil {
			return id, fmt.Errorf("[repo] failed to insert category association: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return id, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

func (r *transactionRep) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var existingIncome, existingOutcome float64
	var existingAccountIncomeID, existingAccountOutcomeID uuid.UUID

	row := tx.QueryRow(ctx, transactionGet, transaction.ID)
	err = row.Scan(&existingIncome, &existingOutcome, &existingAccountIncomeID, &existingAccountOutcomeID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("[repo] %w: %v", &models.NoSuchTransactionError{UserID: transaction.ID}, err)
	} else if err != nil {
		return fmt.Errorf("[repo] failed request db %s, %w", transactionGet, err)
	}

	_, err = tx.Exec(ctx, transacitonUpdateAccount, existingIncome, existingAccountIncomeID)
	if err != nil {
		return fmt.Errorf("failed to update old AccountIncome balance: %s, %w", transacitonUpdateAccount, err)
	}

	_, err = tx.Exec(ctx, transacitonUpdateAccount, -existingOutcome, existingAccountOutcomeID)
	if err != nil {
		return fmt.Errorf("failed to update old AccountOutcome balance: %s, %w", transacitonUpdateAccount, err)
	}

	_, err = tx.Exec(ctx, transacitonUpdateAccount, -transaction.Income, transaction.AccountIncomeID)
	if err != nil {
		return fmt.Errorf("failed to update new AccountIncome balance: %s, %w", transacitonUpdateAccount, err)
	}

	_, err = tx.Exec(ctx, transacitonUpdateAccount, transaction.Outcome, transaction.AccountOutcomeID)
	if err != nil {
		return fmt.Errorf("failed to update new AccountOutcome balance: %s, %w", transacitonUpdateAccount, err)
	}

	_, err = tx.Exec(ctx, transactionUpdate,
		transaction.ID,
		transaction.AccountIncomeID,
		transaction.AccountOutcomeID,
		transaction.Income,
		transaction.Outcome,
		transaction.Date,
		transaction.Payer,
		transaction.Description)
	if err != nil {
		return fmt.Errorf("[repo] failed to update transaction information: %s, %w", transactionUpdate, err)
	}

	_, err = tx.Exec(ctx, transactionDeleteCategory, transaction.ID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete existing category associations: %s, %w", transactionUpdate, err)
	}

	for _, categoryID := range transaction.Categories {
		_, err = tx.Exec(ctx, transactionCreateCategory, transaction.ID, categoryID)
		if err != nil {
			return fmt.Errorf("[repo] failed to insert category associations: %s, %w", transactionUpdate, err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *transactionRep) CheckForbidden(ctx context.Context, transactionID uuid.UUID) (uuid.UUID, error) { // need test
	var userID uuid.UUID
	row := r.db.QueryRow(ctx, TransactionGetUserByID, transactionID)

	err := row.Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return userID, fmt.Errorf("[repo] %w: %v", &models.NoSuchTransactionError{UserID: transactionID}, err)
	} else if err != nil {
		return userID,
			fmt.Errorf("[repo] failed request db %s, %w", TransactionGetUserByID, err)
	}
	return userID, nil
}

func (r *transactionRep) DeleteTransaction(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var existingIncome, existingOutcome float64
	var existingAccountIncomeID, existingAccountOutcomeID uuid.UUID

	row := tx.QueryRow(ctx, transactionGet, transactionID)
	err = row.Scan(&existingIncome, &existingOutcome, &existingAccountIncomeID, &existingAccountOutcomeID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("[repo] %w: %v", &models.NoSuchTransactionError{UserID: transactionID}, err)
	} else if err != nil {
		return fmt.Errorf("[repo] failed request db %s, %w", transactionGet, err)
	}

	_, err = tx.Exec(ctx, transacitonUpdateAccount, existingIncome, existingAccountIncomeID)
	if err != nil {
		return fmt.Errorf("[repo] failed to update old AccountIncome balance: %w", err)
	}

	_, err = tx.Exec(ctx, transacitonUpdateAccount, -existingOutcome, existingAccountOutcomeID)
	if err != nil {
		return fmt.Errorf("[repo] failed to update old AccountOutcome balance: %w", err)
	}

	_, err = tx.Exec(ctx, transactionDeleteCategory, transactionID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete existing category associations: %s, %w", transactionUpdate, err)
	}

	_, err = tx.Exec(ctx, transactionDelete, transactionID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete transaction %s, %w", transactionDelete, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// func (r *transactionRep) Check(ctx context.Context, transactionID uuid.UUID) error {
// 	var exists bool
// 	err := r.db.QueryRow(ctx, transactionCheck, transactionID).Scan(&exists)
// 	if err != nil {
// 		return fmt.Errorf("(repo) failed to exec query: %s, %w", transactionCheck, err)
// 	}

// 	if !exists {
// 		return fmt.Errorf("(repo) %w: %w", &models.NoSuchTransactionError{UserID: transactionID}, err)
// 	}

// 	return nil
// }
