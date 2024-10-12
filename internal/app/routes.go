package app

import (
	"github.com/Sadere/song-depository/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) setupRoutes() (*gin.Engine, error) {
	r := gin.New()

	// Attach logger
	r.Use(middleware.Logger(s.log))

	// Default gin panic recovery middleware
	r.Use(gin.Recovery())

	// Set response content-type
	r.Use(middleware.JSON())

	// Routes
	jsonRoutes := r.Group("")

	jsonRoutes.Use(middleware.CheckJSON())
	{
		r.POST("/list-songs", s.ListSongs)
		r.POST("/song", s.AddSong)
		r.PUT("/song/:id", s.ModifySong)
	}

	r.GET("/song-text", s.GetSongText)
	r.DELETE("/song/:id", s.DeleteSong)

	// Swagger routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r, nil
}
