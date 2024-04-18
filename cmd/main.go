package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/cli"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/cli/views/download"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/cli/views/login"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/cli/views/playlist"
	"log"
)

type Model struct {
	currentPage string
	pages       map[string]cli.View
}

func InitialModel() Model {
	return Model{
		currentPage: cli.LoginView,
		pages:       cli.GetViews(),
	}
}

func (m Model) GetView() cli.View {
	return m.pages[m.currentPage]
}

func (m Model) Init() tea.Cmd {
	return m.GetView().Init()
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		for _, page := range m.pages {
			page.SetSize(msg.Width, msg.Height)
		}

	case login.SuccessMsg:
		m.currentPage = cli.PlaylistView
		return m, tea.Batch(m.GetView().Init(), playlist.SetSession(msg.Session))

	case playlist.FetchSuccessMsg:
		m.currentPage = cli.DownloadView
		return m, tea.Batch(m.GetView().Init(), download.SetPlaylist(msg.Session, msg.Playlist))

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, m.GetView().Update(message)
}

func (m Model) View() string {
	return m.GetView().Render()
}

func main() {
	model := InitialModel()

	program := tea.NewProgram(model)
	program.SetWindowTitle("Not a Spotify Downloader")
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
