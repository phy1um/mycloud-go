package internal

import (
	"context"
	"time"

	"github.com/phy1um/mycloud-go/config"

	"github.com/charmbracelet/wish"
	"github.com/gliderlabs/ssh"
)

// ServerOpts makes sensible default ssh Options based on our config
func ServerOpts(ctx context.Context, cfg *config.AppConfig, extraOpts ...ssh.Option) []ssh.Option {
	auth := NewPublicKeyAuthFromFiles(ctx, cfg.AuthorizedKeyFiles)
	opts := []ssh.Option{
		wish.WithIdleTimeout(2 * time.Minute),
		wish.WithPublicKeyAuth(auth.PublicKeyHandler),
	}

	for _, key := range cfg.HostKeys {
		opts = append(opts, wish.WithHostKeyPath(key))
	}

	return append(opts, extraOpts...)
}
