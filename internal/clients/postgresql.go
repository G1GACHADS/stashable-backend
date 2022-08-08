package clients

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreSQLConnection struct {
	ctx        context.Context
	connection *pgxpool.Pool
}

func NewPostgreSQLClient(ctx context.Context, connString string) (*PostgreSQLConnection, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	connection, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &PostgreSQLConnection{
		ctx:        ctx,
		connection: connection,
	}, nil
}

func (p *PostgreSQLConnection) Exec(sql string, args ...interface{}) (int64, error) {
	result, err := p.connection.Exec(p.ctx, sql, args...)
	return result.RowsAffected(), err
}

func (p *PostgreSQLConnection) Query(sql string, args ...interface{}) (*ConnRows, error) {
	rows, err := p.connection.Query(p.ctx, sql, args...)
	return newDatabaseRows(rows), err
}

func (p *PostgreSQLConnection) QueryRow(sql string, args ...interface{}) *ConnRow {
	row := p.connection.QueryRow(p.ctx, sql, args...)
	return newDatabaseRow(row)
}

func (p *PostgreSQLConnection) BeginTx() (*ConnTx, error) {
	tx, err := p.connection.BeginTx(p.ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	return newDatabaseTransaction(p.ctx, tx), nil
}

type ConnRow struct {
	row pgx.Row
}

func newDatabaseRow(row pgx.Row) *ConnRow {
	return &ConnRow{
		row: row,
	}
}

func (p *ConnRow) Scan(dest ...interface{}) error {
	return p.row.Scan(dest...)
}

type ConnRows struct {
	rows pgx.Rows
}

func newDatabaseRows(rows pgx.Rows) *ConnRows {
	return &ConnRows{
		rows: rows,
	}
}

func (p *ConnRows) Scan(dest ...interface{}) error {
	return p.rows.Scan(dest...)
}

func (p *ConnRows) Next() bool {
	return p.rows.Next()
}

func (p *ConnRows) Close() {
	p.rows.Close()
}

func (p *ConnRows) Err() error {
	return p.rows.Err()
}

type ConnTx struct {
	ctx context.Context
	tx  pgx.Tx
}

func newDatabaseTransaction(ctx context.Context, tx pgx.Tx) *ConnTx {
	return &ConnTx{
		ctx: ctx,
		tx:  tx,
	}
}

func (p *ConnTx) Commit() error {
	return p.tx.Commit(p.ctx)
}

func (p *ConnTx) Rollback() error {
	return p.tx.Rollback(p.ctx)
}

func (p *ConnTx) Exec(sql string, args ...interface{}) (int64, error) {
	result, err := p.tx.Exec(p.ctx, sql, args...)
	return result.RowsAffected(), err
}

func (p *ConnTx) Query(sql string, args ...interface{}) (*ConnRows, error) {
	rows, err := p.tx.Query(p.ctx, sql, args...)
	return newDatabaseRows(rows), err
}

func (p *ConnTx) QueryRow(sql string, args ...interface{}) *ConnRow {
	row := p.tx.QueryRow(p.ctx, sql, args...)
	return newDatabaseRow(row)
}
