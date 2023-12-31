package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	transactionCreate  = "INSERT INTO transaction (user_id, account_income, account_outcome, income, outcome, date, payer, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;"
	transactionGetFeed = `
    	SELECT 
			t.id, 
			t.user_id, 
			t.account_income, 
			t.account_outcome, 
			t.income, 
			t.outcome, 
			t.date, 
			t.payer, 
			t.description
		FROM Transaction t
		JOIN UserAccount ua ON t.account_income = ua.account_id
		WHERE ua.user_id = $1
	`

	transactionUpdate         = "UPDATE transaction set account_income=$2, account_outcome=$3, income=$4, outcome=$5, date=$6, payer=$7, description=$8 WHERE id = $1;"
	transactionGet            = "SELECT income, outcome, account_income, account_outcome FROM transaction WHERE id = $1;"
	TransactionGetUserByID    = "SELECT user_id FROM transaction WHERE id = $1;"
	transactionDelete         = "DELETE FROM transaction WHERE id = $1;"
	transactionGetCategory    = "SELECT tc.category_id, c.name AS category_name FROM TransactionCategory tc JOIN category c ON tc.category_id = c.id WHERE tc.transaction_id = $1;"
	transactionCreateCategory = "INSERT INTO transactionCategory (transaction_id, category_id) VALUES ($1, $2);"
	transactionDeleteCategory = "DELETE FROM transactionCategory WHERE transaction_id = $1;"
	transactionUpdateAccount  = "UPDATE accounts SET balance = balance - $1 WHERE id = $2;"
	transactionCheck          = "SELECT EXISTS( SELECT id FROM transaction WHERE id = $1);"
	transactionCount          = "SELECT COUNT(*) FROM transaction WHERE user_id = $1;"

	transactionGetFeedForExport = ` SELECT 
										t.id,  
										a_income.mean_payment AS account_income_name,
    									a_outcome.mean_payment AS account_outcome_name, 
										t.income, 
										t.outcome, 
										t.date, 
										t.payer, 
										t.description
									FROM Transaction t
									JOIN UserAccount ua ON t.account_income = ua.account_id
									JOIN Accounts a_income ON t.account_income = a_income.id
									JOIN Accounts a_outcome ON t.account_outcome = a_outcome.id
									WHERE ua.user_id = $1
	`
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

func (r *transactionRep) GetCount(ctx context.Context, user_id uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, transactionCount, user_id).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("[repo] %w", err)
	}

	return count, nil
}

func (r *transactionRep) GetFeed(ctx context.Context, user_id uuid.UUID, queryGet *models.QueryListOptions) ([]models.Transaction, error) {
	var transactions []models.Transaction
	count := 1
	var queryParamsSlice []interface{}

	query := transactionGetFeed
	queryParamsSlice = append(queryParamsSlice, user_id.String())

	if queryGet.Account != uuid.Nil {
		count++
		query += " AND (account_income = $" + strconv.Itoa(count) + " OR account_outcome = $" + strconv.Itoa(count) + ")"
		queryParamsSlice = append(queryParamsSlice, queryGet.Account.String())
	}

	if queryGet.Category != uuid.Nil {
		count++
		query += " AND id IN (SELECT transaction_id FROM TransactionCategory WHERE category_id = $" + strconv.Itoa(count) + ")"
		queryParamsSlice = append(queryParamsSlice, queryGet.Category.String())
	}

	if queryGet.Income && queryGet.Outcome {
		query += " AND income > 0 AND outcome > 0"
	}

	if !queryGet.Income && queryGet.Outcome {
		query += " AND outcome > 0 AND income = 0"
	}

	if queryGet.Income && !queryGet.Outcome {
		query += " AND income > 0 AND outcome = 0"
	}

	if !queryGet.StartDate.IsZero() || !queryGet.EndDate.IsZero() {
		count++
		if !queryGet.StartDate.IsZero() && !queryGet.EndDate.IsZero() {
			query += " AND date BETWEEN $" + strconv.Itoa(count) + " AND $" + strconv.Itoa(count+1)
			queryParamsSlice = append(queryParamsSlice, queryGet.StartDate, queryGet.EndDate)
		} else if !queryGet.StartDate.IsZero() {
			query += " AND date >= $" + strconv.Itoa(count)
			queryParamsSlice = append(queryParamsSlice, queryGet.StartDate)
		} else {
			query += " AND date <= $" + strconv.Itoa(count)
			queryParamsSlice = append(queryParamsSlice, queryGet.EndDate)
		}
	}

	query += " ORDER BY date DESC;"
	rows, err := r.db.Query(ctx, query, queryParamsSlice...)
	if err != nil {
		return nil, fmt.Errorf("[repo] %v", err)
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
			return nil, fmt.Errorf("[repo] %w", err)
		}

		categories, err := r.getCategoriesForTransaction(ctx, transaction.ID)
		if err != nil {
			return nil, fmt.Errorf("[repo] %w", err)
		}
		transaction.Categories = categories

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[repo] %w", err)
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("[repo] %w: %v", &models.NoSuchTransactionError{UserID: user_id}, err)
	}

	return transactions, nil
}

