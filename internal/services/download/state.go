package download

import (
	"errors"
	Spotify "github.com/d3xter-dev/not-a-spotify-downloader/internal/spotify"
	"sync"
)

type State struct {
	queue      []*QueueItem
	isPaused   bool
	mutex      sync.Mutex
	pauseMutex sync.Mutex
}

type QueueItem struct {
	retries int
	item    *Spotify.Item
}

func (s *State) setPause(val bool) {
	s.pauseMutex.Lock()
	s.isPaused = val
	s.pauseMutex.Unlock()
}

func (s *State) getIsPaused() bool {
	s.pauseMutex.Lock()
	defer s.pauseMutex.Unlock()
	return s.isPaused
}

func (s *State) addToQueue(item *QueueItem) {
	s.mutex.Lock()
	s.queue = append(s.queue, item)
	s.mutex.Unlock()
}

func (s *State) nextItem() (QueueItem, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.queue) == 0 {
		return QueueItem{}, errors.New("queue is empty")
	}

	var nextItem *QueueItem
	nextItem, s.queue = s.queue[0], s.queue[1:]

	return *nextItem, nil
}

func (s *State) retryItem(item QueueItem) {
	item.retries++
	if item.retries <= 3 {
		s.addToQueue(&item)
	}
}
