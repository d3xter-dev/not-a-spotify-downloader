package download

import (
	"context"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/core"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/utils"
	Spotify "github.com/d3xter-dev/not-a-spotify-downloader/internal/spotify"
	"io"
	"os"
	"path"
	"time"
)

type Strategy interface {
	GetSaveDir() string
	GetItems() []*Spotify.Item
}

type Processor struct {
	session       *core.Session
	strategy      Strategy
	returnChannel chan Item
}

func NewProcessor(session *core.Session, strategy Strategy) *Processor {
	return &Processor{session: session, strategy: strategy}
}

func (p *Processor) StartProcessing() (chan Item, error) {
	p.returnChannel = make(chan Item)

	saveDir := p.strategy.GetSaveDir()
	err := os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	go processTracks(p.strategy, p.session, p.returnChannel)

	return p.returnChannel, nil
}

func processTracks(strategy Strategy, session *core.Session, returnChannel chan Item) {
	state := &State{
		processed:  0,
		running:    0,
		maxRunning: 5,
		isPaused:   false,
	}

	playlistItems := strategy.GetItems()

	go func() {
		for {
			select {
			case <-session.Context().Done():
				state.setPause(true)
				time.Sleep(15 * time.Second)
				state.setPause(false)
			}
		}
	}()

	for state.processed < len(playlistItems) {
		time.Sleep(100 * time.Millisecond)

		if state.getIsPaused() {
			continue
		}

		if state.getRunning() < state.maxRunning {
			state.addRunning()

			// check again because state might have changed
			if state.getIsPaused() {
				state.stopRunning()
				continue
			}

			trackId := utils.ParseIdFromString(playlistItems[state.processed].GetUri())
			track, err := session.Mercury().GetTrack(utils.Base62ToHex(trackId))
			if err != nil {
				continue
			}

			go processTrack(processOptions{
				session:       session,
				track:         track,
				state:         state,
				returnChannel: returnChannel,
				saveDir:       strategy.GetSaveDir(),
			})

			state.processed++
		}
	}
}

type processOptions struct {
	session *core.Session
	track   *Spotify.Track
	state   *State

	returnChannel chan Item
	saveDir       string
}

func processTrack(options processOptions) {
	track := options.track
	session := options.session
	state := options.state
	ctx, cancel := context.WithTimeout(options.session.Context(), 60*time.Second)

	defer func() {
		state.stopRunning()
		cancel()
	}()

	var SendStatusUpdate = func(status Status) {
		options.returnChannel <- mapSpotTrackToDownloadItem(track, status)
	}

	SendStatusUpdate(StatusDownloading)

	if ctx.Err() != nil {
		SendStatusUpdate(StatusFailed)
		return
	}

	bestFile := getBestFileFormat(track)
	if bestFile == nil {
		SendStatusUpdate(StatusFailed)
		return
	}

	file, err := session.Player().LoadTrack(bestFile, track.GetGid())
	if err != nil {
		SendStatusUpdate(StatusFailed)
		return
	}

	buffer := make([]byte, 0, 512)
	for {
		n, err := file.Read(ctx, buffer[len(buffer):cap(buffer)])
		buffer = buffer[:len(buffer)+n]
		if err != nil {
			if err == io.EOF {
				//read until end
				break
			}

			SendStatusUpdate(StatusFailed)
			return
		}

		if len(buffer) == cap(buffer) {
			// Add more capacity (let append pick how much).
			buffer = append(buffer, 0)[:len(buffer)]
		}
	}

	ext := extensionMap[bestFile.GetFormat()]
	filename := sanitizeFilename(track.GetName()) + ext
	out := path.Join(options.saveDir, filename)

	err = os.WriteFile(out, buffer, os.ModePerm)
	if err != nil {
		SendStatusUpdate(StatusFailed)
		return
	}

	SendStatusUpdate(StatusComplete)
}
