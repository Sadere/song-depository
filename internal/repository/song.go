package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Sadere/song-depository/internal/domain"
	"github.com/Sadere/song-depository/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	SongsPerPage = 10
)

// Song storage repository
type SongRepository interface {
	Create(ctx context.Context, song *model.Song) error
	GetById(ctx context.Context, songID uint64) (*model.Song, error)
	ListFiltered(ctx context.Context, filter domain.SongFilter, page uint) (model.Songs, error)
	GetSongText(ctx context.Context, songID uint64) (string, error)
	Update(ctx context.Context, songID uint64, req domain.UpdateSongRequest) error
	Delete(ctx context.Context, songID uint64) error
}

type PgSongRepository struct {
	db *sqlx.DB
}

func NewPgSongRepository(db *sqlx.DB) *PgSongRepository {
	return &PgSongRepository{
		db: db,
	}
}

// Creates new song in DB
func (r *PgSongRepository) Create(ctx context.Context, song *model.Song) error {
	sb := sq.StatementBuilder.
		Insert("songs").
		Columns("song_name", "song_group", "song_text", "release_date", "link").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db)

	sb = sb.Values(
		song.Name,
		song.Group,
		song.Text,
		song.ReleaseDate,
		song.Link,
	)

	_, err := sb.ExecContext(ctx)

	return err
}

func (r *PgSongRepository) GetById(ctx context.Context, songID uint64) (*model.Song, error) {
	var song model.Song

	sb := sq.Select(
		"id",
		"created_at",
		"updated_at",
		"song_name",
		"song_group",
		"song_text",
		"release_date",
		"link",
	).
		From("songs").
		Where(sq.Eq{
			"id": songID,
		}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "repository.GetById")
	}

	err = r.db.QueryRowxContext(ctx, query, args...).StructScan(&song)

	if err != nil {
		return nil, errors.Wrap(err, "repository.GetById")
	}

	return &song, nil
}

// Fetches songs with pagination and filter
func (r *PgSongRepository) ListFiltered(ctx context.Context, filter domain.SongFilter, page uint) (model.Songs, error) {
	var songs model.Songs

	offset := uint64(page * SongsPerPage)

	sb := sq.Select(
		"id",
		"created_at",
		"updated_at",
		"song_name",
		"song_group",
		"song_text",
		"release_date",
		"link",
	).
		From("songs").
		OrderBy("id DESC").
		Limit(SongsPerPage).
		Offset(offset).
		PlaceholderFormat(sq.Dollar)

	// Filter
	if filter.Group != nil {
		sb = sb.Where(sq.Like{
			"song_group": *filter.Group,
		})
	}

	if filter.Name != nil {
		sb = sb.Where(sq.Like{
			"song_name": *filter.Name,
		})
	}

	if filter.Text != nil {
		sb = sb.Where(sq.Like{
			"song_text": "%" + *filter.Text + "%",
		})
	}

	if filter.ReleaseDate != nil {
		sb = sb.Where(sq.Eq{
			"release_date": *filter.ReleaseDate,
		})
	}

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "repository.ListFiltered")
	}

	// Fetch query
	err = r.db.SelectContext(ctx, &songs, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "repository.ListFiltered")
	}

	if len(songs) == 0 {
		return nil, errors.Wrap(domain.ErrNoSongs, "repository.ListFiltered")
	}

	return songs, nil
}

// Fetch song text from DB with provided song ID
func (r *PgSongRepository) GetSongText(ctx context.Context, songID uint64) (string, error) {
	var songText string

	sb := sq.Select(
		"song_text",
	).
		From("songs").
		Where(sq.Eq{
			"id": songID,
		}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sb.ToSql()
	if err != nil {
		return "", errors.Wrap(err, "repository.GetSongText")
	}

	err = r.db.QueryRowxContext(ctx, query, args...).Scan(&songText)

	if err != nil {
		return "", errors.Wrap(err, "repository.GetSongText")
	}

	return songText, nil
}

// Updates song in DB
func (r *PgSongRepository) Update(ctx context.Context, songID uint64, req domain.UpdateSongRequest) error {
	sb := sq.StatementBuilder.
		Update("songs").
		Set("updated_at", time.Now()).
		Where(sq.Eq{
			"id": songID,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db)
	
	if len(req.Group) > 0 {
		sb = sb.Set("song_group", req.Group)
	}

	if len(req.Name) > 0 {
		sb = sb.Set("song_name", req.Name)
	}

	if len(req.Text) > 0 {
		sb = sb.Set("song_text", req.Text)
	}

	if req.ReleaseDate != nil {
		sb = sb.Set("release_date", *req.ReleaseDate)
	}

	if len(req.Link) > 0 {
		sb = sb.Set("link", req.Link)
	}

	_, err := sb.ExecContext(ctx)

	return err
}

// Removes song from DB
func (r *PgSongRepository) Delete(ctx context.Context, songID uint64) error {
	sb := sq.StatementBuilder.
		Delete("songs").
		Where(sq.Eq{
			"id": songID,
		}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db)

	_, err := sb.ExecContext(ctx)

	return err
}
