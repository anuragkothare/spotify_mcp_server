package spotify

type Track struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	URI    string `json:"uri"`
}

type Artist struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
	URI        string `json:"uri"`
}

type SearchResult struct {
	Tracks []Track `json:"tracks"`
	Total  int     `json:"total"`
}

type ArtistSearchResult struct {
	Artists []Artist `json:"artists"`
	Total   int      `json:"total"`
}
