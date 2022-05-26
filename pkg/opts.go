package internal

import (
	"sshtest/config"
	"time"

	"github.com/charmbracelet/wish"
	"github.com/gliderlabs/ssh"
)

func ServerOpts(cfg *config.AppConfig, extraOpts ...ssh.Option) []ssh.Option {
	auth := NewPublicKeyAuthFromFiles(cfg.AuthorizedKeyFiles)
	opts := []ssh.Option{
		wish.WithIdleTimeout(2 * time.Minute),
		wish.WithPublicKeyAuth(auth.PublicKeyHandler),
	}

	for _, key := range cfg.HostKeys {
		opts = append(opts, wish.WithHostKeyPath(key))
	}

	return append(opts, extraOpts...)
}
