package handler

import (
	"net/http"
)

func (h *Handler) addSong(w http.ResponseWriter, r *http.Request) {
	// var req struct {
	// 	Group string `json:"group"`
	// 	Song  string `json:"song"`
	// }

	// // Парсим JSON запрос
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	newErrorResponse(w, http.StatusBadRequest, "Invalid request body")
	// 	return
	// }

	// // Ищем трек по названию и исполнителю через Spotify API
	// songs, err := h.serviceSpotify.SearchTrackByNameAndArtist(req.Song, req.Group)
	// if err != nil {
	// 	newErrorResponse(w, http.StatusInternalServerError, "Failed to search song")
	// 	return
	// }

	// if len(songs) == 0 {
	// 	newErrorResponse(w, http.StatusNotFound, "Song not found")
	// 	return
	// }

	// songs = []models.SpotySong{
	// 	{ID: "1", SongName: "Test Song", Artist: "Test Artist"},
	// }
	// song := songs[0]
	// var songForAdd = postgresrepo.AddSongParams{
	// 	GroupName:   song.Artist,
	// 	SongName:    song.SongName,
	// 	Text:        "xdddddddddddddd",
	// 	ReleaseDate: song.ReleaseDate,
	// 	Link:        song.TrackURL,
	// }

	// // Сохраняем трек в базе данных
	// _, err = h.services.AddSong(context.Background(), songForAdd)
	// if err != nil {
	// 	newErrorResponse(w, http.StatusInternalServerError, "Failed to save song to database")
	// 	return
	// }

	// // Возвращаем успешный ответ
	// sendSuccessResponse(w, http.StatusCreated, map[string]string{"message": "Song added successfully"})
}
