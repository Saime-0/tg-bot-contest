package autotx

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func Commit(tx *sqlx.Tx, err error, r any) {
	if r != nil {
		if err = tx.Rollback(); err != nil {
			slog.Error("tx.Rollback: " + err.Error())
		}
		return
	}

	if err != nil {
		if err := tx.Rollback(); err != nil {
			slog.Error("tx.Rollback: " + err.Error())
		}
		return
	}

	if err = tx.Commit(); err != nil {
		slog.Error("autotx: Commit: tx.Commit: " + err.Error())
	}
}
