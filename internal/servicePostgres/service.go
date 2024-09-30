package servicePostgres

import (
	"context"
	"musPlayer/models"
	postgresrepo "musPlayer/pkg/postgresRepo"
)

type SongService interface {
	AddSong(ctx context.Context, song postgresrepo.AddSongParams) (int, error)
	GetSongText(ctx context.Context, songID, pageSize, pageNumber int) (string, error)
	GetSongs(ctx context.Context, filter string, limit, offset int) ([]models.Song, error)
	DeleteSong(ctx context.Context, songID int64) error
	UpdateSong(ctx context.Context, song models.Song) error
}

type Service struct {
	SongService
}

func NewServicePostgres(repo *postgresrepo.Repository) *Service {
	return &Service{
		SongService: NewSongService(repo.SongRepository),
	}
}
