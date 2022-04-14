package tui

import (
	"fmt"
	"log"
	"sshtest/internal/data"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jmoiron/sqlx"
)

type manageView struct {
	options     menu
	sub         View
	path        string
	displayName string
}

func NewManageView(path string, displayName string) View {
	items := []item{
		fnItem{
			name: "Create Access Key",
			do: func(st *State) (View, tea.Cmd) {
				key, err := makeKeyInternal(path, time.Hour*48, st.db)
				st.PopView()
				st.PushView(createStatusView{
					key: key,
					err: err,
				})
				return nil, nil
			},
		},
		fnItem{
			name: "Set Tag",
			do: func(st *State) (View, tea.Cmd) {
				st.PushView(&setTagView{
					options:     menuFromStrings(st.cfg.Manage.Buckets, setTag(path)),
					path:        path,
					displayName: displayName,
				})
				return nil, nil
			},
		},
	}
	return &manageView{
		options: menu{
			items:      items,
			renderBase: " - [%s] :: %s",
			sel:        "*",
			unsel:      " ",
		},
	}
}

func (m *manageView) Enter() {}
func (m *manageView) Exit()  {}
func (m *manageView) Update(msg tea.Msg, st *State) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.options = m.options.next()
		case "k", "up":
			m.options = m.options.prev()
		case "x", "enter":
			return m.options.action(st)
		}
	}
	return m, nil
}

func (m *manageView) View() string {
	return fmt.Sprintf(`
::: Manage File [%s] :::

%s

::: Enter to confirm, Esc to go back, J/K to increase/decrease time :::`,
		m.displayName,
		m.options.render())
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
	log.Printf("creating access key for %s: %s\n", path, key)
	_, err := db.NamedExec(
		"INSERT INTO access_keys (path,key,user_code,display_name,created,until) VALUES (:path, :key, :user_code, :display_name, :created, :until)",
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

type setTagView struct {
	options     menu
	path        string
	displayName string
}

func (s setTagView) Enter() {}
func (s setTagView) Exit()  {}
func (s setTagView) View() string {
	return fmt.Sprintf(`
::: Set Tags for File [%s] :::

%s

::: --- :::`,
		s.displayName,
		s.options.render(),
	)
}

func (s *setTagView) Update(msg tea.Msg, st *State) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			s.options = s.options.next()
		case "k", "up":
			s.options = s.options.prev()
		case "x", "enter":
			return s.options.action(st)
		}
	}
	return s, nil
}

func setTag(id string) func(tag string, st *State) {
	return func(tag string, st *State) {
		file := data.File{
			Id:  id,
			Tag: tag,
		}
		log.Printf("trying to update tag for: %s to %s", id, tag)
		_, err := st.db.NamedExec("UPDATE files SET tag=:tag WHERE id = \""+id+"\"", &file)
		if err != nil {
			log.Printf("failed to update tag: %s (%s): %s", id, tag, err.Error())
		}
	}
}
