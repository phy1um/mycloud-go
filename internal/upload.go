package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/scp"
	"github.com/gliderlabs/ssh"
)

func MakeUploadServer(host string, port int, dir string, keys []string) (*ssh.Server, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	fileHandler := scp.NewFileSystemHandler(dir)

	opts := []ssh.Option{
		wish.WithAddress(addr),
		wish.WithIdleTimeout(2 * time.Minute),
		wish.WithMiddleware(
			scp.Middleware(fileHandler, fileHandler),
		),
	}

	for _, key := range keys {
		opts = append(opts, wish.WithHostKeyPath(key))
	}

	server, err := wish.NewServer(opts...)

	if err != nil {
		return nil, err
	}

	log.Printf("scp upload server created @ %s\n", addr)

	return server, nil
}
