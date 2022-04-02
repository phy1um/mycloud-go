package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sshtest/cmd"
	"sshtest/config"
	"sshtest/internal"
	"sshtest/internal/tui"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/gliderlabs/ssh"
)

func main() {
	log.Println("start upload service")

	rand.Seed(time.Now().UnixNano())

	configPath := flag.String("c", "./config/local.yaml", "path to load config")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %s\n", err.Error())
		panic(err)
	}

	log.Printf("loaded config: %v\n", cfg)

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Manage.Port)

	opts := internal.ServerOpts(
		cfg.App,
		wish.WithAddress(addr),
		wish.WithMiddleware(
			bm.Middleware(MakeFolderHandler(cfg.App)),
		),
	)

	log.Printf("creating server @ %s\n", addr)
	server, err := wish.NewServer(opts...)

	if err != nil {
		log.Fatalf("failed to create server: %s", err.Error())
		panic(err)
	}

	err = cmd.RunSSHServer(server)
	if err != nil {
		log.Fatalf("scp server error: %s", err.Error())
	}

	os.Exit(0)
}

func MakeFolderHandler(cfg *config.AppConfig) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		_, _, active := s.Pty()
		if !active {
			fmt.Println("no active terminal, skipping")
			return nil, nil
		}

		m := tui.NewState(cfg)
		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}
}
