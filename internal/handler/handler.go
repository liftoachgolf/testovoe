package handler

import (
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
			songs.HandleFunc("/", h.addSong).Methods(http.MethodGet)                 // Путь для добавления песни
			songs.HandleFunc("/search", h.searchSong).Methods(http.MethodPost)       // Изменено на POST для поиска песни
			songs.HandleFunc("/filter", h.getFilteredSongs).Methods(http.MethodPost) // получения песен с фильтрацией
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
