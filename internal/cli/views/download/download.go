package download

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/cli/theme"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/librespot/core"
	"github.com/d3xter-dev/not-a-spotify-downloader/internal/services/download"
	Spotify "github.com/d3xter-dev/not-a-spotify-downloader/internal/spotify"
	"sort"
	"time"
)

type View struct {
	currentlyDownloading map[string]download.Item
	itemsDownloaded      int
	downloadChannel      chan download.Item

	session  *core.Session
	playlist *Spotify.SelectedListContent

	width    int
	height   int
	spinner  spinner.Model
	progress progress.Model
	theme    theme.Theme
}

func NewView(theme theme.Theme) *View {
	return &View{
		theme:                theme,
		currentlyDownloading: map[string]download.Item{},
		spinner:              spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(lipgloss.NewStyle().Foreground(theme.Special()))),
		progress:             progress.New(progress.WithDefaultScaledGradient(), progress.WithoutPercentage()),
	}
}

func (v *View) GetName() string {
	return "Downloading..."
}

func (v *View) Init() tea.Cmd {
	return tea.Batch(v.spinner.Tick)
}

func (v *View) SetSize(width, height int) {
	v.width, v.height = width, height
}

func (v *View) Update(message tea.Msg) tea.Cmd {
	switch msg := message.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		v.spinner, cmd = v.spinner.Update(msg)
		return cmd

	case download.Item:
		if msg.Status == download.StatusComplete {
			v.itemsDownloaded++
			delete(v.currentlyDownloading, msg.Path)

			if v.itemsDownloaded == len(v.playlist.GetContents().GetItems()) {
				return tea.Quit
			}

			return tea.Batch(WatchDownload(v.downloadChannel), tea.Printf("%s %s", lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓"), msg.Name))
		}

		if msg.Status == download.StatusFailed {
			v.itemsDownloaded++
			delete(v.currentlyDownloading, msg.Path)
			return tea.Batch(WatchDownload(v.downloadChannel), tea.Printf("%s %s", lipgloss.NewStyle().Foreground(lipgloss.Color("#F00")).SetString("✗"), msg.Name))
		}

		v.currentlyDownloading[msg.Path] = msg
		return WatchDownload(v.downloadChannel)

	case setPlaylistMsg:
		v.playlist, v.session = msg.playlist, msg.session
		ch, err := download.StartProcessing(v.session, v.playlist)
		if err != nil {
			panic("Failed to start processing playlist: " + err.Error())
		}
		v.downloadChannel = ch
		return WatchDownload(ch)
	}

	return nil
}

func (v *View) Render() string {
	title := v.theme.Text().MarginTop(1).Foreground(v.theme.Special()).Bold(true).Render("Processing playlist: " + v.playlist.GetAttributes().GetName())

	var list []string
	//sort map by start time
	var currentlyDownloading = make([]download.Item, 0, len(v.currentlyDownloading))
	for _, item := range v.currentlyDownloading {
		currentlyDownloading = append(currentlyDownloading, item)
	}

	sort.Slice(currentlyDownloading, func(i, j int) bool {
		return currentlyDownloading[i].StartTime.Before(currentlyDownloading[j].StartTime)
	})

	for _, item := range currentlyDownloading {
		prefix := v.spinner.View()
		suffix := fmt.Sprintf("(%fs)", time.Since(item.StartTime).Seconds())

		maxWidth := max(0, v.width-lipgloss.Width(prefix+suffix)-2)

		itemName := v.theme.Text().Foreground(v.theme.Highlight()).Render(item.Name)
		info := v.theme.Text().Width(maxWidth).Render("Downloading " + itemName)
		row := lipgloss.JoinHorizontal(lipgloss.Center, prefix, info, suffix)

		list = append(list, row)
	}
	content := lipgloss.JoinVertical(lipgloss.Center, list...)

	maxItems := len(v.playlist.GetContents().GetItems())
	downloadProgress := fmt.Sprintf("(%d/%d)", v.itemsDownloaded, maxItems)
	downloadPercentage := float64(v.itemsDownloaded) / float64(maxItems)
	v.progress.Width = v.width - lipgloss.Width(downloadProgress)

	statusbar := lipgloss.JoinHorizontal(lipgloss.Left, v.progress.ViewAs(downloadPercentage), downloadProgress)

	return lipgloss.JoinVertical(lipgloss.Left, title, v.theme.DialogBox().Width(v.width-2).Render(content), statusbar)
}
