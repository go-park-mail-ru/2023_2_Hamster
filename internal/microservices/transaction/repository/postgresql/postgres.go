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

	transactionUpdateBalanceAccount = `UPDATE accounts
						SET balance = CASE
							WHEN id = $1 THEN balance + $2
							WHEN id = $3 THEN balance - $4
							ELSE balance
							END
						WHERE id IN ($1, $3);`

	transactionUpdate         = "UPDATE transaction set account_income=$2, account_outcome=$3, income=$4, outcome=$5, date=$6, payer=$7, description=$8 WHERE id = $1;"
	transactionGet            = "SELECT income, outcome, account_income, account_outcome FROM transaction WHERE id = $1;"
	TransactionGetUserByID    = "SELECT user_id FROM transaction WHERE id = $1;"
	transactionDelete         = "DELETE FROM transaction WHERE $1 = id;"
	transactionGetCategory    = "SELECT category_id FROM TransactionCategory WHERE transaction_id = $1;"
	transactionCreateCategory = "INSERT INTO transactionCategory (transaction_id, category_id) VALUES ($1, $2);"
	transactionDeleteCategory = "DELETE FROM transactionCategory WHERE transaction_id;"
	transacitonUpdateOld      = "UPDATE accounts SET balance = balance - $1 WHERE id = $2;"
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

func (r *transactionRep) CreateTransaction(ctx context.Context, transaction *models.Transaction) (uuid.UUID, error) { // need test
	row := r.db.QueryRow(ctx, transactionCreate,
		transaction.UserID,
		transaction.AccountIncomeID,
		transaction.AccountOutcomeID,
		transaction.Income,
		transaction.Outcome,
		transaction.Date,
		transaction.Payer,
		transaction.Description)
	var id uuid.UUID

	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("[repo] failed create transaction %w", err)
	}

	_, err = r.db.Exec(ctx, transactionUpdateBalanceAccount,
		transaction.AccountIncomeID, transaction.Income,
		transaction.AccountOutcomeID, transaction.Outcome)
	if err != nil {
		return id, fmt.Errorf("[repo] failed to update account balances: %w", err)
	}

	for _, categoryID := range transaction.Categories {
		_, err = r.db.Exec(ctx, transactionCreateCategory, id, categoryID)
		if err != nil {
			return id, fmt.Errorf("[repo] failed to insert category association: %w", err)
		}
	}
	return id, nil
}

func (r *transactionRep) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	var existingIncome, existingOutcome float64
	var existingAccountIncomeID, existingAccountOutcomeID uuid.UUID

	row := r.db.QueryRow(ctx, transactionGet, transaction.ID)
	err := row.Scan(&existingIncome, &existingOutcome, &existingAccountIncomeID, &existingAccountOutcomeID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("[repo] %w: %v", &models.NoSuchTransactionError{UserID: transaction.ID}, err)
	} else if err != nil {
		return fmt.Errorf("[repo] failed request db %s, %w", transactionGet, err)
	}

	_, err = r.db.Exec(ctx, transactionUpdateBalanceAccount, -existingIncome, existingAccountIncomeID)
	if err != nil {
		return fmt.Errorf("failed to update old AccountIncome balance: %w", err)
	}

	_, err = r.db.Exec(ctx, transactionUpdateBalanceAccount, -existingOutcome, existingAccountOutcomeID)
	if err != nil {
		return fmt.Errorf("failed to update old AccountOutcome balance: %w", err)
	}

	_, err = r.db.Exec(ctx, transactionUpdateBalanceAccount, transaction.Income, transaction.AccountIncomeID)
	if err != nil {
		return fmt.Errorf("failed to update new AccountIncome balance: %w", err)
	}

	_, err = r.db.Exec(ctx, transactionUpdateBalanceAccount, transaction.Outcome, transaction.AccountOutcomeID)
	if err != nil {
		return fmt.Errorf("failed to update new AccountOutcome balance: %w", err)
	}

	_, err = r.db.Exec(ctx, transactionUpdate,
		transaction.AccountIncomeID,
		transaction.AccountOutcomeID,
		transaction.Income,
		transaction.Outcome,
		transaction.Date,
		transaction.Payer,
		transaction.Description)
	if err != nil {
		return fmt.Errorf("[repo] failed to update transaction information: %w", err)
	}

	_, err = r.db.Exec(ctx, transactionDeleteCategory, transaction.ID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete existing category associations: %w", err)
	}

	for _, categoryID := range transaction.Categories {
		_, err = r.db.Exec(ctx, transactionCreateCategory, transaction.ID, categoryID)
		if err != nil {
			return fmt.Errorf("[repo] failed to insert category associations: %w", err)
		}
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
	var existingIncome, existingOutcome float64
	var existingAccountIncomeID, existingAccountOutcomeID uuid.UUID

	row := r.db.QueryRow(ctx, transactionGet, transactionID)
	err := row.Scan(&existingIncome, &existingOutcome, &existingAccountIncomeID, &existingAccountOutcomeID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("[repo] %w: %v", &models.NoSuchTransactionError{UserID: transactionID}, err)
	} else if err != nil {
		return fmt.Errorf("[repo] failed request db %s, %w", transactionGet, err)
	}

	_, err = r.db.Exec(ctx, transactionUpdateBalanceAccount, -existingIncome, existingAccountIncomeID)
	if err != nil {
		return fmt.Errorf("failed to update old AccountIncome balance: %w", err)
	}

	_, err = r.db.Exec(ctx, transactionUpdateBalanceAccount, -existingOutcome, existingAccountOutcomeID)
	if err != nil {
		return fmt.Errorf("failed to update old AccountOutcome balance: %w", err)
	}

	return nil
}
