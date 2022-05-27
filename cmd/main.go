package main

import (
	"context"
	"flag"
	"fmt"
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

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := log.Logger
	ctx := log.Logger.WithContext(context.Background())

	logger.Info().Msg("server startup")

	rand.Seed(time.Now().UnixNano())

	configPath := flag.String("c", "./config/local.yaml", "path to load config")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("failed to load config: %s", configPath)
		panic(err)
	}

	log.Ctx(ctx).Info().Msgf("loaded config: %+v", cfg)

	db, err := sqlx.Open("sqlite3", cfg.App.DBFile)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to connect to database")
		panic(err)
	}

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Manage.Port)

	fileHandler := scp.NewFileSystemHandler(cfg.App.FilePath)

	opts := internal.ServerOpts(
		&cfg.App,
		wish.WithAddress(addr),
		wish.WithMiddleware(
			bm.Middleware(newTUIForFolder(ctx, &cfg.App, db)),
			scp.Middleware(fileHandler, scp2.NewWriter(cfg.App.FilePath, db)),
		),
	)

	log.Ctx(ctx).Info().Msgf("creating server @ %s", addr)

	server, err := wish.NewServer(opts...)

	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to create server")
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

	ds := data.NewServer(ctx, client)
	ds.DefineServer(mux)

	httpAddr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Serve.Port)

	err = internal.RunServers(httpAddr, server, mux)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("error running servers")
	}

	os.Exit(0)
}

func newTUIForFolder(ctx context.Context, cfg *config.AppConfig, db *sqlx.DB) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		pty, _, active := s.Pty()
		if !active {
			fmt.Println("no active terminal, skipping")
			return nil, nil
		}

		logg := log.Ctx(ctx).With().
			Str("ssh-user", s.User()).
			Str("ssh-key-type", s.PublicKey().Type()).
			Str("ssh-pty", fmt.Sprintf("PTY(%t): %d x %d", active, pty.Window.Width, pty.Window.Height)).
			Logger()

		wctx := logg.WithContext(s.Context())

		m := tui.NewState(wctx, cfg, db)

		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}
}
