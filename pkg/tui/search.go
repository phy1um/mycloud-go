package tui

import (
	"context"
	"fmt"
	"sshtest/pkg/store"
	"sshtest/pkg/tui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

type SearchKind string

const (
	All  SearchKind = "*"
	Name SearchKind = "name"
	Tag  SearchKind = "tag"
)

type fileSearchView struct {
	tags   []string
	cursor int
	store  store.Client
	box    styles.Box
	input  textinput.Model
	search SearchKind
}

func (b *fileSearchView) Enter(ctx context.Context) {
	b.input = textinput.New()
	b.input.Placeholder = "File Name"
	b.input.Focus()
	b.search = All
}

func (b *fileSearchView) Exit(ctx context.Context) {}

func (b *fileSearchView) Update(ctx context.Context, msg tea.Msg, st *State) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if b.search == Name {
				b.search = Tag
			} else if b.search == Tag {
				b.search = All
			} else {
				b.search = Name
			}
			return b, nil
		case "enter":
			var fn store.CursorFunc
			if b.search == All {
				fn = store.AllFiles()
			} else if b.search == Name {
				fn = store.FileNameSearch(b.input.Value())
			} else if b.search == Tag {
				fn = store.FileTagSearch((b.input.Value()))
			}

			fv := NewFileView(b.store, 10, fn)
			st.PushView(fv)
			return nil, nil
		default:
			log.Ctx(ctx).Info().Msg("updating textinput")
			n, cmd := b.input.Update(msg)
			b.input = n
			return b, cmd
		}
	}

	return b, nil
}

func (b fileSearchView) View() []string {
	title := styles.Title(lipgloss.NewStyle()).
		MarginLeft(3).
		Render("File Search")
	titleC := lipgloss.PlaceHorizontal(b.box.Width, lipgloss.Center, title)

	var main string
	if b.search != All {
		main = styles.Highlight(lipgloss.NewStyle()).
			MarginLeft(3).
			Render(b.input.View())
	}

	kind := lipgloss.NewStyle().
		Background(lipgloss.Color("#000044")).
		Foreground(lipgloss.Color("#FFFFFF")).
		MarginLeft(3).
		Render(fmt.Sprintf("<< %s >>", b.search))

	btm := b.box.Height - (lipgloss.Height(titleC) + lipgloss.Height(main))
	footer := lipgloss.PlaceVertical(btm, lipgloss.Bottom, "-- Use J/K for down/up. Q to go back --")
	return []string{lipgloss.JoinVertical(lipgloss.Top, titleC, main, kind, footer)}
}
