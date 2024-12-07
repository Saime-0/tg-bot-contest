package updatesController

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func InTransaction(db *sqlx.DB, f func(tx *sqlx.Tx) error) (err error) {
	var tx *sqlx.Tx
	if tx, err = db.BeginTxx(context.Background(), nil); err != nil {
		return err
	}

	defer func() { // если во время выполнения действий случилась паника, то откатываем транзакцию
		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				slog.Error("tx.Rollback: " + err.Error())
			}
		}
	}()

	if err = f(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			slog.Error("tx.Rollback: " + err.Error())
		}
		return err
	}

	return tx.Commit()
}
