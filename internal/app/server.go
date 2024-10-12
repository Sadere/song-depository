package app

import (
	"net/http"

	"github.com/Sadere/song-depository/internal/config"
	"github.com/Sadere/song-depository/internal/database"
	"github.com/Sadere/song-depository/internal/repository"
	"github.com/Sadere/song-depository/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Server struct {
	config      *config.Config
	songService service.ISongService
	log         *zap.SugaredLogger
	db          *sqlx.DB
}

func NewServer(
	cfg *config.Config,
	log *zap.SugaredLogger,
	db *sqlx.DB,
) *Server {
	// Init repo
	songRepo := repository.NewPgSongRepository(db)

	// Init service
	songService := service.NewSongService(cfg, songRepo, log)

	return &Server{
		config:      cfg,
		songService: songService,
		log:         log,
		db:          db,
	}
}

func (s *Server) Start() error {
	// Run migrations
	err := database.MigrateUp(s.config.PostgresDSN)
	if err != nil {
		return errors.Wrap(err, "database.MigrateUp")
	}

	// Setup routes
	r, err := s.setupRoutes()
	if err != nil {
		return errors.Wrap(err, "setupRoutes")
	}

	srv := &http.Server{
		Addr:    s.config.Address,
		Handler: r,
	}

	// Run server in background
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatalf("listen: %s\n", err)
		}
	}()

	return nil
}
