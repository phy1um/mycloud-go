package tui

import (
	"context"
	"fmt"
	"sshtest/pkg/data"
	"sshtest/pkg/store"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
)

func intmax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func intmin(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// fileView displays a list of files
// TODO: pagination
type fileView struct {
	cursor   int
	files    []*data.File
	store    store.Client
	dbCursor store.CursorKey
	err      error
}

func NewFileView(ss store.Client, pageSize int, fn store.CursorFunc) *fileView {
	return &fileView{
		cursor:   0,
		store:    ss,
		dbCursor: ss.NewCursor(pageSize, "name", store.Descend, fn),
	}
}

// Enter is called when this view is first shown
func (f *fileView) Enter(ctx context.Context) {
	var err error

	// Fetch the files on startup
	files, err := f.store.GetFiles(ctx, f.dbCursor)
	if err != nil {
		log.Ctx(ctx).Error().Stack().Err(err).Msg("failed to get files")
		f.err = err
		return
	}

	f.files = files

	log.Ctx(ctx).Info().Msgf("found %d files to display\n", len(f.files))
}

// Exit is called when this view is no longer needed
func (f *fileView) Exit(_ context.Context) {
	f.store.DestroyCursor(f.dbCursor)
}

// Update is called for each bubbletea event
// It allows for selecting a file to manage
func (f *fileView) Update(ctx context.Context, msg tea.Msg, st *State) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			f.cursor = intmin(f.cursor+1, len(f.files)-1)
		case "k", "up":
			f.cursor = intmax(f.cursor-1, 0)
		case "x", "enter":
			if f.files == nil {
				return f, nil
			}
			log.Ctx(ctx).Info().Msg("selecting file to view")
			// Show the manage view for the selected file
			st.PushView(NewManageView(
				f.files[f.cursor],
				f.store,
			))
			return nil, nil
		}
	}

	return f, nil
}

// View renders the list of files
func (f *fileView) View() []string {
	if f.err != nil {
		return []string{fmt.Sprintf(":: File View Error: %s", f.err.Error())}
	}
	s := []string{" :: File View ::"}
	for i, file := range f.files {
		if i == f.cursor {
			s = append(s, fmt.Sprintf("[*] %s (%s)", file.Name, file.Path))
		} else {
			s = append(s, fmt.Sprintf("[ ] %s (%s)", file.Name, file.Path))
		}
	}
	s = append(s, "\n -- Use J/K for down/up. Q to quit --")
	return s
}
