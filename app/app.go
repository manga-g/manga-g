package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	padding = 2 // Define padding for layout calculations

	// Primary UI Colors
	primaryColor   = lipgloss.Color("63")
	secondaryColor = lipgloss.Color("240")
	highlightColor = lipgloss.Color("86")
	textColor      = lipgloss.Color("252")
	errorColor     = lipgloss.Color("196")
	successColor   = lipgloss.Color("46")

	// Base Styles
	mainStyle = lipgloss.NewStyle().Padding(1, 2)
	listStyle = lipgloss.NewStyle().Padding(1, 2)
	helpStyle = lipgloss.NewStyle().Foreground(secondaryColor)

	// Text Styles
	highlightTextStyle = lipgloss.NewStyle().Foreground(highlightColor).Bold(true)
	dimmedStyle        = lipgloss.NewStyle().Foreground(secondaryColor)
	errorTextStyle     = lipgloss.NewStyle().Foreground(errorColor)
	successTextStyle   = lipgloss.NewStyle().Foreground(successColor)
	versionStyle       = lipgloss.NewStyle().Foreground(secondaryColor).Align(lipgloss.Right)

	// Styles for different UI elements
	appStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Background(lipgloss.Color("#1F2937")) // Dark blue-gray

	titleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Bold(true).
			Foreground(lipgloss.Color("170")) // Purple

	inputBoxStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(lipgloss.Color("#1F2937")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Width(50)

	menuStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(lipgloss.Color("#1F2937")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Width(50)

	infoStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(lipgloss.Color("#374151")).
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#10B981")).
			Width(50)

	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red

	menuHighlightStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1F2937")).
				Background(primaryColor).
				Bold(true).
				Padding(0, 1)

	downloadHeaderStyle = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	progressAreaStyle   = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	successStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("46")) // Green

	loadingStyle = lipgloss.NewStyle().Padding(1, 2)

	// Download Progress Styles
	downloadStatusContainerStyle = lipgloss.NewStyle().MarginTop(1).Padding(0, 1).Border(lipgloss.RoundedBorder(), true).BorderForeground(secondaryColor)
)

// Item represents a selectable item in our list
type Item struct {
	id    string // Use ID for fetching details/chapters
	title string
	desc  string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

// Define view states
const (
	ViewMain int = iota
	ViewSearchResults
	ViewChapters
	ViewReading
)

// Messages used by the application
type SearchResultsMsg MangaList
type ChaptersMsg MangaChapters
type ErrorMsg struct{ err error }
type searchErrMsg struct{ err error }

func (e searchErrMsg) Error() string { return e.err.Error() }

type searchResultsMsg struct {
	items   []list.Item
	results interface{}
}
type chaptersMsg struct {
	manga    *Manga
	chapters []Item
}

// Add new message type for download results
type atHomeResponseMsg struct {
	info         *AtHomeResponse
	chapterID    string // Include chapterID for context
	mangaTitle   string
	chapterTitle string
}

// --- App State ---
type appState int

const (
	searchState       appState = iota
	mangaDetailsState          // Used when viewing manga info before chapters
	chapterListState           // Used when viewing the chapter list
	loadingState               // Used for any loading action
	errorState                 // Used to display errors
)

// --- KeyMap ---
// (Define KeyMap and DefaultKeyMap)
type KeyMap struct {
	// Shared
	Help key.Binding
	Quit key.Binding
	Back key.Binding

	// Search State
	Submit key.Binding

	// List States (Manga/Chapter)
	Select key.Binding // Select item (view details / download chapter)
	Up     key.Binding
	Down   key.Binding
	Filter key.Binding // Optional for filtering lists
}

var DefaultKeyMap = KeyMap{
	Help: key.NewBinding(key.WithKeys("?", "h"), key.WithHelp("?", "help")),
	Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Back: key.NewBinding(key.WithKeys("esc", "backspace", "b"), key.WithHelp("esc", "back")),

	Submit: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "search")),

	Select: key.NewBinding(key.WithKeys("enter", "d"), key.WithHelp("enter/d", "select/download")),
	Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Filter: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")), // Assuming list supports filtering
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k KeyMap) ShortHelp() []key.Binding {
	// TODO: Make context-aware based on state
	return []key.Binding{k.Help, k.Quit, k.Back}
}

