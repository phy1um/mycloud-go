package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sshtest/config"
	internal "sshtest/pkg"
	"sshtest/pkg/data"
	"sshtest/pkg/tui"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/scp"
	"github.com/gliderlabs/ssh"
	"github.com/jmoiron/sqlx"

	scp2 "sshtest/pkg/scp"
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

	log.Printf("loaded config: %+v\n", cfg)

	db, err := sqlx.Open("sqlite3", cfg.App.DBFile)
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Manage.Port)

	fileHandler := scp.NewFileSystemHandler(cfg.App.FilePath)

	opts := internal.ServerOpts(
		&cfg.App,
		wish.WithAddress(addr),
		wish.WithMiddleware(
			bm.Middleware(newTUIForFolder(&cfg.App, db)),
			scp.Middleware(fileHandler, scp2.NewWriter(cfg.App.FilePath, db)),
		),
	)

	log.Printf("creating server @ %s\n", addr)
	server, err := wish.NewServer(opts...)

	if err != nil {
		log.Fatalf("failed to create server: %s", err.Error())
		panic(err)
	}

	health := internal.Health{
		Version: cfg.Meta.Version,
	}

	mux := http.NewServeMux()
	mux.Handle("/health", health)

	client, err := data.NewClient(db, cfg.App.FilePath)
	if err != nil {
		panic(err)
	}

	ds := data.NewServer(client)
	ds.DefineServer(mux)

	httpAddr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Serve.Port)

	err = internal.RunServers(httpAddr, server, mux)
	if err != nil {
		log.Fatalf("scp server error: %s", err.Error())
	}

	os.Exit(0)
}

func newTUIForFolder(cfg *config.AppConfig, db *sqlx.DB) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		_, _, active := s.Pty()
		if !active {
			fmt.Println("no active terminal, skipping")
			return nil, nil
		}

		m := tui.NewState(cfg, db)
		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}
}
