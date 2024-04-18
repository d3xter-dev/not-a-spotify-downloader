package login

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/cli/theme"
	"strings"
)

type View struct {
	focusIndex   int
	status       string
	stayLoggedIn bool
	inputs       []textinput.Model

	spinner spinner.Model

	width  int
	height int
	theme  theme.Theme
}

func NewView(theme theme.Theme) *View {
	usernameInput := textinput.New()
	usernameInput.Placeholder = "Username"
	usernameInput.Focus()

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Password"
	passwordInput.EchoMode = textinput.EchoPassword

	return &View{
		theme:        theme,
		focusIndex:   0,
		stayLoggedIn: true,
		spinner:      spinner.New(spinner.WithSpinner(spinner.Line)),
		inputs: []textinput.Model{
			usernameInput, passwordInput,
		},
	}
}

func (v *View) GetName() string {
	return "Login"
}

func (v *View) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, CheckExistingLogin(), v.spinner.Tick)
}

func (v *View) SetSize(width, height int) {
	v.width, v.height = width, height
}

func (v *View) Update(message tea.Msg) tea.Cmd {

	switch msg := message.(type) {
	case FailedMsg:
		v.status = ""
	case StatusMsg:
		v.status = msg.Status
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			v.nextFocus()
		case "up":
			v.backFocus()
		case "enter":
			if v.currentFocusElement() == "check-login" {
				v.stayLoggedIn = !v.stayLoggedIn
				return nil
			}

			if v.currentFocusElement() == "btn-cancel" {
				return tea.Quit
			}

			username, password := v.inputs[0].Value(), v.inputs[1].Value()
			if len(username) < 5 || len(password) < 3 {
				return nil
			}

			return StartLogin(username, password, v.stayLoggedIn)
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		v.spinner, cmd = v.spinner.Update(msg)
		return cmd
	}

	commands := make([]tea.Cmd, len(v.inputs))

	for i := range v.inputs {
		v.inputs[i], commands[i] = v.inputs[i].Update(message)
	}

	return tea.Batch(commands...)
}

func (v *View) Render() string {
	var ui = ""
	if v.status == "" {
		ui = v.renderLoginForm()
	} else {
		ui = lipgloss.NewStyle().Width(75).Align(lipgloss.Center).Render(lipgloss.JoinHorizontal(lipgloss.Center, v.spinner.View(), " ", v.status))
	}

	dialog := lipgloss.Place(v.width, v.height,
		lipgloss.Center, lipgloss.Center,
		v.theme.DialogBox().Render(ui),
		lipgloss.WithWhitespaceChars("#"),
		lipgloss.WithWhitespaceForeground(v.theme.Subtle()),
	)

	return dialog
}

func (v *View) renderLoginForm() string {

	question := lipgloss.NewStyle().Width(75).Align(lipgloss.Center).Render("Login with your Spotify account")

	var inputs = ""
	for _, input := range v.inputs {
		input.PromptStyle = v.theme.Input()
		input.TextStyle = v.theme.Input()
		input.Width = 27

		if input.Focused() {
			input.TextStyle = v.theme.InputFocused()
			input.PromptStyle = v.theme.InputFocused()
		}

		inputs = lipgloss.JoinVertical(lipgloss.Center, inputs, input.View())
	}

	checkboxValue := "[x]"
	checkboxStyle := v.theme.Checkbox()

	if !v.stayLoggedIn {
		checkboxValue = "[ ]"
	}
	if v.currentFocusElement() == "check-login" {
		checkboxStyle = v.theme.CheckboxFocused()
	}

	checkbox := checkboxStyle.MarginTop(1).Width(30).Render(checkboxValue + " Stay logged in")

	cancelStyle, okStyle := v.theme.Button(), v.theme.ActiveButton()
	if v.currentFocusElement() == "btn-cancel" {
		cancelStyle = v.theme.ActiveButton()
		okStyle = v.theme.Button()
	}
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okStyle.Render("Login"), "  ", cancelStyle.Render("Cancel"))

	return lipgloss.JoinVertical(lipgloss.Center, question, inputs, checkbox, buttons)
}

var focusableElements = []string{"input-username", "input-password", "check-login", "btn-login", "btn-cancel"}

func (v *View) currentFocusElement() string {
	maxIndex := len(focusableElements) - 1
	if v.focusIndex > maxIndex || v.focusIndex < 0 {
		return ""
	}
	return focusableElements[v.focusIndex]
}

func (v *View) nextFocus() {
	v.focusIndex++
	v.updateFocus()
}

func (v *View) backFocus() {
	v.focusIndex++
	v.updateFocus()
}

func (v *View) updateFocus() {
	maxIndex := len(focusableElements) - 1
	if v.focusIndex < 0 {
		v.focusIndex = maxIndex
		v.updateFocus()
		return
	}

	if v.focusIndex > maxIndex {
		v.focusIndex = 0
		v.updateFocus()
		return
	}

	for i := range v.inputs {
		v.inputs[i].Blur()
	}

	if strings.HasPrefix(v.currentFocusElement(), "input-") {
		v.inputs[v.focusIndex].Focus()
	}
}
