package librespot

import (
	core "github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/core"
)

// Login to Spotify using username and password
func Login(username string, password string, deviceName string) (*core.Session, error) {
	return core.Login(username, password, deviceName)
}

// Login to Spotify using an existing authData blob
func LoginSaved(username string, authData []byte, deviceName string) (*core.Session, error) {
	return core.LoginSaved(username, authData, deviceName)
}

// Login to Spotify using the OAuth method
func LoginOAuth(deviceName string, clientId string, clientSecret string) (*core.Session, error) {
	return core.LoginOAuth(deviceName, clientId, clientSecret)
}
