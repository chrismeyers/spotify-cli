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

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

type choice struct {
	name       string
	searchType string
	selected   bool
}

type model struct {
	sub       chan SearchResults
	client    Client
	textInput textinput.Model
	choices   []choice
	cursor    int
	spinner   spinner.Model
	loading   bool
	results   *SearchResults
}

func initialModel(config *Config) model {
	client := NewClient(*config)

	ti := textinput.New()
	ti.Placeholder = ""
	ti.Prompt = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		sub:       make(chan SearchResults),
		client:    client,
		textInput: ti,
		choices: []choice{
			{name: "Album", searchType: "album", selected: false},
			{name: "Artist", searchType: "artist", selected: true},
			{name: "Playlist", searchType: "playlist", selected: false},
			{name: "Track", searchType: "track", selected: false},
			{name: "Show", searchType: "show", selected: false},
			{name: "Episode", searchType: "episode", selected: false},
			{name: "Audiobook", searchType: "audiobook", selected: false},
		},
		spinner: s,
		loading: false,
	}
}

func waitForActivity(sub chan SearchResults) tea.Cmd {
	return func() tea.Msg {
		return SearchResults(<-sub)
	}
}

func (m model) Init() tea.Cmd {
	return waitForActivity(m.sub)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "right", "left":
			m.choices[m.cursor].selected = !m.choices[m.cursor].selected
		case "enter":
			var types []string
			for _, choice := range m.choices {
				if choice.selected == true {
					types = append(types, choice.searchType)
				}
			}

			input := m.textInput.Value()
			typeStr := strings.Join(types, ",")

			if input != "" && typeStr != "" {
				m.results = nil
				m.loading = !m.loading
				cmd = m.spinner.Tick

				go func() {
					// TODO: add error handling
					results, _ := m.client.search(SearchQuery{Q: input, Type: typeStr})
					m.sub <- *results
				}()
			} else {
				// TODO: handle error notification
			}
		default:
			m.textInput, cmd = m.textInput.Update(msg)
		}
	case SearchResults:
		m.loading = !m.loading
		m.results = &msg
		return m, waitForActivity(m.sub)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m model) View() string {
	s := "Search Spotify: "

	s += m.textInput.View() + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if choice.selected {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.name)
	}

	if m.loading {
		s += fmt.Sprintf("\n%s Loading...\n", m.spinner.View())
	}

	if m.results != nil {
		// TODO: make output nicer using tabs and lists or something
		t := table.New(
			table.WithColumns([]table.Column{
				{Title: "Category", Width: 10},
				{Title: "Hits", Width: 10},
			}),
			table.WithRows([]table.Row{
				{"Albums", strconv.Itoa(m.results.Albums.Total)},
				{"Artists", strconv.Itoa(m.results.Artists.Total)},
				{"Playlists", strconv.Itoa(m.results.Playlists.Total)},
				{"Tracks", strconv.Itoa(m.results.Tracks.Total)},
				{"Shows", strconv.Itoa(m.results.Shows.Total)},
				{"Episodes", strconv.Itoa(m.results.Episodes.Total)},
				{"Audiobooks", strconv.Itoa(m.results.Audiobooks.Total)},
			}),
			table.WithHeight(7),
			table.WithFocused(false),
		)

		t.Blur()

		styles := table.DefaultStyles()
		styles.Selected = styles.Selected.
			Foreground(styles.Cell.GetForeground()).
			Background(styles.Cell.GetBackground()).
			Bold(false)

		t.SetStyles(styles)

		s += fmt.Sprintf("\n%s\n", t.View())
	}

	s += "\nPress Ctrl-C to quit.\n"

	return s
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	config, err := loadConfig("./config.json")
	if err != nil {
		panic(err)
	}
	config.TokenPath = "./token.json"

	p := tea.NewProgram(initialModel(config))
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
