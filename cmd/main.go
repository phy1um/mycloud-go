package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sshtest/config"
	internal "sshtest/pkg"
	"sshtest/pkg/tui"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/scp"
	"github.com/gliderlabs/ssh"
	"github.com/jmoiron/sqlx"

	scp2 "sshtest/pkg/scp"

	"sshtest/pkg/fetch"

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
		log.Ctx(ctx).Fatal().Err(err).Msgf("failed to load config: %s", configPath)
	}

	log.Ctx(ctx).Info().Msgf("loaded config: %+v", cfg)

	annotatedLogger := log.Ctx(ctx).With().
		Str("app-version", cfg.Meta.Version).
		Str("database", cfg.App.DBFile).
		Stack().
		Logger()
	ctx = annotatedLogger.WithContext(ctx)

	db, err := sqlx.Open("sqlite3", cfg.App.DBFile)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("failed to connect to database")
	}

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Manage.Port)

	// Use default SCP outgoing file handler, but we customize the incoming handler
	fileHandler := scp.NewFileSystemHandler(cfg.App.FilePath)

	// Create server options based on sensible defaults
	opts := internal.ServerOpts(
		ctx,
		&cfg.App,
		wish.WithAddress(addr),
		wish.WithMiddleware(
			bm.Middleware(newTUIForConfig(ctx, &cfg.App, db)),
			scp.Middleware(fileHandler, scp2.NewWriter(cfg.App.FilePath, db)),
		),
	)

	log.Ctx(ctx).Info().Msgf("creating server @ %s", addr)

	server, err := wish.NewServer(opts...)

	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("failed to create server")
	}

	// Setup a health/status endpoint
	health := internal.Health{
		Version: cfg.Meta.Version,
	}

	mux := http.NewServeMux()
	mux.Handle("/health", health)

	// Initialize our file retrieval server
	client, err := fetch.NewClient(db, cfg.App.FilePath)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("failed to create data client")
	}

	ds := fetch.NewServer(ctx, client)
	ds.DefineServer(mux)

	httpAddr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Serve.Port)

	// Run everything - this function blocks until we are done
	err = internal.RunServers(httpAddr, server, mux)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("error running servers")
	}

	os.Exit(0)
}

// Create a TUI for managing files based on our app config
func newTUIForConfig(ctx context.Context, cfg *config.AppConfig, db *sqlx.DB) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		pty, _, active := s.Pty()

		logg := log.Ctx(ctx).With().
			Str("ssh-user", s.User()).
			Str("ssh-key-type", s.PublicKey().Type()).
			Str("ssh-pty", fmt.Sprintf("PTY(%t): %d x %d", active, pty.Window.Width, pty.Window.Height)).
			Str("remote-address", s.RemoteAddr().String()).
			Logger()

		if !active {
			logg.Error().Stack().Err(errors.New("no PTY active for connected client"))
			return nil, nil
		}

		wctx := logg.WithContext(s.Context())

		m := tui.NewState(wctx, cfg, db)

		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}
}
