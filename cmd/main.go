package main

import (
	"log"
	musplayer "musPlayer"
	"musPlayer/internal/config"
	"musPlayer/internal/handler"
	"musPlayer/internal/logger"
	servicegenius "musPlayer/internal/serviceGenius"
	"musPlayer/internal/servicePostgres"
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
	geniusSrv := servicegenius.NewGeniusService(cfg.GeniusConfig.ID, cfg.GeniusConfig.Secret, cfg.GeniusConfig.RedirectURI)
	handler := handler.NewHandler(dbSrv, geniusSrv)

	srv := new(musplayer.Server)
	if err := srv.Run(cfg.App.Port, handler); err != nil {
		logrus.Panic("error while running server")
	}
}
