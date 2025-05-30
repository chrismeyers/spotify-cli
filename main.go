package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	API struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	} `json:"api"`
	TokenPath string
}

type RawToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Token struct {
	RawToken
	Expiration int `json:"expiration"`
}

type Client struct {
	Config Config
}

type SearchQuery struct {
	Q               string
	Type            string
	Market          string
	Limit           int
	Offset          int
	IncludeExternal string
}

type SearchResults struct {
	Tracks struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			Album struct {
				AlbumType        string   `json:"album_type"`
				TotalTracks      int      `json:"total_tracks"`
				AvailableMarkets []string `json:"available_markets"`
				ExternalUrls     struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href   string `json:"href"`
				ID     string `json:"id"`
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
				Restrictions         struct {
					Reason string `json:"reason"`
				} `json:"restrictions"`
				Type    string `json:"type"`
				URI     string `json:"uri"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
			} `json:"album"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			DiscNumber       int      `json:"disc_number"`
			DurationMs       int      `json:"duration_ms"`
			Explicit         bool     `json:"explicit"`
			ExternalIds      struct {
				Isrc string `json:"isrc"`
				Ean  string `json:"ean"`
				Upc  string `json:"upc"`
			} `json:"external_ids"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href       string `json:"href"`
			ID         string `json:"id"`
			IsPlayable bool   `json:"is_playable"`
			LinkedFrom struct {
			} `json:"linked_from"`
			Restrictions struct {
				Reason string `json:"reason"`
			} `json:"restrictions"`
			Name        string `json:"name"`
			Popularity  int    `json:"popularity"`
			PreviewURL  string `json:"preview_url"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			URI         string `json:"uri"`
			IsLocal     bool   `json:"is_local"`
		} `json:"items"`
	} `json:"tracks"`
	Artists struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Followers struct {
				Href  string `json:"href"`
				Total int    `json:"total"`
			} `json:"followers"`
			Genres []string `json:"genres"`
			Href   string   `json:"href"`
			ID     string   `json:"id"`
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name       string `json:"name"`
			Popularity int    `json:"popularity"`
			Type       string `json:"type"`
			URI        string `json:"uri"`
		} `json:"items"`
	} `json:"artists"`
	Albums struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			AlbumType        string   `json:"album_type"`
			TotalTracks      int      `json:"total_tracks"`
			AvailableMarkets []string `json:"available_markets"`
			ExternalUrls     struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name                 string `json:"name"`
			ReleaseDate          string `json:"release_date"`
			ReleaseDatePrecision string `json:"release_date_precision"`
			Restrictions         struct {
				Reason string `json:"reason"`
			} `json:"restrictions"`
			Type    string `json:"type"`
			URI     string `json:"uri"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
		} `json:"items"`
	} `json:"albums"`
	Playlists struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			Collaborative bool   `json:"collaborative"`
			Description   string `json:"description"`
			ExternalUrls  struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			Name  string `json:"name"`
			Owner struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href        string `json:"href"`
				ID          string `json:"id"`
				Type        string `json:"type"`
				URI         string `json:"uri"`
				DisplayName string `json:"display_name"`
			} `json:"owner"`
			Public     bool   `json:"public"`
			SnapshotID string `json:"snapshot_id"`
			Tracks     struct {
				Href  string `json:"href"`
				Total int    `json:"total"`
			} `json:"tracks"`
			Type string `json:"type"`
			URI  string `json:"uri"`
		} `json:"items"`
	} `json:"playlists"`
	Shows struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			AvailableMarkets []string `json:"available_markets"`
			Copyrights       []struct {
				Text string `json:"text"`
				Type string `json:"type"`
			} `json:"copyrights"`
			Description     string `json:"description"`
			HTMLDescription string `json:"html_description"`
			Explicit        bool   `json:"explicit"`
			ExternalUrls    struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			IsExternallyHosted bool     `json:"is_externally_hosted"`
			Languages          []string `json:"languages"`
			MediaType          string   `json:"media_type"`
			Name               string   `json:"name"`
			Publisher          string   `json:"publisher"`
			Type               string   `json:"type"`
			URI                string   `json:"uri"`
			TotalEpisodes      int      `json:"total_episodes"`
		} `json:"items"`
	} `json:"shows"`
	Episodes struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			AudioPreviewURL string `json:"audio_preview_url"`
			Description     string `json:"description"`
			HTMLDescription string `json:"html_description"`
			DurationMs      int    `json:"duration_ms"`
			Explicit        bool   `json:"explicit"`
			ExternalUrls    struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			IsExternallyHosted   bool     `json:"is_externally_hosted"`
			IsPlayable           bool     `json:"is_playable"`
			Language             string   `json:"language"`
			Languages            []string `json:"languages"`
			Name                 string   `json:"name"`
			ReleaseDate          string   `json:"release_date"`
			ReleaseDatePrecision string   `json:"release_date_precision"`
			ResumePoint          struct {
				FullyPlayed      bool `json:"fully_played"`
				ResumePositionMs int  `json:"resume_position_ms"`
			} `json:"resume_point"`
			Type         string `json:"type"`
			URI          string `json:"uri"`
			Restrictions struct {
				Reason string `json:"reason"`
			} `json:"restrictions"`
		} `json:"items"`
	} `json:"episodes"`
	Audiobooks struct {
		Href     string `json:"href"`
		Limit    int    `json:"limit"`
		Next     string `json:"next"`
		Offset   int    `json:"offset"`
		Previous string `json:"previous"`
		Total    int    `json:"total"`
		Items    []struct {
			Authors []struct {
				Name string `json:"name"`
			} `json:"authors"`
			AvailableMarkets []string `json:"available_markets"`
			Copyrights       []struct {
				Text string `json:"text"`
				Type string `json:"type"`
			} `json:"copyrights"`
			Description     string `json:"description"`
			HTMLDescription string `json:"html_description"`
			Edition         string `json:"edition"`
			Explicit        bool   `json:"explicit"`
			ExternalUrls    struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			Languages []string `json:"languages"`
			MediaType string   `json:"media_type"`
			Name      string   `json:"name"`
			Narrators []struct {
				Name string `json:"name"`
			} `json:"narrators"`
			Publisher     string `json:"publisher"`
			Type          string `json:"type"`
			URI           string `json:"uri"`
			TotalChapters int    `json:"total_chapters"`
		} `json:"items"`
	} `json:"audiobooks"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func NewClient(config Config) Client {
	return Client{Config: config}
}