func (r *transactionRep) getCategoriesForTransaction(ctx context.Context, transactionID uuid.UUID) ([]models.CategoryName, error) {
	var categoryIDs []models.CategoryName

	rows, err := r.db.Query(ctx, transactionGetCategory, transactionID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var categoryID models.CategoryName
		if err := rows.Scan(&categoryID.ID, &categoryID.Name); err != nil {
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
	defer func() {
		if err != nil {
			if err = tx.Rollback(ctx); err != nil {
				r.logger.Fatal("Rollback transaction Error: %w", err)
			}

		}
	}()

	id, err := r.insertTransaction(ctx, tx, transaction)
	if err != nil {
		return id, err
	}

	if err = r.updateAccountBalances(ctx, tx, transaction); err != nil {
		return id, err
	}

	if err = r.insertCategories(ctx, tx, id, transaction.Categories); err != nil {
		return id, err
	}

	if err = tx.Commit(ctx); err != nil {
		return id, fmt.Errorf("[repo] failed to commit transaction: %w", err)
	}

	return id, nil
}

func (r *transactionRep) insertTransaction(ctx context.Context, tx pgx.Tx, transaction *models.Transaction) (uuid.UUID, error) {
	row := tx.QueryRow(ctx, transactionCreate,
		transaction.UserID,
		transaction.AccountIncomeID,
		transaction.AccountOutcomeID,
		transaction.Income,
		transaction.Outcome,
		transaction.Date,
		transaction.Payer,
		transaction.Description,
	)

	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return id, fmt.Errorf("[repo] failed create transaction: %w", err)
	}

	return id, nil
}

func (r *transactionRep) updateAccountBalances(ctx context.Context, tx pgx.Tx, transaction *models.Transaction) error {
	if err := r.updateAccountBalance(ctx, tx, transaction.AccountIncomeID, -transaction.Income); err != nil {
		return fmt.Errorf("[repo] failed to update old AccountIncome balance: %w", err)
	}

	if err := r.updateAccountBalance(ctx, tx, transaction.AccountOutcomeID, transaction.Outcome); err != nil {
		return fmt.Errorf("[repo] failed to update old AccountIncome balance: %w", err)
	}

	return nil
}

func (r *transactionRep) updateAccountBalance(ctx context.Context, tx pgx.Tx, accountID uuid.UUID, amount float64) error {
	_, err := tx.Exec(ctx, transactionUpdateAccount, amount, accountID)
	return err
}

func (r *transactionRep) insertCategories(ctx context.Context, tx pgx.Tx, transactionID uuid.UUID, categoryIDs []models.CategoryName) (err error) {
	for _, categoryID := range categoryIDs {
		if categoryID.ID == uuid.Nil {
			_, err = tx.Exec(ctx, transactionCreateCategory, transactionID, nil)
		} else {
			_, err = tx.Exec(ctx, transactionCreateCategory, transactionID, categoryID.ID)
		}
		if err != nil {
			return fmt.Errorf("[repo] failed to insert category association: %w", err)
		}
	}
	return nil
}

func (r *transactionRep) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("[repo] failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(ctx); err != nil {
				r.logger.Fatal("Rollback transaction Error: %w", err)
			}

		}
	}()

	existingIncome, existingOutcome, existingAccountIncomeID, existingAccountOutcomeID, err := r.getTransactionInfo(ctx, tx, transaction.ID)
	if err != nil {
		return err
	}

	if err = r.deleteAccountBalance(ctx, tx, existingIncome, existingOutcome, existingAccountIncomeID, existingAccountOutcomeID); err != nil {
		return err
	}

	if err = r.updateAccountBalances(ctx, tx, transaction); err != nil {
		return err
	}

	if err = r.updateTransactionInfo(ctx, tx, transaction); err != nil {
		return err
	}

	if err = r.deleteExistingCategoryAssociations(ctx, tx, transaction.ID); err != nil {
		return err
	}

	if err = r.insertCategories(ctx, tx, transaction.ID, transaction.Categories); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("[repo] failed to commit transaction: %w", err)
	}

	return nil
}

func (r *transactionRep) updateTransactionInfo(ctx context.Context, tx pgx.Tx, transaction *models.Transaction) error {
	_, err := tx.Exec(ctx, transactionUpdate,
		transaction.ID,
		transaction.AccountIncomeID,
		transaction.AccountOutcomeID,
		transaction.Income,
		transaction.Outcome,
		transaction.Date,
		transaction.Payer,
		transaction.Description,
	)
	if err != nil {
		return fmt.Errorf("[repo] failed to update transaction information: %w", err)
	}
	return nil
}

