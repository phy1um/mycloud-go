package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sshtest/config"
	"sshtest/internal"
	"sshtest/internal/data"
)

func main() {
	log.Println("start serve service")

	configPath := flag.String("c", "./config/local.yaml", "path to load config")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
		panic(err)
	}

	health := internal.Health{
		Version: cfg.Meta.Version,
	}
	mux := http.NewServeMux()
	mux.Handle("/health", health)

	client, err := data.NewClient(cfg.App.DBFile, cfg.App.FilePath)
	if err != nil {
		panic(err)
	}

	server := data.NewServer(client)
	server.DefineServer(mux)

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Serve.Port)

	log.Printf("listening on %s\n", addr)

	err = http.ListenAndServe(addr, mux)
	if err != nil {
		panic(err)
	}
}
