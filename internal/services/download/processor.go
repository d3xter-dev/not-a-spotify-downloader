package download

import (
	"context"
	"errors"
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

	go startWorkers(p.strategy, p.session, p.returnChannel)

	return p.returnChannel, nil
}

func startWorkers(strategy Strategy, session *core.Session, returnChannel chan Item) {
	state := &State{
		isPaused: false,
	}

	playlistItems := strategy.GetItems()
	for _, item := range playlistItems {
		state.addToQueue(&QueueItem{
			item: item,
		})
	}

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

	for i := 0; i < 5; i++ {
		go func() {
			for {
				if state.getIsPaused() {
					time.Sleep(150 * time.Millisecond)
					continue
				}

				item, err := state.nextItem()
				if err != nil {
					break
				}

				err = processItem(processOptions{
					session:       session,
					item:          item.item,
					returnChannel: returnChannel,
					saveDir:       strategy.GetSaveDir(),
				})

				if err != nil {
					state.retryItem(item)
				}
			}
		}()
	}
}

type processOptions struct {
	session *core.Session
	item    *Spotify.Item

	returnChannel chan Item
	saveDir       string
}

func processItem(options processOptions) error {
	session := options.session
	item := options.item
	ctx, cancel := context.WithTimeout(options.session.Context(), 60*time.Second)

	trackId := utils.ParseIdFromString(item.GetUri())
	track, err := session.Mercury().GetTrack(utils.Base62ToHex(trackId))

	defer func() {
		cancel()
	}()

	if err != nil {
		return err
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	bestFile := getBestFileFormat(track)
	if bestFile == nil {
		return errors.New("no file found")
	}

	var SendStatusUpdate = func(status Status) {
		options.returnChannel <- mapSpotTrackToDownloadItem(track, status)
	}

	SendStatusUpdate(StatusDownloading)

	file, err := session.Player().LoadTrack(bestFile, track.GetGid())
	if err != nil {
		SendStatusUpdate(StatusFailed)
		return err
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
			return err
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
		return err
	}

	SendStatusUpdate(StatusComplete)
	return nil
}
