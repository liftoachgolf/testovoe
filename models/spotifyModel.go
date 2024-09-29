package models

type SpotifyConfig struct {
	ID          string `json:"client_id"`
	Secret      string `json:"client_secret"`
	RedirectURI string `json:"redirect_uri"`
	AuthURL     string `json:"auth_url"`
	TokenURL    string `json:"token_url"`
	Scope       string `json:"scope"`
}
type SpotySong struct {
	ID          string `json:"id"`
	SongName    string `json:"song_name"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Duration    int    `json:"duration"`
	Popularity  int    `json:"popularity"`
	ReleaseDate string `json:"release_date"` // Дата релиза
	TrackURL    string `json:"track_url"`    // Ссылка на трек
}
