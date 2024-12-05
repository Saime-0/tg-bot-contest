package updatesController

import (
	"context"
	"log/slog"
	"sync"

	"github.com/jmoiron/sqlx"
)

func InLazyTransaction(db *sqlx.DB, f func(tx func() *sqlx.Tx) error) (err error) {
	o := sync.Once{}
	mu := sync.Mutex{}
	var tx *sqlx.Tx
	getTx := func() *sqlx.Tx {
		mu.Lock()
		defer mu.Unlock()
		o.Do(func() {
			if tx, err = db.BeginTxx(context.Background(), nil); err != nil {
				panic(err)
			}
		})
		return tx
	}

	defer func() { // если во время выполнения действий случилась паника, то откатываем транзакцию
		if r := recover(); r != nil {
			if err := getTx().Rollback(); err != nil {
				slog.Error("tx.Rollback: " + err.Error())
			}
		}
	}()

	if err = f(getTx); err != nil {
		mu.Lock()
		if tx == nil {
			mu.Unlock()
			return err
		}
		mu.Unlock()
		if err := getTx().Rollback(); err != nil {
			slog.Error("tx.Rollback: " + err.Error())
		}
		slog.Info("tx.Rollback ended")
		return err
	}
	mu.Lock()
	if tx == nil {
		mu.Unlock()
		return nil
	}
	mu.Unlock()
	return getTx().Commit()
}
