package servicePostgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"musPlayer/internal/logger"
	"musPlayer/models"
	postgresrepo "musPlayer/pkg/postgresRepo"
	"strings"
	"time"
)

type songService struct {
	repo postgresrepo.SongRepository
}

func NewSongService(repo postgresrepo.SongRepository) SongService {
	return &songService{
		repo: repo,
	}
}
func (s *songService) AddSong(ctx context.Context, song postgresrepo.AddSongParams) (int, error) {
	logger.Logger.Debugf("Adding song: %+v", song)
	startTime := time.Now()

	id, err := s.repo.AddSong(ctx, song)
	if err != nil {
		logger.Logger.Error("Error adding song: ", err)
		return 0, err
	}

	logger.Logger.Infof("Song added successfully with ID: %d, execution time: %s", id, time.Since(startTime))
	return id, nil
}

// Получение текста песни с обработкой
func paginateText(text string, pageSize int) []string {
	words := strings.Split(text, " ")
	var pages []string
	var page string

	for _, word := range words {
		if len(page)+len(word)+1 > pageSize {
			pages = append(pages, page)
			page = ""
		}
		if page != "" {
			page += " "
		}
		page += word
	}

	if page != "" {
		pages = append(pages, page)
	}

	return pages
}

func (s *songService) GetSongText(ctx context.Context, songID, pageSize, pageNumber int) (string, error) {
	startTime := time.Now()
	logger.Logger.Debugf("Fetching song text for ID: %d, pageSize: %d, pageNumber: %d", songID, pageSize, pageNumber)
	// Получаем текст песни
	songText, err := s.repo.GetSongText(ctx, songID)
	if err != nil {
		logger.Logger.Error("Error retrieving song text: ", err)
		return "", err
	}

	if songText != "" {
		songText = strings.ReplaceAll(songText, "\\n", "\n")
	}
	pages := paginateText(songText, pageSize)

	if pageNumber <= 0 || pageNumber > len(pages) {
		return "", fmt.Errorf("Invalid page number")
	}
	logger.Logger.Infof("GetSongText executed successfully, execution time: %s", time.Since(startTime))
	return pages[pageNumber-1], nil
}

// Получение песен с фильтром и логированием
func (s *songService) GetSongs(ctx context.Context, filter string, limit, offset int) ([]models.Song, error) {
	startTime := time.Now()
	logger.Logger.Debugf("Retrieving songs with filter: %s, limit: %d, offset: %d", filter, limit, offset)

	songs, err := s.repo.GetSongs(ctx, filter, limit, offset)
	if err != nil {
		logger.Logger.Error("Error retrieving songs: ", err)
		return nil, err
	}

	logger.Logger.Infof("GetSongs executed successfully, retrieved %d songs, execution time: %s", len(songs), time.Since(startTime))
	return songs, nil
}

// Удаление песни с логикой проверки
func (s *songService) DeleteSong(ctx context.Context, songID int64) error {
	startTime := time.Now()
	logger.Logger.Debugf("Attempting to delete song with ID: %d", songID)

	err := s.repo.DeleteSong(ctx, songID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Logger.Warnf("Song with ID %d not found", songID)
			return fmt.Errorf("song with id %d not found", songID)
		}
		logger.Logger.Error("Error deleting song: ", err)
		return err
	}

	logger.Logger.Infof("DeleteSong executed successfully, song ID: %d deleted, execution time: %s", songID, time.Since(startTime))
	return nil
}

// Обновление песни с логикой проверки
func (s *songService) UpdateSong(ctx context.Context, updSong models.SongUpdateParams) error {
	startTime := time.Now()
	logger.Logger.Debugf("Updating song with ID: %d, data: %+v", updSong.ID, updSong)

	err := s.repo.UpdateSong(ctx, updSong)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Logger.Warnf("Song with ID %d not found", updSong.ID)
			return fmt.Errorf("song with id %d not found", updSong.ID)
		}
		logger.Logger.Error("Error updating song: ", err)
		return err
	}

	logger.Logger.Infof("UpdateSong executed successfully, song ID: %d updated, execution time: %s", updSong.ID, time.Since(startTime))
	return nil
}
