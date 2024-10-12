package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Sadere/song-depository/internal/config"
	"github.com/Sadere/song-depository/internal/domain"
	"github.com/Sadere/song-depository/internal/model"
	"github.com/Sadere/song-depository/internal/repository"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ISongService interface {
	Add(song *model.Song) error
	List(filter domain.SongFilter, page uint) (model.Songs, error)
	Song(songID uint64, verse int) (string, error)
	Modify(songID uint64, req domain.UpdateSongRequest) error
	Remove(songID uint64) error
}

type SongService struct {
	config   *config.Config
	songRepo repository.SongRepository
	log      *zap.SugaredLogger
}

func NewSongService(config *config.Config, songRepo repository.SongRepository, log *zap.SugaredLogger) *SongService {
	return &SongService{
		config:   config,
		songRepo: songRepo,
		log:      log,
	}
}

func (s *SongService) Add(song *model.Song) error {
	// Request music info endpoint
	var songDetail domain.SongDetail

	params := url.Values{
		"group": {song.Group},
		"song":  {song.Name},
	}

	infoEndPoint := fmt.Sprintf("%s/info?%s", s.config.MusicInfoAddress, params.Encode())

	s.log.Debug("music info endpoint request: ", infoEndPoint)

	response, err := resty.New().R().
		SetResult(&songDetail).
		Get(infoEndPoint)

	if err != nil {
		return errors.Wrap(err, "songDetail")
	}

	s.log.Debug("music info endpoint response: ", response, " body: ", songDetail)

	if response.StatusCode() != http.StatusOK {
		return domain.ErrSongDetail
	}

	// Process song detail response
	releaseDate, err := time.Parse("02.01.2006", songDetail.ReleaseDate)
	if err != nil {
		return errors.Wrap(err, "time.Parse")
	}
	song.ReleaseDate = releaseDate
	song.Text = songDetail.Text
	song.Link = songDetail.Link

	// Save song to storage
	err = s.songRepo.Create(context.TODO(), song)
	if err != nil {
		return errors.Wrap(err, "songRepo.Create")
	}

	return nil
}

func (s *SongService) List(filter domain.SongFilter, page uint) (model.Songs, error) {
	return s.songRepo.ListFiltered(context.TODO(), filter, page)
}

func (s *SongService) Song(songID uint64, verse int) (string, error) {
	songText, err := s.songRepo.GetSongText(context.TODO(), songID)
	if err != nil {
		return "", errors.Wrap(err, "songRepo.GetSongText")
	}

	// Get required verse
	verses := strings.Split(songText, "\n\n")

	if verse > len(verses) {
		return "", domain.ErrVerseNotFound
	}

	return verses[verse], nil
}

func (s *SongService) Modify(songID uint64, req domain.UpdateSongRequest) error {
	// Check if song exists
	_, err := s.songRepo.GetById(context.TODO(), songID)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrSongNotFound
	}

	if err != nil {
		return err
	}

	return s.songRepo.Update(context.TODO(), songID, req)
}

func (s *SongService) Remove(songID uint64) error {
	// Check if song exists
	_, err := s.songRepo.GetById(context.TODO(), songID)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrSongNotFound
	}

	if err != nil {
		return err
	}

	return s.songRepo.Delete(context.TODO(), songID)
}
