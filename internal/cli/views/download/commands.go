package download

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/core"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/services/download"
	Spotify "github.com/d3xter-dev/not-a-spotify-downloader/internal/spotify"
)

type setPlaylistMsg struct {
	session  *core.Session
	playlist *Spotify.SelectedListContent
}

func SetPlaylist(session *core.Session, playlist *Spotify.SelectedListContent) tea.Cmd {
	return func() tea.Msg {
		return setPlaylistMsg{session, playlist}
	}
}

func WatchDownload(sub chan download.Item) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}
