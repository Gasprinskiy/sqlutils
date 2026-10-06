package session_manager

import "context"

func SetExecutor(ctx context.Context, s Executor) context.Context {
	return context.WithValue(ctx, executorKey{}, s)
}

func MustGetExecutor(ctx context.Context) Executor {
	return ctx.Value(executorKey{}).(Executor)
}

func SetController(ctx context.Context, c Controller) context.Context {
	return context.WithValue(ctx, sessionControllerKey{}, c)
}

func MustGetController(ctx context.Context) Controller {
	return ctx.Value(sessionControllerKey{}).(Controller)
}
