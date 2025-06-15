package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	API struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	} `json:"api"`
	TokenPath string
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

var (
	docStyle      = lipgloss.NewStyle().Margin(1, 2)
	categoryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#1DB954")).Bold(true)
	footerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#767676")).Faint(true)
)

var bands = []string{
	"The Beatles", "Led Zeppelin", "Pink Floyd", "Queen", "The Rolling Stones",
	"Nirvana", "Radiohead", "AC/DC", "The Who", "Black Sabbath",
	"Deep Purple", "Metallica", "Iron Maiden", "The Doors", "Jimi Hendrix",
	"Bob Dylan", "David Bowie", "The Clash", "Sex Pistols", "Ramones",
	"Pearl Jam", "Soundgarden", "Alice in Chains", "Red Hot Chili Peppers",
	"Foo Fighters", "Green Day", "U2", "R.E.M.", "The Smiths",
	"Joy Division", "New Order", "Depeche Mode", "The Cure", "Blur",
	"Oasis", "Coldplay", "Arctic Monkeys", "The Strokes", "The White Stripes",
	"Kings of Leon", "Arcade Fire", "Vampire Weekend", "Tame Impala",
	"Fleetwood Mac", "Eagles", "Lynyrd Skynyrd", "Creedence Clearwater Revival",
}

var songs = []string{
	"Bohemian Rhapsody", "Stairway to Heaven", "Hotel California", "Imagine",
	"Smells Like Teen Spirit", "Sweet Child O' Mine", "Billie Jean",
	"Like a Rolling Stone", "Purple Haze", "Hey Jude", "What's Going On",
	"Born to Run", "Good Vibrations", "Respect", "Johnny B. Goode",
	"Satisfaction", "My Generation", "God Only Knows", "A Change Is Gonna Come",
	"Dancing Queen", "Don't Stop Believin'", "Sweet Caroline", "Piano Man",
	"Wonderwall", "Creep", "Losing My Religion", "Black", "Everlong",
	"Mr. Brightside", "Seven Nation Army", "Crazy", "Hallelujah",
	"The Sound of Silence", "Bridge Over Troubled Water", "Yesterday",
	"What a Wonderful World", "Unchained Melody", "Stand by Me",
	"Lean on Me", "I Want to Hold Your Hand", "Help!", "Come As You Are",
	"Enter Sandman", "Sweet Dreams", "Take On Me", "Don't Stop Me Now",
	"We Will Rock You", "We Are the Champions", "Another Brick in the Wall",
}

func getRandomSearchTerm() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	allItems := make([]string, 0, len(bands)+len(songs))
	allItems = append(allItems, bands...)
	allItems = append(allItems, songs...)

	return allItems[r.Intn(len(allItems))]
}

type choice struct {
	name       string
	searchType string
	selected   bool
}

type resultItem struct {
	category string
	name     string
	detail   string
	url      string
}

func (i resultItem) Title() string { return i.name }
func (i resultItem) Description() string {
	return fmt.Sprintf("%s · %s", categoryStyle.Render(i.category), i.detail)
}
func (i resultItem) FilterValue() string { return i.name }

type ViewState int

const (
	SearchView ViewState = iota
	ResultsView
)

type model struct {
	sub        chan SearchResults
	client     Client
	textInput  textinput.Model
	choices    []choice
	cursor     int
	spinner    spinner.Model
	loading    bool
	results    *SearchResults
	resultList list.Model
	error      string
	view       ViewState
}

