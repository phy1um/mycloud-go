package internal

import (
	"context"
	"io/ioutil"

	"github.com/gliderlabs/ssh"
	"github.com/rs/zerolog/log"
)

type PublicKeyAuth struct {
	allowedKeys []ssh.PublicKey
	ctx         context.Context
}

func NewPublicKeyAuthFromFiles(ctx context.Context, paths []string) PublicKeyAuth {
	auth := PublicKeyAuth{
		ctx: ctx,
	}

	for _, file := range paths {
		key, err := loadPublicKey(file)
		if err != nil {
			log.Ctx(ctx).Error().Stack().Err(err).Msgf("failed to load public key: %s", key)
		} else {
			log.Ctx(ctx).Info().Msgf("loaded public key: %s", key.Type())
			auth.allowedKeys = append(auth.allowedKeys, key)
		}
	}

	return auth
}

func (p PublicKeyAuth) PublicKeyHandler(x ssh.Context, key ssh.PublicKey) bool {

	logg := log.Ctx(p.ctx).With().
		Str("incoming-key-type", key.Type()).
		Str("incoming-ssh-user", x.User()).
		Str("remote-address", x.RemoteAddr().String()).
		Logger()

	for _, allowedKey := range p.allowedKeys {
		if ssh.KeysEqual(key, allowedKey) {
			logg.Info().Msg("connection allowed")
			return true
		}
	}

	logg.Info().Msg("connection denied")
	return false
}

func loadPublicKey(path string) (ssh.PublicKey, error) {

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	key, _, _, _, err := ssh.ParseAuthorizedKey(b)
	return key, err
}
