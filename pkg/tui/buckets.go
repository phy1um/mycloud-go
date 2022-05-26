package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jmoiron/sqlx"
)

type bucketView struct {
	tags   []string
	cursor int
	db     *sqlx.DB
}

func (b bucketView) Enter() {}
func (b bucketView) Exit()  {}

func (b bucketView) Update(msg tea.Msg, st *State) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			return bucketView{
				tags:   b.tags,
				cursor: intmin(b.cursor+1, len(b.tags)-1),
				db:     b.db,
			}, nil
		case "k", "up":
			return bucketView{
				tags:   b.tags,
				cursor: intmax(b.cursor-1, 0),
				db:     b.db,
			}, nil
		case "x", "enter":
			tag := b.tags[b.cursor]

			if tag == "all" {
				tag = ""
			}

			fv := fileView{
				tag:    tag,
				cursor: 0,
				db:     b.db,
			}
			st.PushView(&fv)
			return nil, nil
		}
	}

	return b, nil
}

func (b bucketView) View() string {
	parts := make([]string, len(b.tags)+2)
	parts = append(parts, "=== Select Tag to View ===")
	for i, tag := range b.tags {
		sel := " "
		if i == b.cursor {
			sel = "*"
		}

		parts = append(parts, fmt.Sprintf("[%s] Tag [%s]", sel, tag))
	}
	parts = append(parts, "\n -- Use J/K for down/up. Q to go back --\n")
	return strings.Join(parts, "\n")
}
