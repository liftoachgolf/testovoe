package postgresrepo

import (
	"context"
	"database/sql"
	"musPlayer/models"
)

type SongRepository interface {
	AddSong(ctx context.Context, song AddSongParams) (int, error)
	GetSongs(ctx context.Context, filter string, limit, offset int) ([]models.Song, error)
	DeleteSong(ctx context.Context, songID int64) error
	UpdateSong(ctx context.Context, updSong models.SongUpdateParams) error
	GetSongText(ctx context.Context, songID int) (string, error)
}
type Repository struct {
	SongRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		SongRepository: NewSongRepository(db),
	}
}
