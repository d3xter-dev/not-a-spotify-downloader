package theme

import "github.com/charmbracelet/lipgloss"

type Theme interface {
	DialogBox() lipgloss.Style
	Button() lipgloss.Style
	ActiveButton() lipgloss.Style
	Input() lipgloss.Style
	InputFocused() lipgloss.Style
	Checkbox() lipgloss.Style
	CheckboxFocused() lipgloss.Style

	Text() lipgloss.Style

	Subtle() lipgloss.AdaptiveColor
	Highlight() lipgloss.AdaptiveColor
	Special() lipgloss.AdaptiveColor
}

type DefaultTheme struct{}

func (t DefaultTheme) DialogBox() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 0).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)
}

func (t DefaultTheme) Button() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#888B7E")).
		Padding(0, 4).
		MarginTop(1)
}

func (t DefaultTheme) ActiveButton() lipgloss.Style {
	return t.Button().Copy().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#F25D94")).
		Underline(true)
}

func (t DefaultTheme) Input() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (t DefaultTheme) InputFocused() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
}

func (t DefaultTheme) Checkbox() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (t DefaultTheme) CheckboxFocused() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
}

func (t DefaultTheme) Text() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (t DefaultTheme) Subtle() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
}

func (t DefaultTheme) Highlight() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
}

func (t DefaultTheme) Special() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
}
