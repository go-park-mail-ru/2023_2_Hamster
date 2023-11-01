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

type Conn struct {
	sqlconn *pgx.Conn
}

func (p *Conn) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return p.sqlconn.QueryRow(ctx, sql, args...)
}

func (p *Conn) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return p.sqlconn.Query(ctx, sql, args...)
}

func (p *Conn) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return p.sqlconn.Exec(ctx, sql, arguments...)
}

func (p *Conn) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.sqlconn.BeginTx(ctx, txOptions)
}

func (p *Conn) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.sqlconn.Begin(ctx)
}

func (p *Conn) BeginFunc(ctx context.Context, f func(pgx.Tx) error) (err error) {
	return p.sqlconn.BeginFunc(ctx, f)
}

// func (p *Conn) Close(ctx context.Context) error {
// 	return p.sqlconn.Close(ctx)
// }

// func (p *Conn) Config() *pgx.ConnConfig {
// 	return p.sqlconn.Config()
// }

// func (p *Conn) ConnInfo() *pgtype.ConnInfo {
// 	return p.sqlconn.ConnInfo()
// }

// func (p *Conn) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
// 	return p.sqlconn.CopyFrom(ctx, tableName, columnNames, rowSrc)
// }

// func (p *Conn) Deallocate(ctx context.Context, name string) error {
// 	return p.sqlconn.Deallocate(ctx, name)
// }

// func (p *Conn) IsClosed() bool {
// 	return p.sqlconn.IsClosed()
// }

// func (p *Conn) PgConn() *pgconn.PgConn {
// 	return p.sqlconn.PgConn()
// }

// func (p *Conn) Ping(ctx context.Context) error {
// 	return p.sqlconn.Ping(ctx)
// }

// func (p *Conn) Prepare(ctx context.Context, name string, sql string) (*pgconn.StatementDescription, error) {
// 	return p.sqlconn.Prepare(ctx, name, sql)
// }

// func (p *Conn) QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error) {
// 	return p.sqlconn.QueryFunc(ctx, sql, args, scans, f)
// }

// func (p *Conn) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
// 	return p.sqlconn.SendBatch(ctx, b)
// }

// func (p *Conn) StatementCache() stmtcache.Cache {
// 	return p.sqlconn.StatementCache()
// }

// func (p *Conn) WaitForNotification(ctx context.Context) (*pgconn.Notification, error) {
// 	return p.sqlconn.WaitForNotification(ctx)
// }
