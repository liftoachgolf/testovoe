package todo

import (
	"context"
	"musPlayer/internal/logger"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

// Run запускает HTTP-сервер на заданном порту с указанным обработчиком
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:              ":" + port,
		MaxHeaderBytes:    1 << 20, // 1MB
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           handler,
	}

	// Логируем, что сервер запущен
	logger.Logger.Infof("Server is running on port %s", port)

	// Запускаем сервер и возвращаем ошибку, если она произошла
	return s.httpServer.ListenAndServe()
}

// Shutdown корректно завершает работу сервера
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
