package tui

import (
	"context"
	"fmt"
	"log"
	"sshtest/pkg/data"
	"sshtest/pkg/store"
	"sshtest/pkg/tui/styles"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type manageView struct {
	options menu
	file    *data.File
	tags    data.TagSet
	access  []*data.Access
	box     styles.Box
	store   store.Client
	err     error
}

func NewManageView(file *data.File, store store.Client) View {
	items := []item{
		fnItem{
			name: "Create Access Key",
			do: func(st *State) (View, tea.Cmd) {
				key, err := makeKeyInternal(file.Id, file.Path, time.Hour*48, store)
				st.PopView()
				st.PushView(createStatusView{
					key: key,
					err: err,
				})
				return nil, nil
			},
		},
		fnItem{
			name: "Add Tag",
			do: func(st *State) (View, tea.Cmd) {
				st.PushView(&setTagView{
					store: store,
					file:  file,
				})
				return nil, nil
			},
		},
		fnItem{
			name: "Remove Tags",
			do: func(st *State) (View, tea.Cmd) {
				return nil, nil
			},
		},
	}
	return &manageView{
		file:  file,
		store: store,
		options: menu{
			items:      items,
			renderBase: " - [%s] :: %s",
			sel:        "*",
			unsel:      " ",
		},
	}
}

func (m *manageView) Enter() {
	tags, err := m.store.GetTags(context.Background(), m.file)
	if err != nil {
		m.err = err
		return
	}
	m.tags = tags

	access, err := m.store.GetAccessKeys(context.Background(), m.file)
	if err != nil {
		m.err = err
		return
	}
	m.access = access
}

func (m *manageView) Exit() {}

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

func (m *manageView) View() []string {
	title := styles.Title(m.box.Style()).Render("Manage File")
	details := styles.Body(m.box.Style()).Border(lipgloss.NormalBorder()).Render(
		fmt.Sprintf("Name: %s\nPath: %s\nCreated: %s\n", m.file.Name, m.file.Path, m.file.Created.String()),
	)
	tags := styles.Body(m.box.Style()).Border(lipgloss.NormalBorder()).Render(
		"Tags\n\n" + strings.Join(m.tags, "\n"),
	)
	access := styles.Body(m.box.Style()).Border(lipgloss.NormalBorder()).Render(
		accessString(m.access),
	)
	meat := lipgloss.JoinHorizontal(lipgloss.Left, tags, access)
	body := styles.Body(m.box.Style()).Render(strings.Join(m.options.render(), "\n"))
	return []string{lipgloss.JoinVertical(lipgloss.Top, title, details, meat, body)}
}

type createStatusView struct {
	err error
	key string
}

func (c createStatusView) Enter() {}
func (c createStatusView) Exit()  {}
func (c createStatusView) View() []string {
	if c.err != nil {
		return []string{fmt.Sprintf("ERROR: %s", c.err.Error())}
	}
	return []string{fmt.Sprintf("Created access with key = %s", c.key)}
}

func (c createStatusView) Update(msg tea.Msg, st *State) (View, tea.Cmd) {
	return nil, nil
}

type setTagView struct {
	file  *data.File
	box   styles.Box
	input textinput.Model
	store store.Client
}

func (s *setTagView) Enter() {
	s.input = textinput.New()
	s.input.Placeholder = "New Tag Name"
	s.input.Focus()
}

func (s setTagView) Exit() {}

func (s setTagView) View() []string {
	title := styles.Title(s.box.Style()).Render(fmt.Sprintf("Tag File"))
	details := styles.Body(s.box.Style()).Border(lipgloss.NormalBorder()).Render(
		fmt.Sprintf("Name: %s\nPath: %s\nCreated: %s\n", s.file.Name, s.file.Path, s.file.Created.String()),
	)
	input := s.input.View()

	return []string{lipgloss.JoinVertical(lipgloss.Top, title, details, input)}
}

func (s *setTagView) Update(msg tea.Msg, st *State) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if len(s.input.Value()) > 0 {
				s.store.AddTag(context.Background(), s.file.Id, s.input.Value())
				st.PopView()
				st.PopView()
				return nil, nil
			}
		}
	}
	input, rs := s.input.Update(msg)
	s.input = input
	return s, rs
}

func accessString(keys []*data.Access) string {
	var sb strings.Builder
	sb.WriteString(" Access Keys\n\n")
	for _, k := range keys {
		b := fmt.Sprintf("[%s] (until %s)\n", k.Key, k.Until.String())
		sb.WriteString(b)
	}
	return sb.String()
}

func makeKeyInternal(id string, path string, duration time.Duration, store store.Client) (string, error) {
	key := data.RandomKey()
	access := data.Access{
		FileId:   id,
		Key:      key,
		UserCode: "",
		Created:  time.Now(),
		Until:    time.Now().Add(duration),
	}
	log.Printf("creating access key for %s: %s\n", path, key)
	err := store.CreateAccessKey(&access)

	if err != nil {
		return "", err
	}

	return key, nil
}
