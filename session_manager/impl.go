package session_manager

import (
	"github.com/jmoiron/sqlx"
)

type manager struct {
	db *sqlx.DB
}

func NewManager(db *sqlx.DB) Manager {
	return &manager{db: db}
}

func (s *manager) CreateController() Controller {
	return newController(s.db)
}

type controller struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func newController(db *sqlx.DB) Controller {
	return &controller{db: db}
}

func (s *controller) Start() (err error) {
	if s.tx != nil {
		return ErrTransactionAlreadyStarted
	}

	s.tx, err = s.db.Beginx()
	return err
}

func (s *controller) Rollback() error {
	if s.tx == nil {
		return nil
	}

	if err := s.tx.Rollback(); err != nil {
		return err
	}

	s.tx = nil
	return nil
}

func (s *controller) Commit() error {
	if s.tx == nil {
		return nil
	}

	if err := s.tx.Commit(); err != nil {
		return err
	}

	s.tx = nil
	return nil
}

func (s *controller) Executor() Executor {
	if s.tx != nil {
		return s.tx
	}

	return s.db
}
