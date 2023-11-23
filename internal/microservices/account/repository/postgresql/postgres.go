package postgresql

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/models"
	"github.com/google/uuid"
)

const (
	AccountGetUserByID = `SELECT EXISTS(
    						SELECT 1 FROM UserAccount 
    						WHERE account_id = $1 AND user_id = $2);`

	AccountUpdate     = "UPDATE accounts SET balance = $1, accumulation = $2, balance_enabled = $3, mean_payment = $4;"
	AccountDelete     = "DELETE FROM account WHERE id = $1;"
	UserAccountDelete = "DELETE FROM userAccount WHERE account_id = $1;"
	AccountCreate     = "INSERT INTO account (balance, accumulation, balance_enabled, mean_paymment) VALUES ($1, $2, $3, $4);"
	AccountUserCreate = "INSERT INTO userAccount (user_id, account_id) VALUES ($1, $2);"
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

func (r *AccountRep) CheckForbidden(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) error { // need test
	var result bool
	row := r.db.QueryRow(ctx, AccountGetUserByID, transactionID, userID)

	err := row.Scan(&result)

	return err
}

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

	row := tx.QueryRow(ctx, AccountCreate, account.Balance, account.Accumulation, account.BalanceEnabled, account.MeanPayment)
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

func (r *AccountRep) UpdateAccount(ctx context.Context, userID uuid.UUID, account *models.Accounts) error {
	_, err := r.db.Exec(ctx, AccountUpdate, account.Balance, account.Accumulation, account.BalanceEnabled, account.MeanPayment)
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

	_, err = tx.Exec(ctx, UserAccountDelete, accountID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete from UserAccount table: %w", err)
	}

	_, err = tx.Exec(ctx, AccountDelete, accountID)
	if err != nil {
		return fmt.Errorf("[repo] failed to delete transaction %s, %w", AccountDelete, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("[repo] failed to commit transaction: %w", err)
	}

	return nil
}
