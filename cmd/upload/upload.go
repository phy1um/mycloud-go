package main

import (
	"flag"
	"log"
	"os"

	"sshtest/cmd"
	"sshtest/config"
	"sshtest/internal"
)

func main() {
	log.Println("start upload service")

	configPath := flag.String("c", "./upload/local.yaml", "path to load config")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
		panic(err)
	}

	log.Printf("loaded config: %v", cfg)

	server, err := internal.MakeUploadServer(
		cfg.App.Host,
		cfg.App.Upload.Port,
		cfg.App.FilePath,
		cfg.App.Keys,
	)

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
