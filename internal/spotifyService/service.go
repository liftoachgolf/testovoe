package spotifyService

import (
	"encoding/json"
	"fmt"
	"musPlayer/internal/logger"
	"musPlayer/models"
	"net/http"
	"net/url"
	"strings"
)

type SpotifyService struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AccessToken  string
}

func NewSpotifyService(clientID, clientSecret, redirectURI string) *SpotifyService {
	return &SpotifyService{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
}

// RedirectUser перенаправляет пользователя на страницу авторизации Spotify
func (s *SpotifyService) RedirectUser(w http.ResponseWriter, r *http.Request) {
	if s.ClientID == "" || s.RedirectURI == "" {
		http.Error(w, "Service is not configured properly", http.StatusInternalServerError)
		return
	}

	authURL := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=user-read-private user-read-email",
		s.ClientID,
		url.QueryEscape(s.RedirectURI),
	)

	http.Redirect(w, r, authURL, http.StatusFound)
}

// GetAccessToken получает токен доступа
func (s *SpotifyService) GetAccessToken(code string) error {
	form := url.Values{
		"code":          {code},
		"client_id":     {s.ClientID},
		"client_secret": {s.ClientSecret},
		"redirect_uri":  {s.RedirectURI},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.PostForm("https://accounts.spotify.com/api/token", form)
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to request access token")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("failed to fetch access token: %s", resp.Status)
		logger.Logger.Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Logger.WithError(err).Error("Failed to decode access token response")
		return err
	}

	s.AccessToken = result.AccessToken
	logger.Logger.Info("Successfully obtained access token")
	return nil
}

// FetchSongDetails получает информацию о треке по его ID
func (s *SpotifyService) FetchSongDetails(trackID string) (models.Song, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to create request")
		return models.Song{}, err
	}
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to execute request")
		return models.Song{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("failed to fetch song details: %s", resp.Status)
		logger.Logger.Error(errMsg)
		return models.Song{}, fmt.Errorf(errMsg)
	}

	var song models.Song
	if err := json.NewDecoder(resp.Body).Decode(&song); err != nil {
		logger.Logger.WithError(err).Error("Failed to decode response")
		return models.Song{}, err
	}

	logger.Logger.Infof("Successfully fetched details for song: %s", song.SongName)
	return song, nil
}

func (s *SpotifyService) SearchTrack(query string) ([]models.SpotySong, error) {
	// Формируем URL для поиска треков
	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(query))

	// Создаем новый GET-запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to create search request")
		return nil, err
	}

	// Устанавливаем заголовок Authorization с токеном доступа
	req.Header.Set("Authorization", "Bearer "+s.AccessToken)

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to execute search request")
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Failed to search tracks: %s", resp.Status)
		logger.Logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// Декодируем ответ
	var result struct {
		Tracks struct {
			Items []struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
				Album struct {
					Name        string `json:"name"`
					ReleaseDate string `json:"release_date"` // Дата релиза альбома
				} `json:"album"`
				ExternalUrls struct {
					Spotify string `json:"spotify"` // Ссылка на трек
				} `json:"external_urls"`
				DurationMs int `json:"duration_ms"`
				Popularity int `json:"popularity"`
			} `json:"items"`
		} `json:"tracks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Logger.WithError(err).Error("Failed to decode search response")
		return nil, err
	}

	// Формируем результат
	var songs []models.SpotySong
	for _, item := range result.Tracks.Items {
		song := models.SpotySong{
			ID:          item.ID,
			SongName:    item.Name,
			Artist:      item.Artists[0].Name,
			Album:       item.Album.Name,
			Duration:    item.DurationMs,
			Popularity:  item.Popularity,
			ReleaseDate: item.Album.ReleaseDate,    // Дата релиза
			TrackURL:    item.ExternalUrls.Spotify, // Ссылка на трек
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (s *SpotifyService) SearchTrackByNameAndArtist(songName, artistName string) ([]models.SpotySong, error) {
	// Формируем URL для поиска треков по названию и исполнителю
	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=track:%s artist:%s&type=track",
		strings.ReplaceAll(url.QueryEscape(songName), "+", "%20"),
		strings.ReplaceAll(url.QueryEscape(artistName), "+", "%20"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search tracks: %s", resp.Status)
	}

	var result struct {
		Tracks struct {
			Items []struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
				Album struct {
					Name string `json:"name"`
				} `json:"album"`
				DurationMs   int `json:"duration_ms"`
				Popularity   int `json:"popularity"`
				ExternalURLs struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
			} `json:"items"`
		} `json:"tracks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var songs []models.SpotySong
	for _, item := range result.Tracks.Items {
		song := models.SpotySong{
			ID:         item.ID,
			SongName:   item.Name,
			Artist:     item.Artists[0].Name,
			Album:      item.Album.Name,
			Duration:   item.DurationMs,
			Popularity: item.Popularity,
			TrackURL:   item.ExternalURLs.Spotify, // Ссылка на трек
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (s *SpotifyService) RefreshAccessToken(refreshToken string) error {
	// Формируем URL для обновления токена
	tokenURL := "https://accounts.spotify.com/api/token"

	// Создаем url.Values для отправки данных
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", s.ClientID)
	data.Set("client_secret", s.ClientSecret)

	// Создаем новый POST-запрос
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовок Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to refresh access token: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to refresh token: %s", resp.Status)
	}

	// Декодируем ответ
	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`              // Время действия нового токена в секундах
		RefreshToken string `json:"refresh_token,omitempty"` // Иногда может прийти новый refresh token
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	// Обновляем токен в сервисе
	s.AccessToken = tokenResponse.AccessToken

	return nil
}
