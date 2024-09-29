package postgresrepo

import (
	"context"
	"database/sql"
	"musPlayer/models"
)

type songRepository struct {
	db *sql.DB
}

func NewSongRepository(db *sql.DB) SongRepository {
	return &songRepository{
		db: db,
	}
}

type AddSongParams struct {
	GroupName   string
	SongId      int
	SongName    string
	Text        string
	ReleaseDate string
	Link        string
}

// Добавление песни
func (r *songRepository) AddSong(ctx context.Context, song AddSongParams) (int, error) {
	var id int
	query := `INSERT INTO songs (song_id, group_name, song_name, text, release_date, link ) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, song.SongId, song.GroupName, song.SongName, song.Text, song.ReleaseDate, song.Link).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Получение текста песни
func (r *songRepository) GetSongText(ctx context.Context, songID int) (string, error) {
	query := `SELECT text FROM songs WHERE id = $1`

	var songText string
	err := r.db.QueryRowContext(ctx, query, songID).Scan(&songText)
	if err != nil {
		return "", err
	}

	return songText, nil
}

// Получение списка песен
func (r *songRepository) GetSongs(ctx context.Context, filter string, limit, offset int) ([]models.Song, error) {
	query := `SELECT id, group_name, song_name, text, release_date FROM songs WHERE group_name ILIKE $1 LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, "%"+filter+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.Text, &song.ReleaseDate); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

// Удаление песни
func (r *songRepository) DeleteSong(ctx context.Context, songID int64) error {
	query := `DELETE FROM songs WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, songID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Обновление песни
func (r *songRepository) UpdateSong(ctx context.Context, song models.Song) error {
	query := `UPDATE songs 
              SET group_name = $1, song_name = $2, text = $3, release_date = $4
              WHERE id = $5`

	result, err := r.db.ExecContext(ctx, query, song.GroupName, song.SongName, song.Text, song.ReleaseDate, song.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
