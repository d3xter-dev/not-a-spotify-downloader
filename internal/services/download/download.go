package download

import (
	"fmt"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/core"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/utils"
	Spotify "github.com/d3xter-dev/not-a-spotify-downloader/internal/spotify"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

type Status string

const (
	StatusDownloading Status = "Downloading"
	StatusComplete    Status = "Complete"
	StatusFailed      Status = "Failed"
)

type Item struct {
	Name      string
	Path      string
	StartTime time.Time
	Status    Status
}

func mapSpotTrackToDownloadItem(item *Spotify.Track, status Status) Item {
	return Item{
		Name:      item.GetName(),
		Path:      utils.ConvertTo62(item.GetGid()),
		StartTime: time.Now(),
		Status:    status,
	}
}

var extMap = map[Spotify.AudioFile_Format]string{
	Spotify.AudioFile_OGG_VORBIS_96:  ".96.ogg",
	Spotify.AudioFile_OGG_VORBIS_160: ".160.ogg",
	Spotify.AudioFile_OGG_VORBIS_320: ".320.ogg",
	Spotify.AudioFile_MP3_256:        ".256.mp3",
	Spotify.AudioFile_MP3_320:        ".320.mp3",
	Spotify.AudioFile_MP3_160:        ".160.mp3",
	Spotify.AudioFile_MP3_96:         ".96.mp3",
	Spotify.AudioFile_MP3_160_ENC:    ".160enc.mp3",
	Spotify.AudioFile_AAC_24:         ".24.aac",
	Spotify.AudioFile_AAC_48:         ".48.aac",
}

func StartProcessing(session *core.Session, playlist *Spotify.SelectedListContent) (chan Item, error) {
	ch := make(chan Item)

	playlistName := playlist.GetAttributes().GetName()
	saveDir := fmt.Sprintf("./playlist/%s", playlistName)
	err := os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	go processPlaylist(session, playlist, ch, saveDir)

	return ch, nil
}

func getBestFileFormat(track *Spotify.Track) *Spotify.AudioFile {
	bestFormats := []Spotify.AudioFile_Format{
		Spotify.AudioFile_OGG_VORBIS_320,
		Spotify.AudioFile_MP3_320,
		Spotify.AudioFile_MP3_256,
		Spotify.AudioFile_OGG_VORBIS_160,
		Spotify.AudioFile_MP3_160,
	}

	for _, format := range bestFormats {
		for _, file := range track.File {
			if file.GetFormat() == format {
				return file
			}
		}

		for _, alt := range track.Alternative {
			for _, file := range alt.File {
				if file.GetFormat() == format {
					return file
				}
			}
		}
	}

	return nil
}

type processingState struct {
	processed  int
	running    int
	maxRunning int
	mutex      sync.Mutex
	waitGroup  sync.WaitGroup
}

func (s *processingState) getRunning() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.running
}

func (s *processingState) addRunning() {
	s.mutex.Lock()
	s.running++
	s.mutex.Unlock()
}

func (s *processingState) stopRunning() {
	s.mutex.Lock()
	s.running--
	s.mutex.Unlock()
}

func processPlaylist(session *core.Session, playlist *Spotify.SelectedListContent, ch chan Item, playlistDir string) {
	state := processingState{
		processed:  0,
		running:    0,
		maxRunning: 5,
	}

	playlistItems := playlist.GetContents().GetItems()

	for state.processed < len(playlistItems) {
		if state.getRunning() < state.maxRunning {
			state.waitGroup.Add(1)
			state.addRunning()

			trackId := utils.ParseIdFromString(playlistItems[state.processed].GetUri())
			track, err := session.Mercury().GetTrack(utils.Base62ToHex(trackId))
			if err != nil {
				return
			}

			go func(item *Spotify.Track, session *core.Session, playlistDir string, ch chan Item) {
				defer func() {
					state.waitGroup.Done()
					state.stopRunning()
				}()

				ch <- mapSpotTrackToDownloadItem(item, StatusDownloading)

				bestFile := getBestFileFormat(item)
				if bestFile == nil {
					ch <- mapSpotTrackToDownloadItem(item, StatusFailed)
					return
				}

				file, err := session.Player().LoadTrack(bestFile, item.GetGid())
				if err != nil {
					ch <- mapSpotTrackToDownloadItem(item, StatusFailed)
					return
				}

				buffer, err := io.ReadAll(file)
				if err != nil {
					ch <- mapSpotTrackToDownloadItem(item, StatusFailed)
					return
				}

				ext := extMap[bestFile.GetFormat()]
				err = os.WriteFile(path.Join(playlistDir, item.GetName()+ext), buffer, os.ModePerm)
				if err != nil {
					ch <- mapSpotTrackToDownloadItem(item, StatusFailed)
					return
				}

				ch <- mapSpotTrackToDownloadItem(item, StatusComplete)
			}(track, session, playlistDir, ch)

			state.processed++
		}
		time.Sleep(1000 * time.Millisecond)
	}

	state.waitGroup.Wait()
}
