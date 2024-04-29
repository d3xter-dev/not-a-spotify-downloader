package download

import (
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/utils"
	Spotify "github.com/d3xter-dev/not-a-spotify-downloader/internal/spotify"
	"regexp"
	"time"
)

func mapSpotTrackToDownloadItem(item *Spotify.Track, status Status) Item {
	return Item{
		Name:      item.GetName(),
		Path:      utils.ConvertTo62(item.GetGid()),
		StartTime: time.Now(),
		Status:    status,
	}
}

var extensionMap = map[Spotify.AudioFile_Format]string{
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

func getBestFileFormat(track *Spotify.Track) *Spotify.AudioFile {
	bestFormats := []Spotify.AudioFile_Format{
		Spotify.AudioFile_MP3_320,
		Spotify.AudioFile_OGG_VORBIS_320,
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

func sanitizeFilename(filename string) string {
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	sanitized := invalidChars.ReplaceAllString(filename, "")

	if len(sanitized) == 0 {
		sanitized = "unnamed"
	}

	return sanitized
}
