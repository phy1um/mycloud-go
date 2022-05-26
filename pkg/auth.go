package internal

import (
	"io/ioutil"
	"log"

	"github.com/gliderlabs/ssh"
)

type PublicKeyAuth struct {
	allowedKeys []ssh.PublicKey
}

func NewPublicKeyAuthFromFiles(paths []string) PublicKeyAuth {
	auth := PublicKeyAuth{}

	for _, file := range paths {
		key, err := loadPublicKey(file)
		if err != nil {
			log.Printf("failed to load public key %s: %s\n", file, err.Error())
		} else {
			log.Printf("loaded public key %s\n", file)
			auth.allowedKeys = append(auth.allowedKeys, key)
		}
	}

	return auth
}

func (p PublicKeyAuth) PublicKeyHandler(_ ssh.Context, key ssh.PublicKey) bool {
	log.Printf("testing key: %s\n", key.Type())

	for _, allowedKey := range p.allowedKeys {
		if ssh.KeysEqual(key, allowedKey) {
			return true
		}
	}
	log.Printf("denied access to %s\n", key.Type())
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