// FullHelp returns keybindings for the expanded help view.
func (k KeyMap) FullHelp() [][]key.Binding {
	// TODO: Make context-aware based on state
	return [][]key.Binding{
		{k.Help, k.Quit, k.Back},           // Global
		{k.Submit},                         // Search state
		{k.Select, k.Up, k.Down, k.Filter}, // List states
	}
}

// --- Model ---
// Uses the refactored structure
type Model struct {
	keys         KeyMap
	help         help.Model
	state        appState
	input        textinput.Model
	list         list.Model // Used for both search results and chapters
	spinner      spinner.Model
	loadingMsg   string
	selectedItem Item   // Store the selected manga or chapter Item
	mangaInfo    *Manga // Store detailed manga info when fetched
	chapters     []Item // Store parsed chapter list items
	downloader   *Downloader
	err          error
	width        int
	height       int
	basedApiUrl  string // MangaDex API base URL

	// Download tracking
	progressModels  map[string]progress.Model
	activeDownloads map[string]DownloadProgressInfo

	// Store necessary info temporarily for commands
	lastQuery      string // Store last search query for titles
	lastMangaID    string // Store ID of manga when viewing chapters
	lastMangaTitle string // Store title of manga when viewing chapters
}

// InitialModel creates the initial state of the application.
func InitialModel(downloader *Downloader) Model {
	searchInput := textinput.New()
	searchInput.Placeholder = "Search Manga..."
	searchInput.Focus()
	searchInput.CharLimit = 156
	searchInput.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(secondaryColor)

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = highlightTextStyle.Copy()
	delegate.Styles.SelectedDesc = dimmedStyle.Copy()
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Search Results"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()

	h := help.New()
	h.ShowAll = false

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 24
	}

	m := Model{
		keys:            DefaultKeyMap,
		help:            h,
		state:           searchState,
		input:           searchInput,
		list:            l,
		spinner:         s,
		downloader:      downloader,
		width:           width,
		height:          height,
		progressModels:  make(map[string]progress.Model),
		activeDownloads: make(map[string]DownloadProgressInfo),
		// Zero-valued fields: loadingMsg, selectedItem, mangaInfo, chapters, err, lastQuery, lastMangaID, lastMangaTitle
	}
	return m
}

