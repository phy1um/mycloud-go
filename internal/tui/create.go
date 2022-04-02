package tui

import (
	"fmt"
	"sshtest/internal/data"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jmoiron/sqlx"
)

type createKey struct {
	path     string
	duration time.Duration
	back     tea.Model
}

func (c createKey) View() string {
	return fmt.Sprintf(`

 ::: Creating Access Key for File :::

  TARGET FILE: %s
  DURATION: %s

  ::: Enter to confirm, Esc to go back, J/K to increase/decrease time :::
	`, c.path, c.duration.String())
}

func (c createKey) Update(msg tea.Msg, st *State) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "x", "enter":
			key, err := makeKeyInternal(c.path, c.duration, st.db)
			st.PopView()
			st.PushView(createStatusView{
				key: key,
				err: err,
			})
		}
	}
	return nil, nil
}

func (c createKey) Enter() {
}

func (c createKey) Exit() {

}

func makeKeyInternal(path string, duration time.Duration, db *sqlx.DB) (string, error) {
	key := data.RandomKey()
	access := data.Access{
		Path:     path,
		Key:      key,
		UserCode: "",
		Created:  time.Now(),
		Until:    time.Now().Add(duration),
	}
	_, err := db.NamedExec(
		"INSERT INTO access_keys (path,key,user_code,created,until) VALUES (:path, :key, :user_code, :created, :until)",
		&access,
	)

	if err != nil {
		return "", err
	}

	return key, nil
}

type createStatusView struct {
	err error
	key string
}

func (c createStatusView) Enter() {}
func (c createStatusView) Exit()  {}
func (c createStatusView) View() string {
	if c.err != nil {
		return fmt.Sprintf("ERROR: %s", c.err.Error())
	}
	return fmt.Sprintf("Created access with key = %s", c.key)
}

func (c createStatusView) Update(msg tea.Msg, st *State) (View, tea.Cmd) {
	return nil, nil
}
