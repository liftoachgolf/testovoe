package servicegenius

import (
	"encoding/json"
	"fmt"
	"io"
	"musPlayer/internal/logger"
	"musPlayer/models"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type GeniusService struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AccessToken  string
}

func NewGeniusService(clientID, clientSecret, redirectURI string) *GeniusService {
	return &GeniusService{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
}

// RedirectUser перенаправляет пользователя на страницу авторизации Genius
func (g *GeniusService) RedirectUser(w http.ResponseWriter, r *http.Request) {
	if g.ClientID == "" || g.RedirectURI == "" {
		http.Error(w, "Service is not configured properly", http.StatusInternalServerError)
		return
	}

	authURL := fmt.Sprintf("https://api.genius.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code",
		g.ClientID,
		url.QueryEscape(g.RedirectURI),
	)

	http.Redirect(w, r, authURL, http.StatusFound)
}

// GetAccessToken получает токен доступа
func (g *GeniusService) GetAccessToken(code string) error {
	form := url.Values{
		"code":          {code},
		"client_id":     {g.ClientID},
		"client_secret": {g.ClientSecret},
		"redirect_uri":  {g.RedirectURI},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.PostForm("https://api.genius.com/oauth/token", form)
	if err != nil {
		fmt.Printf("Failed to request access token: %v\n", err)
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
		logger.Logger.Error("Failed to decode access token response: ", err)
		return err
	}

	g.AccessToken = result.AccessToken
	fmt.Println("Successfully obtained access token")
	return nil
}

// SearchSong ищет песню по названию и заполняет структуру Song
func (g *GeniusService) SearchSong(title, artist string) (*models.Song, error) {
	query := url.QueryEscape(title)
	if artist != "" {
		query += " " + url.QueryEscape(artist)
	}

	url := fmt.Sprintf("https://api.genius.com/search?q=%s", query)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Error("Failed to perform search: ", err)
		return nil, fmt.Errorf("failed to perform search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch song: %s", resp.Status)
	}

	var result struct {
		Response struct {
			Hits []struct {
				Result struct {
					ID            int    `json:"id"`
					Title         string `json:"title"`
					PrimaryArtist struct {
						Name string `json:"name"`
					} `json:"primary_artist"`
					ReleaseDate string `json:"release_date"`
					URL         string `json:"url"`
				} `json:"result"`
			} `json:"hits"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Logger.Error("Failed to decode search response: ", err)
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	if len(result.Response.Hits) == 0 {
		return nil, fmt.Errorf("no song found for title: %s and artist: %s", title, artist)
	}

	songID := result.Response.Hits[0].Result.ID
	songText, err := g.GetSongText(result.Response.Hits[0].Result.URL)
	if err != nil {
		return nil, err
	}

	song := &models.Song{
		ID:          songID,
		GroupName:   result.Response.Hits[0].Result.PrimaryArtist.Name,
		SongName:    result.Response.Hits[0].Result.Title,
		ReleaseDate: result.Response.Hits[0].Result.ReleaseDate,
		Link:        result.Response.Hits[0].Result.URL,
		Text:        songText,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return song, nil
}

func extractTextFromHTML(htmlBody string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга HTML: %v", err)
	}

	var text string
	var found bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, attr := range n.Attr {
				if attr.Key == "data-lyrics-container" && attr.Val == "true" {
					found = true
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.TextNode {
							line := strings.TrimSpace(c.Data)
							if line != "" {
								text += line + "\n"
							}
						} else if c.Type == html.ElementNode && c.Data == "br" {
							text += "\n"
						}
					}
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if !found {
		return "", fmt.Errorf("не удалось найти элемент с data-lyrics-container")
	}

	if text == "" {
		return "", fmt.Errorf("не удалось извлечь текст из найденного элемента")
	}

	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n\n", "\n")

	return text, nil
}

// GetSongText получает текст песни по идентификатору песни
func (g *GeniusService) GetSongText(url string) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+g.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Error("Failed to perform song text request: ", err)
		return "", fmt.Errorf("failed to perform song text request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch song text: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Error("Failed to read response body: ", err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	songText, err := extractTextFromHTML(string(body))
	if err != nil {
		logger.Logger.Error("Failed to extract song text: ", err)
		return "", fmt.Errorf("failed to extract song text: %w", err)
	}

	return songText, nil
}
