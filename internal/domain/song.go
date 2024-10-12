package domain

import (
	"errors"
	"time"
)

var (
	ErrNoSongs       = errors.New("no songs found")
	ErrSongDetail    = errors.New("failed to retrieve song detail")
	ErrSongNotFound  = errors.New("song with provided ID not found")
	ErrVerseNotFound = errors.New("requested verse doesn't exist in this song")
)

type ErrorResponse struct {
	Error string `json:"error" example:"Fatal error"`
} // @name ErrorResponse

type SongFilter struct {
	Name        *string    `json:"name"`
	Group       *string    `json:"group"`
	Text        *string    `json:"text"`
	ReleaseDate *time.Time `json:"release_date" example:"2024-10-29T15:04:05.000Z"`
}

type ListSongsRequest struct {
	Filter SongFilter `json:"filter"`
	Page   uint       `json:"page"`
}

type AddSongRequest struct {
	Name  string `json:"song" validate:"required,min=1"`
	Group string `json:"group" validate:"required,min=1"`
}

type UpdateSongRequest struct {
	Name        string     `json:"song"`
	Group       string     `json:"group"`
	Text        string     `json:"text"`
	ReleaseDate *time.Time `json:"release_date" example:"2024-10-29T15:04:05.000Z"`
	Link        string     `json:"link" validate:"http_url"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
