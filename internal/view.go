package internal

import (
	"fmt"
	"io/fs"

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
	prev   string
}

func (f fileView) Init() tea.Cmd {
	return nil
}

func (f fileView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			return fileView{
				files:  f.files,
				cursor: intmin(f.cursor+1, len(f.files)-1),
				prev:   f.prev,
			}, nil
		case "k":
			return fileView{
				files:  f.files,
				cursor: intmax(f.cursor-1, 0),
				prev:   f.prev,
			}, nil
		case "q":
			return f, tea.Quit
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
