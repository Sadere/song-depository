package app

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sadere/song-depository/internal/domain"
	"github.com/Sadere/song-depository/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ListSongs godoc
//
//	@Summary		List songs
//	@Description	list songs based on filter and page
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			message	body	domain.ListSongsRequest	true	"List songs request"
//	@Success		200	{array}		model.Song
//	@Success		204
//	@Failure		400	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/list-songs [post]
func (s *Server) ListSongs(c *gin.Context) {
	var request domain.ListSongsRequest

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.log.Debug("list song request: ", request)

	songs, err := s.songService.List(request.Filter, request.Page)

	if errors.Is(err, domain.ErrNoSongs) {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"songs": songs})
}

// AddSong godoc
//
//	@Summary		Add song
//	@Description	Add new song to depository
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			message	body		domain.AddSongRequest	true	"Add new song request"
//	@Success		201
//	@Failure		400	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Failure		503	{object}	ErrorResponse
//	@Router			/song [post]
func (s *Server) AddSong(c *gin.Context) {
	var request domain.AddSongRequest

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.log.Debug("add song request: ", request)

	validate := validator.New()

	err := validate.Struct(request)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("validation errors: %s", errors)})
		return
	}

	song := model.Song{
		Name:  request.Name,
		Group: request.Group,
	}

	err = s.songService.Add(&song)

	if errors.Is(err, domain.ErrSongDetail) {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// ModifySong godoc
//
//	@Summary		Edit song info
//	@Description	Edit any song info
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			song_id	path		int		true	"Song ID"
//	@Param			message	body		domain.UpdateSongRequest	true	"Edit song request"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/song/{song_id} [put]
func (s *Server) ModifySong(c *gin.Context) {
	var request domain.UpdateSongRequest

	i := c.Param("id")

	songID, err := strconv.Atoi(i)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.log.Debug("update song request: ", request)

	validate := validator.New()

	err = validate.Struct(request)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("validation errors: %s", errors)})
		return
	}

	// Modify song
	err = s.songService.Modify(uint64(songID), request)

	if errors.Is(err, domain.ErrSongNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		s.log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}

	c.Status(http.StatusOK)
}

// GetSongText godoc
//
//	@Summary		Get song text
//	@Description	Gets specific verse of requested song's text
//	@Tags			songs
//	@Produce		json
//	@Param			id			query		int		true	"Song ID"
//	@Param			verse		query		int		false	"Number of verse to return (starting at 0)"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/song-text/		[get]
func (s *Server) GetSongText(c *gin.Context) {
	i := c.Query("id")
	v := c.DefaultQuery("verse", "0")

	// Validate input data
	if len(i) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "please provide song ID"})
		return
	}

	songID, err := strconv.Atoi(i)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verse, err := strconv.Atoi(v)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch song text
	text, err := s.songService.Song(uint64(songID), verse)

	// Verse not found
	if errors.Is(err, domain.ErrVerseNotFound) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Other errors
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return JSON result
	c.JSON(http.StatusOK, gin.H{
		"text": text,
	})
}

// DeleteSong godoc
//
//	@Summary		Delete song
//	@Description	Deletes song with song_id from depository
//	@Tags			songs
//	@Produce		json
//	@Param			song_id	path		int		true	"Song ID"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/song/{song_id} [delete]
func (s *Server) DeleteSong(c *gin.Context) {
	i := c.Param("id")

	songID, err := strconv.Atoi(i)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.log.Debug("request to delete song id: ", songID)

	err = s.songService.Remove(uint64(songID))

	if errors.Is(err, domain.ErrSongNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}

	c.Status(http.StatusOK)
}
