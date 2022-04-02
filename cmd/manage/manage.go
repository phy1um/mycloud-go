package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sshtest/cmd"
	"sshtest/config"
	"sshtest/internal"

	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
)

func main() {
	log.Println("start upload service")

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
			bm.Middleware(internal.MakeFolderHandler(cfg.App.FilePath)),
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
