package handler

import (
	"context"
	"encoding/json"
	"musPlayer/models"
	postgresrepo "musPlayer/pkg/postgresRepo"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type SongRequest struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

func (h *Handler) addSong(w http.ResponseWriter, r *http.Request) {
	var songRequest SongRequest

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&songRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if songRequest.Title == "" {
		http.Error(w, "Missing song title", http.StatusBadRequest)
		return
	}

	// Вызов метода для поиска песни по названию и артисту
	song, err := h.serviceGenius.SearchSong(songRequest.Title, songRequest.Artist)
	if err != nil {
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

func (h *Handler) searchSong(w http.ResponseWriter, r *http.Request) {
	var songRequest SongRequest

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&songRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if songRequest.Title == "" {
		http.Error(w, "Missing song title", http.StatusBadRequest)
		return
	}

	// Вызов метода для поиска песни по названию и артисту
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

type FilterParams struct {
	Filter string `json:"filter"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func (h *Handler) getFilteredSongs(w http.ResponseWriter, r *http.Request) {
	var params FilterParams

	// Декодируем JSON-запрос в структуру
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Получаем контекст и вызываем метод сервиса
	ctx := r.Context()
	songs, err := h.services.SongService.GetSongs(ctx, params.Filter, params.Limit, params.Offset)
	if err != nil {
		http.Error(w, "Failed to retrieve songs", http.StatusInternalServerError)
		return
	}

	// Успешно возвращаем список песен
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(songs); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

type GetTextWithPaginationParams struct {
	Id       int `json:"id"`        // ID песни
	PageSize int `json:"page_size"` // Количество символов на странице
	Page     int `json:"page"`      // Номер страницы
}

func (h *Handler) getTextWithPagination(w http.ResponseWriter, r *http.Request) {
	var params GetTextWithPaginationParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if params.PageSize == 0 {
		params.PageSize = 100 // значение по умолчанию
	}

	if params.Page == 0 {
		params.Page = 1
	}

	// Получаем контекст и вызываем метод сервиса
	ctx := r.Context()
	text, err := h.services.GetSongText(ctx, params.Id, params.PageSize, params.Page)
	if err != nil {
		http.Error(w, "Failed to get text", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"text": text}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteSong(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID песни из параметров
	vars := mux.Vars(r)
	id := vars["id"]
	idd, _ := strconv.Atoi(id)
	// Вызов сервиса для удаления песни
	err := h.services.DeleteSong(r.Context(), int64(idd))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Успешное удаление
}

// Обработчик для обновления данных песни
func (h *Handler) updateSong(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID песни из параметров
	vars := mux.Vars(r)
	id := vars["id"]

	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	song.ID = id // Устанавливаем ID для обновления

	// Вызов сервиса для обновления песни
	if err := h.services.UpdateSong(r.Context(), song); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // Успешное обновление
}
