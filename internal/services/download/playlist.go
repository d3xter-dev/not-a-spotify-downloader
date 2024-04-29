package download

import (
	"fmt"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/core"
	Spotify "github.com/d3xter-dev/not-a-spotify-downloader/internal/spotify"
)

type PlaylistStrategy struct {
	session  *core.Session
	playlist *Spotify.SelectedListContent
}

func NewPlaylistStrategy(session *core.Session, playlist *Spotify.SelectedListContent) *PlaylistStrategy {
	return &PlaylistStrategy{
		session:  session,
		playlist: playlist,
	}
}

func (s *PlaylistStrategy) GetSaveDir() string {
	return fmt.Sprintf("./playlist/%s", sanitizeFilename(s.playlist.GetAttributes().GetName()))
}

func (s *PlaylistStrategy) GetItems() []*Spotify.Item {
	return s.playlist.GetContents().GetItems()
}