func (c *Client) fetchToken() (*Token, error) {
	if _, err := os.Stat(c.Config.TokenPath); !errors.Is(err, os.ErrNotExist) {
		data, err := os.ReadFile(c.Config.TokenPath)
		if err != nil {
			return nil, err
		}

		var cachedToken Token
		err = json.Unmarshal(data, &cachedToken)
		if err != nil {
			return nil, err
		}

		if cachedToken.Expiration >= int(time.Now().Unix()) {
			return &cachedToken, nil
		}
	}

	reqBody := url.Values{
		"grant_type":    []string{"client_credentials"},
		"client_id":     []string{c.Config.API.ClientID},
		"client_secret": []string{c.Config.API.ClientSecret},
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://accounts.spotify.com/api/token",
		strings.NewReader(reqBody.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawToken RawToken
	err = json.Unmarshal(respBody, &rawToken)
	if err != nil {
		return nil, err
	}

	token := Token{
		RawToken:   rawToken,
		Expiration: int(time.Now().Unix()) + rawToken.ExpiresIn - 15,
	}

	tokenStr, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(c.Config.TokenPath, tokenStr, 0644)

	return &token, nil
}

func (c *Client) search(s SearchQuery) (*SearchResults, error) {
	token, err := c.fetchToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/search", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	q := req.URL.Query()
	q.Add("q", s.Q)
	q.Add("type", s.Type)
	if s.Market != "" {
		q.Add("market", s.Market)
	}
	if s.Limit != 0 {
		q.Add("limit", strconv.Itoa(s.Limit))
	}
	if s.Offset != 0 {
		q.Add("offset", strconv.Itoa(s.Offset))
	}
	if s.IncludeExternal != "" {
		q.Add("include_external", s.IncludeExternal)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var results SearchResults
	err = json.Unmarshal(respBody, &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

func main() {
	config, err := loadConfig("./config.json")
	if err != nil {
		panic(err)
	}
	config.TokenPath = "./token.json"

	client := NewClient(*config)

	results, err := client.search(SearchQuery{Q: os.Args[1], Type: os.Args[2]})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", results)
}
