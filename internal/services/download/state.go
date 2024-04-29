package download

import "sync"

type State struct {
	processed  int
	running    int
	maxRunning int
	isPaused   bool
	mutex      sync.Mutex
	pauseMutex sync.Mutex
}

func (s *State) setPause(val bool) {
	s.pauseMutex.Lock()
	s.isPaused = val
	defer s.pauseMutex.Unlock()
}

func (s *State) getIsPaused() bool {
	s.pauseMutex.Lock()
	defer s.pauseMutex.Unlock()
	return s.isPaused
}

func (s *State) getRunning() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.running
}

func (s *State) addRunning() {
	s.mutex.Lock()
	s.running++
	s.mutex.Unlock()
}

func (s *State) stopRunning() {
	s.mutex.Lock()
	s.running--
	s.mutex.Unlock()
}
