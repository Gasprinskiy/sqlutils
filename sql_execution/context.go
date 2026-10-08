package sql_execution

import "context"

func extractErrorNoData(ctx context.Context) error {
	err, exists := ctx.Value(errNoDataKey{}).(error)
	if !exists {
		err = ErrNoData
	}

	return err
}

func SetCustomErrNoDataCtx(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errNoDataKey{}, err)
}
