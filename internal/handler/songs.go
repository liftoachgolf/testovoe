package handler

import (
	"context"
	"encoding/json"
	postgresrepo "musPlayer/pkg/postgresRepo"
	"net/http"
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
