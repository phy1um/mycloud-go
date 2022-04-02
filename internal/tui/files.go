package tui

import (
	"fmt"
	"io/fs"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func intmax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func intmin(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

type fileView struct {
	files  []fs.FileInfo
	cursor int
}

func (f fileView) Enter() {
}

func (f fileView) Exit() {

}

func (f fileView) Update(msg tea.Msg, st *State) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			return fileView{
				files:  f.files,
				cursor: intmin(f.cursor+1, len(f.files)-1),
			}, nil
		case "k", "up":
			return fileView{
				files:  f.files,
				cursor: intmax(f.cursor-1, 0),
			}, nil
		case "x", "enter":
			st.PushView(createKey{
				path:     f.files[f.cursor].Name(),
				duration: 48 * time.Hour,
			})
			return nil, nil
		}
	}

	return f, nil
}

func (f fileView) View() string {
	s := " :: File View :: \n"
	for i, file := range f.files {
		if i == f.cursor {
			s += fmt.Sprintf("[*] %s\n", file.Name())
		} else {
			s += fmt.Sprintf("[ ] %s\n", file.Name())
		}
	}
	s += "\n -- Use J/K for down/up. Q to quit --\n"
	return s
}
