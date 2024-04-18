package playlist

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/cli/theme"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/core"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/utils"
)

type View struct {
	focusIndex    int
	playlistInput textinput.Model
	session       *core.Session

	spinner spinner.Model

	width  int
	height int
	theme  theme.Theme
	title  string
}

func NewView(theme theme.Theme) *View {
	playlistInput := textinput.New()
	playlistInput.Placeholder = "PlayList URL or ID"
	playlistInput.Focus()

	return &View{
		title:         "Which playlist do you want to download?",
		theme:         theme,
		focusIndex:    0,
		spinner:       spinner.New(spinner.WithSpinner(spinner.Line)),
		playlistInput: playlistInput,
	}
}

func (v *View) GetName() string {
	return "Find PlayList"
}

func (v *View) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, v.spinner.Tick)
}

func (v *View) SetSize(width, height int) {
	v.width, v.height = width, height
}

func (v *View) Update(message tea.Msg) tea.Cmd {

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			id := utils.ParseIdFromString(v.playlistInput.Value())
			if id == "" {
				return nil
			}
			return FetchPlaylist(v.session, id)
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		v.spinner, cmd = v.spinner.Update(msg)
		return cmd
	case FetchFailedMsg:
		v.title = "Could not fetch playlist, try again:"
	case setSessionMsg:
		v.session = msg.session
	}

	var command tea.Cmd
	v.playlistInput, command = v.playlistInput.Update(message)

	return command
}

func (v *View) Render() string {
	question := lipgloss.NewStyle().Width(75).Align(lipgloss.Center).Render(v.title)

	input := v.playlistInput
	input.PromptStyle = v.theme.Input()
	input.TextStyle = v.theme.Input()

	if input.Focused() {
		input.TextStyle = v.theme.InputFocused()
		input.PromptStyle = v.theme.InputFocused()
	}

	button := v.theme.ActiveButton().Render("Fetch")

	ui := lipgloss.JoinVertical(lipgloss.Center, question, input.View(), button)

	dialog := lipgloss.Place(v.width, v.height,
		lipgloss.Center, lipgloss.Center,
		v.theme.DialogBox().Render(ui),
		lipgloss.WithWhitespaceChars("#"),
		lipgloss.WithWhitespaceForeground(v.theme.Subtle()),
	)

	return dialog
}
