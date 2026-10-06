package gen_sql_executor

import "errors"

type errNoDataKey struct{}

var (
	ErrNoData = errors.New("no data found")
)
