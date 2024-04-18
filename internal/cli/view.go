package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/cli/theme"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/cli/views/download"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/cli/views/login"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/cli/views/playlist"
)

type View interface {
	GetName() string
	Init() tea.Cmd
	Update(msg tea.Msg) tea.Cmd
	Render() string
	SetSize(width, height int)
}

const (
	DownloadView = "download"
	LoginView    = "login"
	PlaylistView = "playlist"
)

func GetViews() map[string]View {
	defaultTheme := theme.DefaultTheme{}

	return map[string]View{
		DownloadView: download.NewView(defaultTheme),
		LoginView:    login.NewView(defaultTheme),
		PlaylistView: playlist.NewView(defaultTheme),
	}
}
