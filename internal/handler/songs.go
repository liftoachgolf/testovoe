package handler

import (
	"context"
	"encoding/json"
	"musPlayer/internal/logger"
	"musPlayer/models"
	postgresrepo "musPlayer/pkg/postgresRepo"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// SongRequest представляет запрос на добавление новой песни.
type SongRequest struct {
	Title  string `json:"song"`
	Artist string `json:"group"`
}

// @Summary Добавить новую песню
// @Description Добавляет новую песню в базу данных
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body SongRequest true "Данные о песне"
// @Success 200 {object} models.Song
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/songs [post]
func (h *Handler) addSong(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Debugf("Incoming request to %s with method: %s", r.URL.Path, r.Method)

	var songRequest SongRequest
	if err := json.NewDecoder(r.Body).Decode(&songRequest); err != nil {
		logger.Logger.Errorf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if songRequest.Title == "" {
		http.Error(w, "Missing song title", http.StatusBadRequest)
		return
	}

	logger.Logger.Debugf("Received song request: %+v", songRequest)

	song, err := h.serviceGenius.SearchSong(songRequest.Title, songRequest.Artist)
	if err != nil {
		logger.Logger.Errorf("Error searching song: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if song == nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	h.services.AddSong(context.Background(), postgresrepo.AddSongParams{
		SongId:      song.ID,
		GroupName:   song.GroupName,
		SongName:    song.SongName,
		Text:        song.Text,
		Link:        song.Link,
		ReleaseDate: song.ReleaseDate,
	})

	logger.Logger.Debugf("Song added successfully: %s by %s", songRequest.Title, songRequest.Artist)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// @Summary Найти песню
// @Description Ищет песню по заголовку и исполнителю
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body SongRequest true "Данные о песне"
// @Success 200 {object} models.Song
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "Song not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/songs/search [post]
func (h *Handler) searchSong(w http.ResponseWriter, r *http.Request) {
	var songRequest SongRequest

	if err := json.NewDecoder(r.Body).Decode(&songRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if songRequest.Title == "" {
		http.Error(w, "Missing song title", http.StatusBadRequest)
		return
	}
	song, err := h.serviceGenius.SearchSong(songRequest.Title, songRequest.Artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if song == nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	// Возвращаем найденную песню
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// FilterParams представляет параметры фильтрации для получения песен.
type FilterParams struct {
	Filter string `json:"filter"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// @Summary Получить отфильтрованные песни
// @Description Получает песни, основываясь на заданных фильтрах
// @Tags songs
// @Accept  json
// @Produce  json
// @Param filter body FilterParams true "Параметры фильтрации"
// @Success 200 {array} models.Song
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Failed to retrieve songs"
// @Router /api/songs/filter [post]
func (h *Handler) getFilteredSongs(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Debugf("Incoming request to %s with method: %s", r.URL.Path, r.Method)

	var params FilterParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.Logger.Errorf("Failed to decode request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	logger.Logger.Debugf("Received filter params: %+v", params)

	ctx := r.Context()
	songs, err := h.services.SongService.GetSongs(ctx, params.Filter, params.Limit, params.Offset)
	if err != nil {
		logger.Logger.Errorf("Failed to retrieve songs: %v", err)
		http.Error(w, "Failed to retrieve songs", http.StatusInternalServerError)
		return
	}

	logger.Logger.Debugf("Retrieved %d songs", len(songs))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(songs); err != nil {
		logger.Logger.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetTextWithPaginationParams представляет параметры для получения текста песни с пагинацией.
type GetTextWithPaginationParams struct {
	Id       int `json:"id"`
	PageSize int `json:"page_size"`
	Page     int `json:"page"`
}

// @Summary Получить текст песни с пагинацией
// @Description Получает текст песни по идентификатору с возможностью пагинации
// @Tags songs
// @Accept  json
// @Produce  json
// @Param params body GetTextWithPaginationParams true "Параметры запроса"
// @Success 200 {object} map[string]string "Текст песни"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Failed to get text"
// @Router /api/songs/text [post]
func (h *Handler) getTextWithPagination(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Debugf("Incoming request to %s with method: %s", r.URL.Path, r.Method)

	var params GetTextWithPaginationParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.Logger.Errorf("Failed to decode request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if params.PageSize == 0 {
		params.PageSize = 100 // значение по умолчанию
	}

	if params.Page == 0 {
		params.Page = 1
	}

	logger.Logger.Debugf("Retrieving text for song ID: %d with pagination: %+v", params.Id, params)

	ctx := r.Context()
	text, err := h.services.GetSongText(ctx, params.Id, params.PageSize, params.Page)
	if err != nil {
		logger.Logger.Errorf("Failed to get text: %v", err)
		http.Error(w, "Failed to get text", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"text": text}); err != nil {
		logger.Logger.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.Logger.Debugf("Successfully retrieved text for song ID: %d", params.Id)
}

// @Summary Удалить песню
// @Description Удаляет песню по идентификатору
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Идентификатор песни"
// @Success 204 {string} string "No Content"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/songs/{id} [delete]
func (h *Handler) deleteSong(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Debugf("Incoming request to %s with method: %s", r.URL.Path, r.Method)

	vars := mux.Vars(r)
	id := vars["id"]
	idd, err := strconv.Atoi(id)
	if err != nil {
		logger.Logger.Errorf("Invalid song ID: %s, error: %v", id, err)
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	logger.Logger.Debugf("Deleting song with ID: %d", idd)

	err = h.services.DeleteSong(r.Context(), int64(idd))
	if err != nil {
		logger.Logger.Errorf("Failed to delete song: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Logger.Debugf("Song with ID: %d deleted successfully", idd)
}

// GetSongUpdateParams представляет параметры для обновления данных о песне.
type GetSongUpdateParams struct {
	GroupName   string `json:"group"`
	SongName    string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
}

// @Summary Обновить данные о песне
// @Description Обновляет данные о песне по идентификатору
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Идентификатор песни"
// @Param params body GetSongUpdateParams true "Данные для обновления"
// @Success 200 {string} string "Song updated successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/songs/{id} [put]
func (h *Handler) updateSong(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Debugf("Incoming request to %s with method: %s", r.URL.Path, r.Method)

	vars := mux.Vars(r)
	id := vars["id"]
	idd, err := strconv.Atoi(id)
	if err != nil {
		logger.Logger.Errorf("Invalid song ID: %s, error: %v", id, err)
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	var params GetSongUpdateParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.Logger.Errorf("Failed to decode request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	logger.Logger.Debugf("Updating song with ID: %d, params: %+v", idd, params)

	if err := h.services.UpdateSong(r.Context(), models.SongUpdateParams{
		GroupName:   params.GroupName,
		SongName:    params.SongName,
		ReleaseDate: params.ReleaseDate,
		Text:        params.Text,
		ID:          idd,
	}); err != nil {
		logger.Logger.Errorf("Failed to update song: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Logger.Debugf("Song with ID: %d updated successfully", idd)
}
