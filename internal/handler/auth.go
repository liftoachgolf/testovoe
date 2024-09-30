package handler

import (
	"fmt"
	"net/http"
)

// @Summary Добавить новую песню
// @Description Добавляет новую песню в базу данных
// @Tags songs
// @Accept json
// @Produce json
// @Param song body SongRequest true "Данные о песне"
// @Success 200 {object} models.Song
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/songs [post]
func (h *Handler) callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		h.handleError(w, fmt.Errorf("code is missing"), http.StatusBadRequest, "Code is missing")
		return
	}

	// Получение токена доступа
	if err := h.serviceGenius.GetAccessToken(code); err != nil {
		h.handleError(w, err, http.StatusInternalServerError, "Failed to obtain access token")
		return
	}

	// Вывод токена на экран
	response := map[string]string{
		"access_token": h.serviceGenius.AccessToken,
	}
	sendSuccessResponse(w, http.StatusOK, response) // Используем sendSuccessResponse
}
