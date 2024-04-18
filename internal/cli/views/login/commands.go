package login

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/core"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/services/auth"
)

type SuccessMsg struct {
	Session *core.Session
}

type FailedMsg struct {
	Error error
}

type StatusMsg struct {
	Status string
}

func StartLogin(username, password string, remember bool) tea.Cmd {
	login := func() tea.Msg {
		session, err := auth.LoginWithUser(username, password, remember)
		if err != nil {
			return FailedMsg{err}
		}
		return SuccessMsg{session}
	}

	return tea.Batch(SendLoginStatus("Trying to login..."), login)
}

func CheckExistingLogin() tea.Cmd {
	checkExistingLogin := func() tea.Msg {
		session, err := auth.LoginWithBlob()
		if err != nil {
			return FailedMsg{err}
		}
		return SuccessMsg{session}
	}

	return tea.Batch(SendLoginStatus("Trying to login with existing account..."), checkExistingLogin)
}

func SendLoginStatus(status string) tea.Cmd {
	return func() tea.Msg {
		return StatusMsg{status}
	}
}