func initialModel(config *Config) model {
	client := NewClient(*config)

	ti := textinput.New()
	ti.Placeholder = getRandomSearchTerm()
	ti.Prompt = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	s := spinner.New()
	s.Spinner = spinner.Dot

	items := []list.Item{}
	l := list.New(items, list.NewDefaultDelegate(), 40, 2)
	l.Title = "Search Results"
	l.DisableQuitKeybindings()

	return model{
		sub:       make(chan SearchResults),
		client:    client,
		textInput: ti,
		choices: []choice{
			{name: "Album", searchType: "album", selected: false},
			{name: "Artist", searchType: "artist", selected: false},
			{name: "Playlist", searchType: "playlist", selected: false},
			{name: "Track", searchType: "track", selected: false},
			{name: "Show", searchType: "show", selected: false},
			{name: "Episode", searchType: "episode", selected: false},
			{name: "Audiobook", searchType: "audiobook", selected: false},
		},
		spinner:    s,
		loading:    false,
		resultList: l,
		error:      "",
		view:       SearchView,
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
			if m.view == SearchView && m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.view == SearchView && m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "right", "left":
			if m.view == SearchView {
				m.choices[m.cursor].selected = !m.choices[m.cursor].selected
			}
		case "enter":
			if m.view == SearchView {
				var types []string
				for _, choice := range m.choices {
					if choice.selected == true {
						types = append(types, choice.searchType)
					}
				}

				input := m.textInput.Value()
				typeStr := strings.Join(types, ",")

				if input == "" {
					m.error = "Please enter a search term"
					return m, nil
				}
				if typeStr == "" {
					m.error = "Please select at least one category"
					return m, nil
				}

				m.error = ""
				m.results = nil
				m.loading = !m.loading
				cmd = m.spinner.Tick

				go func() {
					results, _ := m.client.search(SearchQuery{Q: input, Type: typeStr})
					m.sub <- *results
				}()
			}
		case "esc":
			if m.view == ResultsView {
				if m.resultList.FilterState() == list.Filtering {
					m.resultList.ResetFilter()
				} else {
					m.view = SearchView
				}
				return m, nil
			}
		default:
			if m.view == SearchView {
				m.textInput, cmd = m.textInput.Update(msg)
			} else if m.view == ResultsView {
				if msg.String() == "q" && m.resultList.FilterState() != list.Filtering {
					m.view = SearchView
					return m, nil
				}
			}
		}
	case SearchResults:
		m.loading = !m.loading
		m.results = &msg

		const maxItems = 10

		var items []list.Item
		for i, a := range msg.Albums.Items {
			if i >= maxItems {
				break
			}
			artistNames := []string{}
			for _, ar := range a.Artists {
				artistNames = append(artistNames, ar.Name)
			}
			items = append(items, resultItem{
				category: "Album",
				name:     a.Name,
				detail:   fmt.Sprintf("by %s · Released: %s", strings.Join(artistNames, ", "), a.ReleaseDate),
				url:      a.ExternalUrls.Spotify,
			})
		}
		for i, a := range msg.Artists.Items {
			if i >= maxItems {
				break
			}
			items = append(items, resultItem{
				category: "Artist",
				name:     a.Name,
				detail:   fmt.Sprintf("Genres: %s", strings.Join(a.Genres, ", ")),
				url:      a.ExternalUrls.Spotify,
			})
		}
		for i, p := range msg.Playlists.Items {
			if i >= maxItems {
				break
			}
			items = append(items, resultItem{
				category: "Playlist",
				name:     p.Name,
				detail:   fmt.Sprintf("by %s · %d tracks", p.Owner.DisplayName, p.Tracks.Total),
				url:      p.ExternalUrls.Spotify,
			})
		}
		for i, t := range msg.Tracks.Items {
			if i >= maxItems {
				break
			}
			artistNames := []string{}
			for _, a := range t.Artists {
				artistNames = append(artistNames, a.Name)
			}
			items = append(items, resultItem{
				category: "Track",
				name:     t.Name,
				detail:   fmt.Sprintf("by %s · Album: %s", strings.Join(artistNames, ", "), t.Album.Name),
				url:      t.ExternalUrls.Spotify,
			})
		}
		for i, s := range msg.Shows.Items {
			if i >= maxItems {
				break
			}
			items = append(items, resultItem{
				category: "Show",
				name:     s.Name,
				detail:   fmt.Sprintf("by %s", s.Publisher),
				url:      s.ExternalUrls.Spotify,
			})
		}
		for i, e := range msg.Episodes.Items {
			if i >= maxItems {
				break
			}
			items = append(items, resultItem{
				category: "Episode",
				name:     e.Name,
				detail:   fmt.Sprintf("by %s", e.Name),
				url:      e.ExternalUrls.Spotify,
			})
		}
		for i, a := range msg.Audiobooks.Items {
			if i >= maxItems {
				break
			}
			items = append(items, resultItem{
				category: "Audiobook",
				name:     a.Name,
				detail:   fmt.Sprintf("by %s", a.Authors[0].Name),
				url:      a.ExternalUrls.Spotify,
			})
		}

		m.resultList.SetItems(items)
		m.resultList.Select(0)
		m.view = ResultsView

		return m, waitForActivity(m.sub)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.resultList.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.view == ResultsView {
		m.resultList, cmd = m.resultList.Update(msg)
	}

	return m, cmd
}

func (m model) searchView() string {
	var s strings.Builder

	s.WriteString("Spotify Search\n\n")

	s.WriteString("Search: ")
	s.WriteString(m.textInput.View())
	s.WriteString("\n\n")

	s.WriteString("Search Types:\n")
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if choice.selected {
			checked = "x"
		}

		s.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.name))
	}

	if m.error != "" {
		s.WriteString(fmt.Sprintf("\nError: %s\n", m.error))
	}

	if m.loading {
		s.WriteString(fmt.Sprintf("\n%s Loading...\n", m.spinner.View()))
	}

	s.WriteString(footerStyle.Render("\nUse arrow keys to select categories and Enter to search."))
	s.WriteString(footerStyle.Render("\nPress Ctrl-C to quit."))

	return s.String()
}

func (m model) resultsView() string {
	var s strings.Builder

	s.WriteString("\n")
	s.WriteString(m.resultList.View())

	s.WriteString(footerStyle.Render("\nPress 'q' or 'esc' to go back. Press Ctrl-C to quit."))

	return s.String()
}

func (m model) View() string {
	if m.view == ResultsView {
		return m.resultsView()
	}
	return m.searchView()
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	possiblePaths := []string{
		os.Getenv("XDG_CONFIG_HOME") + "/spotify-cli/config.json",
		os.Getenv("HOME") + "/.config/spotify-cli/config.json",
		"./config.json",
	}

	var config *Config
	var err error
	for _, path := range possiblePaths {
		if _, statErr := os.Stat(path); statErr == nil {
			config, err = loadConfig(path)
			if err != nil {
				panic(err)
			}
			config.TokenPath = filepath.Dir(path) + "/token.json"
			break
		}
	}
	if config == nil {
		panic("config.json not found in any common location")
	}

	p := tea.NewProgram(
		initialModel(config),
		tea.WithAltScreen(),       // Enable full screen mode
		tea.WithMouseCellMotion(), // Enable mouse support
	)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
