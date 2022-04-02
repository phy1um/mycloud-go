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
	"github.com/charmbracelet/wish/scp"
)

func main() {
	log.Println("start upload service")

	configPath := flag.String("c", "./config/local.yaml", "path to load config")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
		panic(err)
	}

	log.Printf("loaded config: %v", cfg)

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Upload.Port)
	fileHandler := scp.NewFileSystemHandler(cfg.App.FilePath)

	opts := internal.ServerOpts(
		cfg.App,
		wish.WithAddress(addr),
		wish.WithMiddleware(
			scp.Middleware(fileHandler, fileHandler),
		),
	)

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
