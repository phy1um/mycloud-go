package scp

import (
	"fmt"
	"io/fs"
	"log"

	"github.com/charmbracelet/wish/scp"
	"github.com/gliderlabs/ssh"
)

type serverReader struct {
	root        string
	allowedTags []string
}

func NewReader(at string) *serverReader {
	return &serverReader{
		root:        at,
		allowedTags: []string{},
	}
}

func (c serverReader) Glob(_ ssh.Session, s string) ([]string, error) {
	return []string{s}, nil
}

func (c serverReader) WalkDir(session ssh.Session, dir string, fn fs.WalkDirFunc) error {
	log.Printf("trying to walk %s\n", dir)
	return fmt.Errorf("failed to walk %s", dir)
}

func (c serverReader) NewDirEntry(session ssh.Session, dir string) (*scp.DirEntry, error) {
	return nil, fmt.Errorf("failed to create dir %s", dir)
}

func (c serverReader) NewFileEntry(session ssh.Session, file string) (*scp.FileEntry, func() error, error) {
	log.Printf("creating file %s\n", file)
	return nil, nil, fmt.Errorf("failed to create file %s", file)
}
