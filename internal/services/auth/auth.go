package auth

import (
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/core"
	"os"
	"strings"
)

type AuthBlob struct {
	Username string
	Device   string
	Blob     []byte
}

func (b *AuthBlob) ToFile() error {
	blobString := b.Username + ":" + b.Device + ":" + string(b.Blob)
	err := os.WriteFile("auth_blob", []byte(blobString), 0644)
	return err
}

func GetAuthBlob() (*AuthBlob, error) {
	file, err := os.ReadFile("auth_blob")
	if err != nil {
		return nil, err
	}

	authData := strings.Split(string(file), ":")

	return &AuthBlob{
		Username: authData[0],
		Device:   authData[1],
		Blob:     []byte(authData[2]),
	}, nil
}

func LoginWithBlob() (*core.Session, error) {
	blob, err := GetAuthBlob()
	if err != nil {
		return nil, err
	}

	sess, err := librespot.LoginSaved(blob.Username, blob.Blob, blob.Device)
	return sess, err
}

func LoginWithUser(username string, password string, remember bool) (*core.Session, error) {
	sess, err := librespot.Login(username, password, "deviceName")
	if err != nil {
		return nil, err
	}

	if remember {
		blob := &AuthBlob{
			Username: username,
			Device:   "deviceName",
			Blob:     sess.ReusableAuthBlob(),
		}

		err := blob.ToFile()
		if err != nil {
			return nil, err
		}
	}

	return sess, nil
}
