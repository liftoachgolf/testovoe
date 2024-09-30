package handler

import (
	geniusService "musPlayer/internal/serviceGenius"
	"musPlayer/internal/servicePostgres"
	"net/http"

	"github.com/gorilla/mux"

	_ "musPlayer/docs" // путь к сгенерированной документации

	httpSwagger "github.com/swaggo/http-swagger" // импортируем swagger
)

type Handler struct {
	services      *servicePostgres.Service
	serviceGenius *geniusService.GeniusService
}

func NewHandler(services *servicePostgres.Service, serviceGenius *geniusService.GeniusService) *Handler {
	return &Handler{
		services:      services,
		serviceGenius: serviceGenius,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := h.InitRoutes()
	router.ServeHTTP(w, r)
}

// @Summary Инициализация маршрутов
// @Description Инициализирует маршруты для обработчиков
// @Tags routes
// @Router /api [get]
func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	{
		songs := api.PathPrefix("/songs").Subrouter()
		{
			songs.HandleFunc("/", h.addSong).Methods(http.MethodPost)
			songs.HandleFunc("/search", h.searchSong).Methods(http.MethodPost)
			songs.HandleFunc("/filter", h.getFilteredSongs).Methods(http.MethodPost)
			songs.HandleFunc("/text", h.getTextWithPagination).Methods(http.MethodGet)
			songs.HandleFunc("/{id:[0-9]+}", h.updateSong).Methods(http.MethodPut)
			songs.HandleFunc("/{id:[0-9]+}", h.deleteSong).Methods(http.MethodDelete)
		}
		router.HandleFunc("/callback", h.callbackHandler).Methods(http.MethodGet)
		router.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
			h.serviceGenius.RedirectUser(w, r)
		}).Methods(http.MethodGet)
	}

	// Добавляем маршрут для Swagger-документации
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return router
}
