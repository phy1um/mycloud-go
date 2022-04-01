package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gliderlabs/ssh"
)

func RunSSHServer(server *ssh.Server) error {
	interrupted := make(chan os.Signal, 1)
	done := make(chan error, 1)

	signal.Notify(interrupted, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		done <- err
	}()

	go func() {
		<-interrupted
		done <- nil
	}()

	err := <-done
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer func() { cancel() }()

	err = server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
