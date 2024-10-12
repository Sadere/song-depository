package model

import "time"

type Song struct {
	ID          uint64    `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Name        string    `db:"song_name"`
	Group       string    `db:"song_group"`
	Text        string    `db:"song_text"`
	ReleaseDate time.Time `db:"release_date"`
	Link        string    `db:"link"`
}

type Songs []*Song