func (r *transactionRep) deleteExistingCategoryAssociations(ctx context.Context, tx pgx.Tx, transactionID uuid.UUID) error {
	_, err := tx.Exec(ctx, transactionDeleteCategory, transactionID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete existing category associations: %w", err)
	}
	return nil
}

func (r *transactionRep) deleteAccountBalance(ctx context.Context, tx pgx.Tx, existingIncome float64, existingOutcome float64, existingAccountIncomeID uuid.UUID, existingAccountOutcomeID uuid.UUID) error {
	if err := r.updateAccountBalance(ctx, tx, existingAccountIncomeID, existingIncome); err != nil {
		return fmt.Errorf("[repo] failed to update old AccountIncome balance: %w", err)
	}

	if err := r.updateAccountBalance(ctx, tx, existingAccountOutcomeID, -existingOutcome); err != nil {
		return fmt.Errorf("[repo] failed to update old AccountIncome balance: %w", err)
	}

	return nil
}

func (r *transactionRep) getTransactionInfo(ctx context.Context, tx pgx.Tx, transactionID uuid.UUID) (float64, float64, uuid.UUID, uuid.UUID, error) {
	var existingIncome, existingOutcome float64
	var existingAccountIncomeID, existingAccountOutcomeID uuid.UUID

	row := tx.QueryRow(ctx, transactionGet, transactionID)
	err := row.Scan(&existingIncome, &existingOutcome, &existingAccountIncomeID, &existingAccountOutcomeID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0, uuid.Nil, uuid.Nil, fmt.Errorf("[repo] %w: %v", &models.NoSuchTransactionError{UserID: transactionID}, err)
	} else if err != nil {
		return 0, 0, uuid.Nil, uuid.Nil, fmt.Errorf("[repo] failed request db %s, %w", transactionGet, err)
	}

	return existingIncome, existingOutcome, existingAccountIncomeID, existingAccountOutcomeID, nil
}

func (r *transactionRep) DeleteTransaction(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("[repo] failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(ctx); err != nil {
				r.logger.Fatal("Rollback transaction Error: %w", err)
			}

		}
	}()

	existingIncome, existingOutcome, existingAccountIncomeID, existingAccountOutcomeID, err := r.getTransactionInfo(ctx, tx, transactionID)
	if err != nil {
		return err
	}

	if err = r.deleteAccountBalance(ctx, tx, existingIncome, existingOutcome, existingAccountIncomeID, existingAccountOutcomeID); err != nil {
		return err
	}

	if err = r.deleteExistingCategoryAssociations(ctx, tx, transactionID); err != nil {
		return err
	}

	_, err = tx.Exec(ctx, transactionDelete, transactionID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete transaction %s, %w", transactionDelete, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("[repo] failed to commit transaction: %w", err)
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

func (r *transactionRep) GetTransactionForExport(ctx context.Context, userId uuid.UUID, queryGet *models.QueryListOptions) ([]models.TransactionExport, error) {
	var transactions []models.TransactionExport
	count := 1
	var queryParamsSlice []interface{}

	query := transactionGetFeedForExport
	queryParamsSlice = append(queryParamsSlice, userId.String())

	if !queryGet.StartDate.IsZero() || !queryGet.EndDate.IsZero() {
		count++
		if !queryGet.StartDate.IsZero() && !queryGet.EndDate.IsZero() {
			query += " AND date BETWEEN $" + strconv.Itoa(count) + " AND $" + strconv.Itoa(count+1)
			queryParamsSlice = append(queryParamsSlice, queryGet.StartDate, queryGet.EndDate)
		} else if !queryGet.StartDate.IsZero() {
			query += " AND date >= $" + strconv.Itoa(count)
			queryParamsSlice = append(queryParamsSlice, queryGet.StartDate)
		} else {
			query += " AND date <= $" + strconv.Itoa(count)
			queryParamsSlice = append(queryParamsSlice, queryGet.EndDate)
		}
	}

	query += " ORDER BY date DESC;"

	rows, err := r.db.Query(ctx, query, queryParamsSlice...)
	if err != nil {
		return nil, fmt.Errorf("[repo] %v", err)
	}

	for rows.Next() {
		var transaction models.TransactionExport
		if err := rows.Scan(
			&transaction.ID,
			&transaction.AccountIncome,
			&transaction.AccountOutcome,
			&transaction.Income,
			&transaction.Outcome,
			&transaction.Date,
			&transaction.Payer,
			&transaction.Description,
		); err != nil {
			return nil, fmt.Errorf("[repo] %w", err)
		}

		categories, err := r.getCategoriesForTransaction(ctx, transaction.ID)
		if err != nil {
			return nil, fmt.Errorf("[repo] %w", err)
		}
		for _, data := range categories {
			transaction.Categories = append(transaction.Categories, data.Name)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
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
