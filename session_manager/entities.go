package session_manager

import "errors"

var (
	ErrTransactionAlreadyStarted = errors.New("transaction already started")
)

type executorKey struct{}
type sessionControllerKey struct{}
