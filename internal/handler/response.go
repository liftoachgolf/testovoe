package handler

import (
	"encoding/json"
	"musPlayer/internal/logger"
	"net/http"
)

func sendSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, "failed to encode response")
	}
}

// Обработчик для ошибок
func newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Функция для обработки ошибок
func (h *Handler) handleError(w http.ResponseWriter, err error, statusCode int, message string) {
	// Записываем ошибку в лог
	logger.Logger.WithError(err).Error(message)
	newErrorResponse(w, statusCode, message)
}
