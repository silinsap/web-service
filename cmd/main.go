package main

import (
	"log/slog"
	"os"
	"web-service/internal/config"
	"web-service/internal/httpserver"
	"web-service/internal/storage/sqlite"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.NewConfig()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

	log.Info("initializing server", slog.String("address", cfg.HTTPServer.Address)) // Помимо сообщения выведем параметр с адресом

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", slog.String("err", err.Error()))
	}

	storage.SaveLink("abc123", "https://example.com")

	router := gin.Default()

	router.GET("/links", httpserver.GetLinks(log, storage))
	router.POST("/links", httpserver.PostLinks(log, storage))
	router.GET("/links/:short_code", httpserver.GetLinkByShortCode(log, storage))
	router.DELETE("/links/:short_code", httpserver.DeleteLinkByShortCode(log, storage))
	router.GET("/links/:short_code/stats", httpserver.GetStatsByShortCode(log, storage))

	router.Run(cfg.HTTPServer.Address)
}
