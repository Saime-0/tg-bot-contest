package ue

import (
	"database/sql"
	"errors"

	"tgBotContest/l10n"
)

type Err struct {
	text string
}

func (u *Err) Error() string {
	return u.text
}

func New(s string) *Err {
	return &Err{text: s}
}

func Sql(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return New(l10n.RequestedDataNotFound)
	}

	return err
}
