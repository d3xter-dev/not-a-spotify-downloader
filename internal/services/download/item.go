package download

import "time"

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
