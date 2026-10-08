package sql_execution

import "errors"

type errNoDataKey struct{}

var (
	ErrNoData = errors.New("no data found")
)
