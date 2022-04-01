package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"sshtest/config"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/gliderlabs/ssh"
)

func MakeManageServer(cfg *config.AppConfig) (*ssh.Server, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Manage.Port)

	opts := []ssh.Option{
		wish.WithAddress(addr),
		wish.WithIdleTimeout(2 * time.Minute),
		wish.WithMiddleware(
			bm.Middleware(makeFolderHandler(cfg.FilePath)),
		),
	}

	for _, key := range cfg.Keys {
		opts = append(opts, wish.WithHostKeyPath(key))
	}

	server, err := wish.NewServer(opts...)

	if err != nil {
		return nil, err
	}

	log.Printf("manage upload server created @ %s\n", addr)

	return server, nil
}

func makeFolderHandler(path string) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
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
