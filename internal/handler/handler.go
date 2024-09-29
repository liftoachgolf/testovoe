package handler

import (
	"encoding/json"
	"fmt"
	geniusService "musPlayer/internal/serviceGenius" // Импортируем пакет для работы с Genius API
	"musPlayer/internal/servicePostgres"
	"net/http"

	"github.com/gorilla/mux"
)

// Объявляем структуру Handler
type Handler struct {
	services      *servicePostgres.Service
	serviceGenius *geniusService.GeniusService
}

// Новый конструктор для Handler
func NewHandler(services *servicePostgres.Service, serviceGenius *geniusService.GeniusService) *Handler {
	return &Handler{
		services:      services,
		serviceGenius: serviceGenius,
	}
}

// Реализуем метод ServeHTTP для интерфейса http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := h.InitRoutes() // Инициализируем маршруты
	router.ServeHTTP(w, r)   // Передаем запрос в маршрутизатор
}

// Инициализация маршрутов
func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	{
		songs := api.PathPrefix("/songs").Subrouter()
		{
			songs.HandleFunc("/", h.addSong).Methods(http.MethodGet)          // Путь для добавления песни
			songs.HandleFunc("/search", h.searchSong).Methods(http.MethodGet) // Путь для поиска песни
		}

		// Добавление маршрута для обработки callback
		router.HandleFunc("/callback", h.callbackHandler).Methods(http.MethodGet)

		// Добавление маршрута для перенаправления на авторизацию
		router.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
			h.serviceGenius.RedirectUser(w, r) // Редирект пользователя для авторизации в Genius
		}).Methods(http.MethodGet)
	}

	return router
}

// Обработчик для callback
func (h *Handler) callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is missing", http.StatusBadRequest)
		return
	}

	// Получение токена доступа
	if err := h.serviceGenius.GetAccessToken(code); err != nil {
		http.Error(w, "Failed to obtain access token", http.StatusInternalServerError)
		return
	}

	// Вывод токена на экран
	response := map[string]string{
		"access_token": h.serviceGenius.AccessToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Fprintf(w, "Access token: %s", h.serviceGenius.AccessToken) // Также выводим на экран
}

// Обработчик для поиска песни
func (h *Handler) searchSong(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	artist := r.URL.Query().Get("artist")

	if title == "" && artist == "" {
		http.Error(w, "Missing song title and artist", http.StatusBadRequest)
		return
	}

	// Вызов метода для поиска песни по названию и артисту
	song, err := h.serviceGenius.SearchSong(title, artist)
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