// Init is the first command that runs when the application starts.
func Init() {
	// Create text input
	ti := textinput.New()
	ti.Placeholder = "Enter manga title..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40
	ti.Prompt = "❯ "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(secondaryColor)
	ti.TextStyle = lipgloss.NewStyle().Foreground(textColor)

	// Create menu items
	items := []list.Item{
		Item{title: "Search Manga", desc: "Find manga by title"},
		Item{title: "Exit", desc: "Quit the application"},
	}

	// Setup menu list
	menuList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	menuList.Title = "Menu"
	menuList.Styles.Title = titleStyle.Copy().Width(40)
	menuList.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(secondaryColor)
	menuList.Styles.HelpStyle = helpStyle

	// Setup delegate for manga and chapter lists
	mangaDelegate := list.NewDefaultDelegate()
	mangaDelegate.Styles.SelectedTitle = highlightTextStyle.Copy()
	mangaDelegate.Styles.SelectedDesc = highlightTextStyle.Copy().Faint(true)

	// Create empty manga list
	mangaList := list.New([]list.Item{}, mangaDelegate, 0, 0)
	mangaList.Title = "Search Results"
	mangaList.Styles.Title = titleStyle.Copy().Width(40)
	mangaList.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(secondaryColor)

	// Create empty chapter list
	chapterList := list.New([]list.Item{}, mangaDelegate, 0, 0)
	chapterList.Title = "Chapters"
	chapterList.Styles.Title = titleStyle.Copy().Width(40)

	// Create spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(secondaryColor)

	// Create channel for download progress
	downloadProgChan := make(chan DownloadProgressInfo, 10) // Buffered channel

	// Create initial model
	initialModel := Model{
		keys:            DefaultKeyMap,
		help:            help.New(),
		state:           searchState,
		input:           ti,
		list:            menuList,
		spinner:         s,
		loadingMsg:      "Searching...",
		downloader:      NewDownloader(3, "", "", downloadProgChan),
		progressModels:  make(map[string]progress.Model),
		activeDownloads: make(map[string]DownloadProgressInfo),
		selectedItem:    Item{},
		mangaInfo:       nil,
		chapters:        []Item{},
		err:             nil,
		width:           80,                         // Initial width, will be updated
		height:          24,                         // Initial height, will be updated
		basedApiUrl:     "https://api.mangadex.org", // Default API URL
	}

	// Start the program
	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil { // Changed from p.Start() to p.Run()
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}

// listenForDownloads creates a command that listens for progress messages
// on the provided channel and returns them as tea.Msg to the main update loop.
func listenForDownloads(ch <-chan DownloadProgressInfo) tea.Cmd {
	return func() tea.Msg {
		progressMsg := <-ch
		return progressMsg
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, textinput.Blink, m.spinner.Tick)
	if m.downloader != nil && m.downloader.ProgressChan != nil {
		cmds = append(cmds, listenForDownloads(m.downloader.ProgressChan))
	}
	return tea.Batch(cmds...)
}

func (m Model) searchManga() tea.Cmd {
	// Change app state to loading before starting the search
	m.state = loadingState
	m.loadingMsg = fmt.Sprintf("Searching for '%s'...", m.input.Value())

	return func() tea.Msg {
		query := url.QueryEscape(m.input.Value()) // Use input value instead of undefined m.query
		apiSearch := fmt.Sprintf("%s/manga?title=%s&limit=20", m.basedApiUrl, query)

		// Add debug output
		fmt.Println("DEBUG: Requesting URL:", apiSearch)

		results, err := CustomRequest(apiSearch)
		if err != nil {
			fmt.Println("DEBUG: API request error:", err)
			return searchErrMsg{err: err}
		}

		// Debug - print first part of response
		fmt.Println("DEBUG: API response (first 300 chars):", results[:min(300, len(results))])

		var mangaList MangaList
		ParseMangaSearch(results, &mangaList)

		// Debug - print parsed results
		fmt.Printf("DEBUG: Parsed %d manga items\n", len(mangaList))
		for i, manga := range mangaList {
			if i < 3 { // Just show first 3 for brevity
				titleText := "Unknown"
				if t, ok := manga.Attributes.Title["en"]; ok {
					titleText = t
				}
				fmt.Printf("DEBUG: Manga %d: ID=%s, Title=%s\n", i, manga.ID, titleText)
			}
		}

		// Convert to search results with list items
		var items []list.Item
		for _, manga := range mangaList {
			var title string
			if t, ok := manga.Attributes.Title["en"]; ok {
				title = t
			} else {
				// Fallback to the first available title
				for _, t := range manga.Attributes.Title {
					title = t
					break
				}
			}

			var desc string
			if d, ok := manga.Attributes.Description["en"]; ok {
				desc = d
			} else {
				for _, d := range manga.Attributes.Description {
					desc = d
					break
				}
			}

			// Truncate long descriptions
			if len(desc) > 100 {
				desc = desc[:100] + "..."
			}

			items = append(items, Item{
				id:    manga.ID,
				title: title,
				desc:  desc,
			})
		}

		return searchResultsMsg{
			items:   items,
			results: mangaList,
		}
	}
}

// Helper function for debug output
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// loadChapters loads chapters for the selected manga
func (m Model) loadChapters(mangaId string) tea.Cmd {
	return func() tea.Msg {
		// Use the standard MangaDex v5 feed endpoint
		chapterUrl := fmt.Sprintf("%s/manga/%s/feed?translatedLanguage[]=en&limit=100&order[chapter]=asc", m.basedApiUrl, mangaId)
		fmt.Println("DEBUG: Requesting Chapters URL:", chapterUrl) // Add debug log for URL
		results, err := CustomRequest(chapterUrl)
		if err != nil {
			fmt.Println("DEBUG: Chapter request error:", err) // Add debug log for error
			return ErrorMsg{err: err}
		}

		// Debug response
		fmt.Println("DEBUG: Chapter response (first 300 chars):", results[:min(300, len(results))])

		var mangaChapters MangaChapters
		ParseChapters(results, &mangaChapters) // This should now work with the correct response format
		fmt.Printf("DEBUG: Parsed %d chapters from feed\n", len(mangaChapters.Chapters))
		return ChaptersMsg(mangaChapters)
	}
}

// fetchAtHomeInfoCmd fetches the at-home server details for a chapter.
func (m Model) fetchAtHomeInfoCmd(chapterID, mangaTitle string) tea.Cmd {
	return func() tea.Msg {
		fmt.Printf("DEBUG: Fetching at-home for chapter %s\n", chapterID) // Debug log
		apiEndpoint := fmt.Sprintf("%s/at-home/server/%s", m.basedApiUrl, chapterID)
		atHomeResults, err := CustomRequest(apiEndpoint)
		if err != nil {
			fmt.Printf("DEBUG: Error fetching at-home: %v\n", err)
			return ErrorMsg{fmt.Errorf("failed to get image server: %w", err)}
		}

		var atHome AtHomeResponse
		if err := json.Unmarshal([]byte(atHomeResults), &atHome); err != nil {
			fmt.Printf("DEBUG: Error parsing at-home JSON: %v\n", err)
			return ErrorMsg{fmt.Errorf("failed to parse image server response: %w", err)}
		}

		if len(atHome.Chapter.Data) == 0 {
			fmt.Printf("DEBUG: No images found for chapter %s\n", chapterID)
			return ErrorMsg{fmt.Errorf("no image data found for chapter %s", chapterID)}
		}

		fmt.Printf("DEBUG: Got at-home info for chapter %s, %d images\n", chapterID, len(atHome.Chapter.Data))
		// Send the successful response back as a message
		return atHomeResponseMsg{
			info:         &atHome,
			chapterID:    chapterID,
			mangaTitle:   mangaTitle,
			chapterTitle: mangaTitle,
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	// Ensure the download listener is always running if channel exists
	if m.downloader != nil && m.downloader.ProgressChan != nil {
		cmds = append(cmds, listenForDownloads(m.downloader.ProgressChan))
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update dimensions
		m.width = msg.Width
		m.height = msg.Height
		h, v := listStyle.GetFrameSize() // Use a consistent style for frame size
		listWidth := msg.Width - h
		listHeight := msg.Height - v - 5 // Adjust height for input/help/progress
		m.list.SetSize(listWidth, listHeight)
		m.input.Width = listWidth - 2 // Make input slightly smaller than list
		m.help.Width = msg.Width - h

		// Update progress bar widths
		for id := range m.progressModels {
			prog := m.progressModels[id]
			prog.Width = listWidth - 4 // Adjust width
			m.progressModels[id] = prog
		}
		return m, nil // Return early after resize

	case tea.KeyMsg:
		// Handle help toggle first
		if key.Matches(msg, m.keys.Help) {
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
		// Handle Quit globally
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

		// Check for Enter key on search input
		if m.state == searchState && m.input.Focused() && key.Matches(msg, m.keys.Submit) {
			m.lastQuery = m.input.Value() // Store query
			return m, m.searchManga()     // Trigger search
		}

		// Skip further key handling if input is focused (Bubble Tea handles it)
		if m.input.Focused() && m.state == searchState {
			break // Let input handle other keys like backspace etc.
		}

		// State-specific key handling
		switch m.state {
		case searchState:
			if key.Matches(msg, m.keys.Submit) {
				cmd = m.searchManga()
				return m, cmd
			}
		case mangaDetailsState: // This state might merge with chapterListState
			if key.Matches(msg, m.keys.Back) {
				m.state = searchState
				m.list.Title = fmt.Sprintf("Search Results for \"%s\"", m.lastQuery)
				// Reload previous search results? Need to store them.
				return m, nil
			}
		case chapterListState:
			if key.Matches(msg, m.keys.Select) {
				selectedItem, ok := m.list.SelectedItem().(Item)
				if ok {
					// Check if already downloading
					if _, exists := m.activeDownloads[selectedItem.id]; exists {
						m.err = fmt.Errorf("chapter %s already downloading/downloaded", selectedItem.title)
						return m, nil // Don't start again
					}
					m.selectedItem = selectedItem // Store selected chapter
					// *** Initiate Download Fetch ***
					cmd = m.fetchChapterServerCmd(selectedItem.id)
					return m, cmd
				}
			} else if key.Matches(msg, m.keys.Back) {
				m.state = searchState
				m.err = nil
				m.list.Title = fmt.Sprintf("Search Results for \"%s\"", m.lastQuery)
				// TODO: Restore previous search list items if needed
				return m, nil
			}
		}

	// --- Message Handling ---
	case spinner.TickMsg:
		if m.state == loadingState {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case progress.FrameMsg:
		progressCmds := make([]tea.Cmd, 0, len(m.progressModels))
		for id, pm := range m.progressModels {
			if info, ok := m.activeDownloads[id]; ok && !info.Done {
				newPm, cmd := pm.Update(msg)
				if newModel, ok := newPm.(progress.Model); ok {
					m.progressModels[id] = newModel
					progressCmds = append(progressCmds, cmd)
				}
			}
		}
		return m, tea.Batch(progressCmds...)

	case searchErrMsg:
		m.state = errorState
		m.err = msg
		return m, nil // Stop loading etc.

	case searchResultsMsg:
		m.state = searchState // Go back to search state to display results
		m.list.Title = fmt.Sprintf("Search Results for \"%s\" (%d)", m.lastQuery, len(msg.items))
		m.list.SetItems(msg.items)
		// Store raw results if needed later? m.mangaResults = msg.results
		return m, nil

	case chaptersMsg:
		m.state = chapterListState
		m.mangaInfo = msg.manga   // Store manga details
		m.chapters = msg.chapters // Store parsed chapters
		m.list.Title = fmt.Sprintf("Chapters for %s (%d)", m.lastMangaTitle, len(msg.chapters))

		var listItems []list.Item
		for _, item := range msg.chapters {
			listItems = append(listItems, item)
		}
		m.list.SetItems(listItems) // Update list with chapters
		return m, nil

	case atHomeResponseMsg:
		// Received server info, start download
		m.state = chapterListState // Stay in chapter list view
		m.err = nil                // Clear any previous error
		chapterID := msg.info.Chapter.Hash
		chapterTitle := m.selectedItem.title // Get title from stored selected item

		// Create and store the progress bar
		prog := progress.New(
			progress.WithDefaultGradient(),
			progress.WithWidth(m.width-padding*2-4),
			progress.WithoutPercentage(),
		)
		m.progressModels[chapterID] = prog

		// Add to active downloads list
		m.activeDownloads[chapterID] = DownloadProgressInfo{
			ChapterID: chapterID,
			Total:     len(msg.info.Chapter.Data), // Set total pages
			Done:      false,
			Error:     nil,
		}

		// Start the download command
		cmd = m.downloadChapterCmd(msg.info, chapterTitle)
		cmds = append(cmds, cmd)
		// Also start the progress bar animation
		cmds = append(cmds, func() tea.Msg { return progress.FrameMsg{} })

	case DownloadProgressInfo:
		// Update active download info
		if _, ok := m.activeDownloads[msg.ChapterID]; ok {
			m.activeDownloads[msg.ChapterID] = msg
		} else {
			fmt.Printf("WARN: Progress received for untracked chapter %s\n", msg.ChapterID)
			return m, tea.Batch(cmds...)
		}

		// Update progress bar model
		if prog, ok := m.progressModels[msg.ChapterID]; ok {
			var progressCmd tea.Cmd
			if msg.Done {
				// Final update on completion/error
				progressCmd = prog.SetPercent(1.0)
				// TODO: Maybe remove from maps after a short delay?
				// For now, keep it showing final state.
			} else {
				progressCmd = prog.SetPercent(float64(msg.Completed) / float64(msg.Total))
			}
			m.progressModels[msg.ChapterID] = prog
			cmds = append(cmds, progressCmd)
		}
	}

	// --- Update focused component (List or Input) ---
	if m.state == searchState {
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.state == mangaDetailsState || m.state == chapterListState {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update help
	m.help.ShowAll = false // Always show short help unless toggled
	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.err != nil {
		// Simple error view
		return fmt.Sprintf("\nError: %v\n\nPress 'q' to quit.", m.err)
	}

	var viewContent string

	switch m.state {
	case searchState:
		viewContent = lipgloss.JoinVertical(lipgloss.Left,
			m.input.View(),
			listStyle.Render(m.list.View()), // Show search results list
		)
	case loadingState:
		viewContent = loadingStyle.Render(fmt.Sprintf("%s %s", m.spinner.View(), m.loadingMsg))
	case mangaDetailsState:
		// TODO: Display manga details nicely
		viewContent = "Manga Details View (Not Implemented Yet)"
	case chapterListState:
		viewContent = listStyle.Render(m.list.View()) // Show chapter list
	case errorState:
		viewContent = errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	default:
		viewContent = "Unknown state"
	}

	// --- Download Progress View ---
	progressView := ""
	if len(m.activeDownloads) > 0 {
		var progressBars []string
		progressBars = append(progressBars, downloadHeaderStyle.Render("Downloads:")) // Header

		// Sort IDs for consistent display order?
		var sortedIDs []string
		for id := range m.activeDownloads {
			sortedIDs = append(sortedIDs, id)
		}
		// Sort logic here if needed (e.g., alphabetically or by start time)

		for _, chapterID := range sortedIDs {
			info := m.activeDownloads[chapterID]
			progModel, exists := m.progressModels[chapterID]
			title := "Chapter " + chapterID
			if title == "" {
				title = "Chapter " + chapterID
			} // Fallback title
			if exists {
				bar := progModel.View()
				status := ""
				if info.Done {
					if info.Error != nil {
						status = errorStyle.Render(fmt.Sprintf(" ❌ Error: %v", info.Error))
					} else {
						status = successStyle.Render(" ✅ Done")
					}
				} else if info.Total > 0 { // Show progress only if total is known
					status = fmt.Sprintf(" (%d/%d)", info.Completed, info.Total)
				} else {
					status = " (Starting...)" // Initial state
				}
				progressBars = append(progressBars, fmt.Sprintf("%s%s\n%s", title, status, bar))
			} else {
				progressBars = append(progressBars, fmt.Sprintf("%s: (Initializing...)", title))
			}
		}
		progressView = lipgloss.JoinVertical(lipgloss.Left, progressBars...)
		progressView = progressAreaStyle.Render(progressView) // Add padding/border
	}

	// Combine main content, progress, and help
	finalContent := lipgloss.JoinVertical(lipgloss.Left,
		viewContent,
		progressView,                          // Add progress view below main content
		helpStyle.Render(m.help.View(m.keys)), // Show help at the bottom
	)

	// Add padding around the entire view
	return lipgloss.NewStyle().Padding(1, 2).Render(finalContent)
}

func EndMessage() {
	fmt.Println(infoStyle.Render("\nManga-g has completed.\nStart program again to search for another manga."))
}

// QueryCheck validates the user's query input
func QueryCheck(query string) error {
	if strings.TrimSpace(query) == "" {
		return fmt.Errorf("query cannot be empty, please enter a valid search term")
	}
	return nil
}

// getFirstValue returns the first value from a map
func getFirstValue(m map[string]string) string {
	for _, value := range m {
		return value
	}
	return ""
}

// fetchChapterServer sends a command to get the Manga@Home server URL for a chapter.
func (m *Model) fetchChapterServer(chapterID string) tea.Cmd {
	// This is a placeholder - implementation needed
	return nil
}

// downloadChapterCmd sends a command to start the actual chapter download.
// Assumes DownloadChapter is now part of the Downloader struct and sends progress via its channel.
func (m *Model) downloadChapterCmd(atHomeInfo *AtHomeResponse, chapterTitle string) tea.Cmd {
	mangaTitle := m.lastMangaTitle // Use stored manga title
	chapterID := atHomeInfo.Chapter.Hash

	return func() tea.Msg {
		fmt.Printf("DEBUG: Starting download goroutine for Chapter %s (ID: %s)\n", chapterTitle, chapterID)
		// Call the DownloadChapter method with the correct number of arguments
		err := m.downloader.DownloadChapter(atHomeInfo, mangaTitle, chapterID)
		if err != nil {
			fmt.Printf("DEBUG: Downloader finished for %s with error: %v\n", chapterID, err)
			// Downloader sends final error via channel, no explicit msg needed here
		} else {
			fmt.Printf("DEBUG: Downloader finished successfully for %s.\n", chapterID)
		}
		return nil // No message needed, rely on progress channel
	}
}

// fetchChapterServerCmd sends a command to get the Manga@Home server URL for a chapter.
func (m *Model) fetchChapterServerCmd(chapterID string) tea.Cmd {
	m.state = loadingState
	m.loadingMsg = fmt.Sprintf("Fetching download info for %s...", m.selectedItem.title)

	return func() tea.Msg {
		fmt.Printf("DEBUG: Fetching at-home for chapter %s\n", chapterID)
		// Replace with direct implementation
		fmt.Printf("DEBUG: Placeholder implementation for at-home server fetch\n")
		return searchErrMsg{fmt.Errorf("not implemented: fetch chapter server")}
	}
}

// Add a helper for truncating filenames
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return "..." // Avoid panic if maxLen is too small
	}
	return s[:maxLen-3] + "..."
}

// --- Helper Functions ---
// Helper to create a command for searching manga
func (m *Model) searchMangaCmd() tea.Cmd {
	queryString := m.input.Value() // Use current input value
	if queryString == "" {
		return func() tea.Msg { return searchErrMsg{fmt.Errorf("search query cannot be empty")} }
	}
	m.state = loadingState
	m.loadingMsg = fmt.Sprintf("Searching for '%s'...", queryString)
	m.list.SetItems([]list.Item{}) // Clear previous results
	m.lastQuery = queryString      // Store the query for later use (e.g., list title)

	return func() tea.Msg {
		// Replace with direct API call implementation
		// Placeholder implementation
		return searchResultsMsg{items: []list.Item{}, results: nil}
	}
}

// Helper to create a command for fetching manga details (including chapters)
func (m *Model) fetchMangaDetailsCmd(mangaID, mangaTitle string) tea.Cmd {
	m.state = loadingState
	m.loadingMsg = fmt.Sprintf("Loading chapters for %s...", mangaTitle)
	m.list.SetItems([]list.Item{}) // Clear manga list
	m.lastMangaID = mangaID        // Store manga context
	m.lastMangaTitle = mangaTitle

	return func() tea.Msg {
		// Replace with direct API call implementation
		// Placeholder implementation
		return chaptersMsg{chapters: []Item{}}
	}
}
