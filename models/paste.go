package models

import (
	// stdlib
	"time"

	// other
	"github.com/jmoiron/sqlx"
)

const (
	PASTE_KEEP_FOR_MINUTES = 1
	PASTE_KEEP_FOR_HOURS   = 2
	PASTE_KEEP_FOR_DAYS    = 3
	PASTE_KEEP_FOR_MONTHS  = 4
)

var (
	PASTE_KEEPS_CORELLATION = map[string]int{
		"M": PASTE_KEEP_FOR_MINUTES,
		"h": PASTE_KEEP_FOR_HOURS,
		"d": PASTE_KEEP_FOR_DAYS,
		"m": PASTE_KEEP_FOR_MONTHS,
	}
)

type Paste struct {
	ID              int        `db:"id"`
	Title           string     `db:"title"`
	Data            string     `db:"data"`
	CreatedAt       *time.Time `db:"created_at"`
	KeepFor         int        `db:"keep_for"`
	KeepForUnitType int        `db:"keep_for_unit_type"`
}

func (p *Paste) GetByID(db *sqlx.DB) error {
	err := db.Get(p, db.Rebind("SELECT * FROM `pastes` WHERE id=?"), p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Paste) GetPagedPastes(db *sqlx.DB, page int) ([]Paste, error) {
	var pastes []Paste

	// Pagination - 30 pastes on page.
	var startPagination = 0
	if page > 1 {
		startPagination = (page - 1) * 30
	}

	err := db.Select(&pastes, db.Rebind("SELECT * FROM `pastes` ORDER BY id DESC LIMIT 30 OFFSET ?"), startPagination)
	if err != nil {
		return nil, err
	}

	return pastes, nil

}

func (p *Paste) Save(db *sqlx.DB) (int64, error) {
	result, err := db.NamedExec("INSERT INTO `pastes` (title, data, created_at, keep_for, keep_for_unit_type) VALUES (:title, :data, :created_at, :keep_for, :keep_for_unit_type)", p)
	if err != nil {
		return 0, err
	}

	ID, err1 := result.LastInsertId()
	if err1 != nil {
		return 0, err
	}

	return ID, nil
}
