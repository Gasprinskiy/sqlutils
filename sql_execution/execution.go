package sql_execution

import (
	"context"
	"database/sql"

	"github.com/Gasprinskiy/sqlutils/session_manager"
)

func Exec(
	ctx context.Context,
	e session_manager.Executor,
	sqlQuery string,
	args ...any,
) error {
	_, err := e.ExecContext(ctx, sqlQuery, args...)
	return err
}

func ExecNamed(
	ctx context.Context,
	e session_manager.Executor,
	sqlQuery string,
	data any,
) error {
	_, err := e.NamedExecContext(ctx, sqlQuery, data)
	return err
}

func ExecNamedReturningFields[F comparable](
	ctx context.Context,
	e session_manager.Executor,
	sqlQuery string,
	data any,
) (F, error) {
	var fields F

	stmt, err := e.PrepareNamed(sqlQuery)
	if err != nil {
		return fields, err
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &fields, data)
	if err != nil {
		return fields, err
	}

	return fields, nil
}

func Get[T any](
	ctx context.Context,
	e session_manager.Executor,
	sqlQuery string,
	args ...any,
) (T, error) {
	var data T

	err := e.GetContext(ctx, &data, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return data, extractErrorNoData(ctx)
		}

		return data, err
	}

	return data, nil
}

func Select[T any](
	ctx context.Context,
	e session_manager.Executor,
	sqlQuery string,
	args ...any,
) ([]T, error) {
	var data []T

	err := e.SelectContext(ctx, &data, sqlQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, extractErrorNoData(ctx)
		}

		return nil, err
	}

	if len(data) == 0 {
		return nil, extractErrorNoData(ctx)
	}

	return data, nil
}

func SelectNamed[T any](
	ctx context.Context,
	e session_manager.Executor,
	sqlQuery string,
	param any,
) ([]T, error) {
	var data []T

	stmt, err := e.PrepareNamedContext(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.SelectContext(ctx, &data, param)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, extractErrorNoData(ctx)
		}

		return nil, err
	}

	if len(data) == 0 {
		return nil, extractErrorNoData(ctx)
	}

	return data, nil
}
