package playlist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/core"
	Spotify "github.com/d3xter-dev/not-a-spotify-downloader/internal/spotify"
)

type setSessionMsg struct {
	session *core.Session
}

type FetchSuccessMsg struct {
	Session  *core.Session
	Playlist *Spotify.SelectedListContent
}

type FetchFailedMsg struct {
	Error error
}

func SetSession(session *core.Session) tea.Cmd {
	return func() tea.Msg {
		return setSessionMsg{session: session}
	}
}

func FetchPlaylist(session *core.Session, id string) tea.Cmd {
	return func() tea.Msg {
		playlist, err := session.Mercury().GetPlaylist(id)
		if err != nil || playlist.GetAttributes() == nil {
			return FetchFailedMsg{err}
		}
		return FetchSuccessMsg{session, playlist}
	}
}
