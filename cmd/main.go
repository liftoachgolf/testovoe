package main

import (
	"log"
	musplayer "musPlayer"
	"musPlayer/internal/config"
	"musPlayer/internal/handler"
	"musPlayer/internal/logger"
	servicegenius "musPlayer/internal/serviceGenius"
	"musPlayer/internal/servicePostgres"
	"musPlayer/internal/spotifyService"
	postgresrepo "musPlayer/pkg/postgresRepo"

	"github.com/sirupsen/logrus"
)

func main() {

	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatal("Error loading configuration: ", err)
	}

	// Инициализация логгера
	logger.InitLogger(cfg.Logging.Level)

	// Пример использования конфигурации
	logger.Logger.Infof("Starting application on port %s", cfg.App.Port)

	db, err := postgresrepo.NewPostgresDb(cfg.Database)
	if err != nil {
		logrus.Fatalf("error while connecting to db: %v", err)
	}

	dbRepo := postgresrepo.NewRepository(db)
	dbSrv := servicePostgres.NewServicePostgres(dbRepo)
	_ = spotifyService.NewSpotifyService(cfg.SpotifyApi.ID, cfg.SpotifyApi.Secret, cfg.SpotifyApi.RedirectURI)
	geniusSrv := servicegenius.NewGeniusService("FfBcmLXIRKkSUMIoNgjGZTQLIVt4EJhAuBiMCwrKa6oKPR0QRBjt6maOw5wQYB8e", "UB3CDuzus25Hf5sxF4qCaAaSti2YTsOUkLJTmxqqpsTPnugcWoQVh-tzVom38UO6Js8BuuG9dUNcPso32i1Xfw", "http://localhost:8000/callback")
	handler := handler.NewHandler(dbSrv, geniusSrv)

	srv := new(musplayer.Server)
	if err := srv.Run("8000", handler); err != nil {
		logrus.Panic("error while running server")
	}
}
