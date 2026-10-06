package session_manager

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Manager interface {
	CreateController() Controller
}

type Controller interface {
	Start() error
	Rollback() error
	Commit() error
	Executor() Executor
}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
