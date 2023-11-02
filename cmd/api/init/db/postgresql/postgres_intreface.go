package postgresql

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type DbConn interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) (err error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) (err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)

	// Close(ctx context.Context) error
	// Config() *pgx.ConnConfig
	// ConnInfo() *pgtype.ConnInfo
	// CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	// Deallocate(ctx context.Context, name string) error
	// IsClosed() bool
	// PgConn() *pgconn.PgConn
	// Ping(ctx context.Context) error
	// Prepare(ctx context.Context, name string, sql string) (sd *pgconn.StatementDescription, err error)
	// QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
	// SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	// StatementCache() stmtcache.Cache
	// WaitForNotification(ctx context.Context) (*pgconn.Notification, error)
}
