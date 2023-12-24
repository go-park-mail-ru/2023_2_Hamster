package postgresql

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	AccountGetUserByID = ` SELECT * FROM UserAccount 
    						WHERE account_id = $1 AND user_id = $2;`

	DeleteTransactionCategory = `
		DELETE FROM TransactionCategory
		WHERE transaction_id IN (SELECT id FROM Transaction WHERE user_id = $1 AND (account_income = $2 OR account_outcome = $2));
	`

	DeleteUserTransactions = `
		DELETE FROM Transaction
		WHERE user_id = $1 AND (account_income = $2 OR account_outcome = $2);
	`

	AccountSharingCheck       = `SELECT COUNT(*) FROM accounts WHERE sharing_id = $1 AND id = $2;`
	AccountUpdate             = "UPDATE accounts SET balance = $1, accumulation = $2, balance_enabled = $3, mean_payment = $4 WHERE id = $5;"
	AccountDelete             = "DELETE FROM accounts WHERE id = $1;"
	UserAccountDelete         = "DELETE FROM userAccount WHERE account_id = $1;"
	AccountCreate             = "INSERT INTO accounts (balance, accumulation, balance_enabled, mean_payment, sharing_id) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
	AccountUserCreate         = "INSERT INTO userAccount (user_id, account_id) VALUES ($1, $2);"
	TransactionCategoryDelete = "DELETE FROM TransactionCategory WHERE transaction_id IN (SELECT id FROM Transaction WHERE account_income = $1 OR account_outcome = $1)"
	AccountTransactionDelete  = "DELETE FROM Transaction WHERE account_income = $1 OR account_outcome = $1"
	Unsubscribe               = "DELETE FROM userAccount WHERE account_id = $1 AND user_id = $2"
)

type AccountRep struct {
	db     postgresql.DbConn
	logger logger.Logger
}

func NewRepository(db postgresql.DbConn, log logger.Logger) *AccountRep {
	return &AccountRep{
		db:     db,
		logger: log,
	}
}

func (r *AccountRep) SharingCheck(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error {
	var count int
	row := r.db.QueryRow(ctx, AccountSharingCheck, userID, accountID)

	err := row.Scan(&count)
	if err != nil {
		return fmt.Errorf("[repo] failed %w, %v", &models.ForbiddenUserError{}, err)
	}
	if count == 0 {
		return fmt.Errorf("[repo] failed %w", &models.ForbiddenUserError{})
	}

	return nil
}

func (r *AccountRep) Unsubscribe(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("[repo] failed to start transaction: %w", err)
	}

	_, err = tx.Exec(ctx, DeleteTransactionCategory, userID, accountID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("[repo] failed to delete transaction category: %w", err)
	}

	_, err = tx.Exec(ctx, DeleteUserTransactions, userID, accountID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("[repo] failed to delete transaction table: %w", err)
	}

	_, err = tx.Exec(ctx, Unsubscribe, accountID, userID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("[repo] failed to delete from UserAccount table: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("[repo] failed to commit transaction: %w", err)
	}

	return nil
}

func (r *AccountRep) DeleteUserInAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	_, err := r.db.Exec(ctx, Unsubscribe, accountID, userID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete from UserAccount table: %w", err)
	}
	return nil
}

func (r *AccountRep) CheckDuplicate(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	row := r.db.QueryRow(ctx, AccountGetUserByID, accountID, userID)

	err := row.Scan(&userID, &accountID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		} else {
			return fmt.Errorf("[repo] query error: %w", err)
		}
	}

	return &models.DuplicateError{}
}

func (r *AccountRep) CheckForbidden(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error {
	var result bool
	row := r.db.QueryRow(ctx, AccountGetUserByID, accountID, userID)

	err := row.Scan(&result)

	return err
}

// (balance, accumulation, balance_enabled, mean_paymment)
func (r *AccountRep) CreateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) (uuid.UUID, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("[repo] failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(ctx); err != nil {
				r.logger.Fatal("Rollback account Error: %w", err)
			}

		}
	}()

	row := tx.QueryRow(ctx, AccountCreate, account.Balance, account.Accumulation, account.BalanceEnabled, account.MeanPayment, userID)
	var id uuid.UUID

	err = row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("[repo] error request %s, %w", AccountCreate, err)
	}

	_, err = tx.Exec(ctx, AccountUserCreate, userID, id)
	if err != nil {
		return id, fmt.Errorf("[repo] can't create accountUser %s, %w", AccountUserCreate, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return id, fmt.Errorf("[repo] failed to commit account: %w", err)
	}

	return id, nil
}

func (r *AccountRep) AddUserInAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	_, err := r.db.Exec(ctx, AccountUserCreate, userID, accountID)
	if err != nil {
		return fmt.Errorf("[repo] can't create accountUser %s, %w", AccountUserCreate, err)
	}
	return nil
}

func (r *AccountRep) UpdateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) error {
	_, err := r.db.Exec(ctx, AccountUpdate, account.Balance, account.Accumulation, account.BalanceEnabled, account.MeanPayment, account.ID)
	if err != nil {
		return fmt.Errorf("[repo] failed update account %w", err)
	}

	return nil
}

func (r *AccountRep) DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("[repo] failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(ctx); err != nil {
				r.logger.Fatal("Rollback account Error: %w", err)
			}

		}
	}()

	_, err = tx.Exec(ctx, TransactionCategoryDelete, accountID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete from TransactionCategory table: %w", err)
	}

	_, err = tx.Exec(ctx, AccountTransactionDelete, accountID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete from Transaction table: %w", err)
	}

	_, err = tx.Exec(ctx, UserAccountDelete, accountID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete from UserAccount table: %w", err)
	}

	_, err = tx.Exec(ctx, AccountDelete, accountID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete account %s, %w", AccountDelete, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("[repo] failed to commit account: %w", err)
	}

	return nil
}
