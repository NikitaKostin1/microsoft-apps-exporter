package database

import (
	"database/sql"
	"fmt"
)

// withTransaction executes a function within a database transaction.
// It handles rollback and commit logic automatically.
func (db *Database) withTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := db.Connection.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer handleTxRollback(tx, &err)

	if err = fn(tx); err != nil {
		return err
	}

	return nil
}

// handleTxRollback handles the transaction rollback and commit logic.
func handleTxRollback(tx *sql.Tx, err *error) {
	if p := recover(); p != nil {
		*err = tx.Rollback()
		panic(p) // Re-throw panic after rollback
	} else if *err != nil {
		*err = tx.Rollback()
	} else {
		*err = tx.Commit()
	}
}
