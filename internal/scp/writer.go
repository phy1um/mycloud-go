package scp

import (
	"fmt"
	"io"
	"log"
	"os"
	"sshtest/internal/data"
	"time"

	"github.com/charmbracelet/wish/scp"
	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type serverWriter struct {
	root        string
	allowedTags []string
	db          *sqlx.DB
	failOnMkdir bool
}

func NewWriter(root string, db *sqlx.DB) *serverWriter {
	return &serverWriter{
		root:        root,
		allowedTags: []string{},
		db:          db,
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

	trackedFile := data.File{
		Id:      fileName,
		Path:    file.Filepath,
		Created: time.Now(),
		Tag:     "NONE",
	}

	// tell the database about this new file
	_, err = w.db.NamedExec(`INSERT INTO files (id,path,created,tag) VALUES(:id, :path, :created, :tag)`, &trackedFile)
	if err != nil {
		log.Printf("failed to create DB entry for file %s\n", fileName)
	}

	return copied, err
}

func (w serverWriter) prefixedFile(file string) string {
	return w.root + file
}
