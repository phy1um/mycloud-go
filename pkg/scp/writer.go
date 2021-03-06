package scp

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/phy1um/mycloud-go/pkg/data"
	"github.com/phy1um/mycloud-go/pkg/store"

	"github.com/charmbracelet/wish/scp"
	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// serverWriter is a front-door for files uploaded with SCP
type serverWriter struct {
	root        string
	allowedTags []string
	store       store.Client
	failOnMkdir bool
}

// NewWriter creates a writer that registers files in a database and moves everything to one directory
func NewWriter(root string, db *sqlx.DB) *serverWriter {
	return &serverWriter{
		root:        root,
		allowedTags: []string{},
		store:       store.NewClient(db),
		failOnMkdir: true,
	}
}

func (w serverWriter) Mkdir(_ ssh.Session, dir *scp.DirEntry) error {
	log.Printf("cannot mkdir: %s\n", dir.Name)

	if w.failOnMkdir {
		return fmt.Errorf("directory creation is disabled")
	}

	return nil
}

func (w serverWriter) Write(s ssh.Session, file *scp.FileEntry) (int64, error) {
	// All incoming files get a random UUID path
	fileName := uuid.New().String()
	f, err := os.OpenFile(w.prefixedFile(fileName), os.O_TRUNC|os.O_RDWR|os.O_CREATE, file.Mode)
	defer f.Close()

	if err != nil {
		return 0, fmt.Errorf("failed to open file: %q: %w", file.Filepath, err)
	}

	copied, err := io.Copy(f, file.Reader)
	if err != nil {
		return 0, fmt.Errorf("failed to write file: %q: %w", file.Filepath, err)
	}

	// Tag all incoming files as "new"
	tags := data.TagSet([]string{"new"})

	// Create the file and tag entries in the database
	w.store.CreateFile(s.Context(), fileName, file.Filepath, tags)

	if err != nil {
		return 0, err
	}

	return copied, nil
}

func (w serverWriter) prefixedFile(file string) string {
	return w.root + file
}
