package models_dao

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

type TransactionDao struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) *TransactionDao {
	return &TransactionDao{
		db: db,
	}
}

func (t *TransactionDao) StartTransaction() (*sqlx.Tx, error) {
	return t.db.Beginx()
}

func (t *TransactionDao) ShutDown(tx *sqlx.Tx, err error) error {
	if rollErr := tx.Rollback(); err != nil {
		return errors.Join(err, rollErr)
	}
	return err
}

func (t *TransactionDao) Commit(tx *sqlx.Tx) error {
	if err := tx.Commit(); err != nil {
		return t.ShutDown(tx, err)
	}
	return nil
}
