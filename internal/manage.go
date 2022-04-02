package internal

import (
	"fmt"
	"io/ioutil"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gliderlabs/ssh"
)

func MakeFolderHandler(path string) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		_, _, active := s.Pty()
		if !active {
			fmt.Println("no active terminal, skipping")
			return nil, nil
		}

		dir, err := ioutil.ReadDir(path)
		if err != nil {
			fmt.Println("failed to read path")
			return nil, nil
		}

		m := fileView{
			files:  dir,
			prev:   ".",
			cursor: 0,
		}
		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}
}
