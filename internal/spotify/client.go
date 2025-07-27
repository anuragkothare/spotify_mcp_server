package spotify

import (
	"context"
	"fmt"

	"github.com/anuragkothare/spotify_mcp_server/internal/config"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	client *spotify.Client
	auth   *spotifyauth.Authenticator
}

func NewClient(cfg config.SpotifyConfig) (*Client, error) {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(cfg.ClientID),
		spotifyauth.WithClientSecret(cfg.ClientSecret),
		spotifyauth.WithRedirectURL(cfg.RedirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadEmail,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopeUserTopRead,
		),
	)

	// Use client credentials flow for app-only access
	config := &clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	httpClient := spotifyauth.New().Client(context.Background(), token)
	client := spotify.New(httpClient)

	return &Client{
		client: client,
		auth:   auth,
	}, nil
}

func (c *Client) SearchTracks(query string, limit int) (*SearchResult, error) {
	results, err := c.client.Search(context.Background(), query, spotify.SearchTypeTrack, spotify.Limit(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to search tracks: %w", err)
	}

	tracks := make([]Track, len(results.Tracks.Tracks))
	for i, track := range results.Tracks.Tracks {
		// Handle case where track might not have artists
		artistName := "Unknown Artist"
		if len(track.Artists) > 0 {
			artistName = track.Artists[0].Name
		}

		tracks[i] = Track{
			ID:     string(track.ID),
			Name:   track.Name,
			Artist: artistName,
			Album:  track.Album.Name,
			URI:    string(track.URI),
		}
	}

	return &SearchResult{
		Tracks: tracks,
		Total:  int(results.Tracks.Total), // Convert spotify.Numeric to int
	}, nil
}

func (c *Client) SearchArtists(query string, limit int) (*ArtistSearchResult, error) {
	results, err := c.client.Search(context.Background(), query, spotify.SearchTypeArtist, spotify.Limit(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to search artists: %w", err)
	}

	artists := make([]Artist, len(results.Artists.Artists))
	for i, artist := range results.Artists.Artists {
		artists[i] = Artist{
			ID:         string(artist.ID),
			Name:       artist.Name,
			Popularity: int(artist.Popularity), // Convert spotify.Numeric to int
			URI:        string(artist.URI),
		}
	}

	return &ArtistSearchResult{
		Artists: artists,
		Total:   int(results.Artists.Total), // Convert spotify.Numeric to int
	}, nil
}

func (c *Client) GetTrack(trackID string) (*Track, error) {
	track, err := c.client.GetTrack(context.Background(), spotify.ID(trackID))
	if err != nil {
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	// Handle case where track might not have artists
	artistName := "Unknown Artist"
	if len(track.Artists) > 0 {
		artistName = track.Artists[0].Name
	}

	return &Track{
		ID:     string(track.ID),
		Name:   track.Name,
		Artist: artistName,
		Album:  track.Album.Name,
		URI:    string(track.URI),
	}, nil
}
