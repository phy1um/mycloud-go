package tui

import (
	"fmt"
	"log"
	"sshtest/internal/data"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jmoiron/sqlx"
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
	tag    string
	cursor int
	files  []data.File
	db     *sqlx.DB
	err    error
}

func (f *fileView) Enter() {
	var err error
	if f.tag != "" {
		err = f.db.Select(&f.files, "SELECT * FROM files WHERE tag=$1", f.tag)
	} else {
		err = f.db.Select(&f.files, "SELECT * FROM files", f.tag)
	}

	if err != nil {
		f.err = err
	}

	log.Printf("found %d files to display\n", len(f.files))
	for _, f := range f.files {
		log.Printf(" - %s\n", f.Path)
	}
}

func (f *fileView) Exit() {

}

func (f *fileView) Update(msg tea.Msg, st *State) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			return &fileView{
				files:  f.files,
				cursor: intmin(f.cursor+1, len(f.files)-1),
			}, nil
		case "k", "up":
			return &fileView{
				files:  f.files,
				cursor: intmax(f.cursor-1, 0),
			}, nil
		case "x", "enter":
			st.PushView(NewManageView(
				f.files[f.cursor].Id,
				f.files[f.cursor].Path,
			))
			return nil, nil
		}
	}

	return f, nil
}

func (f *fileView) View() string {
	if f.err != nil {
		return fmt.Sprintf(":: File View Error: %s", f.err.Error())
	}
	s := " :: File View :: \n"
	for i, file := range f.files {
		if i == f.cursor {
			s += fmt.Sprintf("[*] %s\n", file.Path)
		} else {
			s += fmt.Sprintf("[ ] %s\n", file.Path)
		}
	}
	s += "\n -- Use J/K for down/up. Q to quit --\n"
	return s
}
