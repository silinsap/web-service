package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web-service/internal/config"
	"web-service/internal/httpserver"
	"web-service/internal/memory"
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

	memory := memory.NewMemoryStorage()

	router := gin.Default()

	// Создание сервера
	srv := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	router.GET("/links", httpserver.GetLinks(log, storage))
	router.POST("/links", httpserver.PostLinks(log, storage))
	router.GET("/links/:short_code", httpserver.GetLinkByShortCode(log, storage, memory))
	router.DELETE("/links/:short_code", httpserver.DeleteLinkByShortCode(log, storage, memory))
	router.GET("/links/:short_code/stats", httpserver.GetStatsByShortCode(log, storage, memory))

	// Запуск сервера в отдельной горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Ошибка: %v", slog.String("err", err.Error()))
		}
	}()

	// Ожидание сигналов ОС для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown с таймаутом 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Ошибка: %v", slog.String("err", err.Error()))
	}

}
