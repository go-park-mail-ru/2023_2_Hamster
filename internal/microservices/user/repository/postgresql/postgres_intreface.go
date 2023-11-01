package postgresql

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type FakeConn struct {
	sqlconn *pgx.Conn
}

type Fake interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func (p *FakeConn) QueryRow(ctx context.Context, sql string, args ...interface{}) {
	return p.sqlconn.Query()
}