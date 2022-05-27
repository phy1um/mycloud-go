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
}

func NewManageView(file *data.File, store store.Client) View {
	items := []item{
		fnItem{
			name: "Create Access Key",
			do: func(st *State) (View, tea.Cmd) {
				key, err := makeKeyInternal(file.Path, time.Hour*48, store)
				st.PopView()
				st.PushView(createStatusView{
					key: key,
					err: err,
				})
				return nil, nil
			},
		},
		fnItem{
			name: "Manage Tags",
			do: func(st *State) (View, tea.Cmd) {
				st.PushView(&setTagView{
					file:  file,
					store: store,
				})
				return nil, nil
			},
		},
	}
	return &manageView{
		file: file,
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

func (m *manageView) View() []string {
	res := []string{
		fmt.Sprintf("::: Manage File [%s] :::", m.file.Name),
		"",
	}
	res = append(res, m.options.render()...)
	res = append(res, "", "::: Enter to confirm, Esc to go back, J/K to increase/decrease time :::")
	return res
}

func makeKeyInternal(path string, duration time.Duration, store store.Client) (string, error) {
	key := data.RandomKey()
	access := data.Access{
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
	file        *data.File
	box         styles.Box
	input       textinput.Model
	store       store.Client
	tags        []string
	err         error
	focusAddTag bool
	cursor      int
}

func (s *setTagView) Enter() {
	tags, err := s.store.GetTags(context.Background(), s.file)
	if err != nil {
		s.err = err
		return
	}
	s.tags = tags
	s.input = textinput.New()
	s.input.Placeholder = "New Tag Name"
	s.input.Focus()
}

func (s setTagView) Exit() {}

func (s setTagView) View() []string {
	if s.err != nil {
		return []string{s.err.Error()}
	}

	title := styles.Title(s.box.Style()).Render(fmt.Sprintf("Manage File [%s]", s.file.Name))
	details := styles.Body(s.box.Style()).Border(lipgloss.NormalBorder()).Render(
		fmt.Sprintf("Path: %s\nCreated: %s\n", s.file.Path, s.file.Created.String()),
	)
	tags := styles.Body(s.box.Style()).Border(lipgloss.NormalBorder()).Render(
		"Tags\n\n" + strings.Join(s.tags, "\n"),
	)
	input := s.input.View()

	return []string{lipgloss.JoinVertical(lipgloss.Top, title, details, tags, input)}
}

func (s *setTagView) Update(msg tea.Msg, st *State) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			newTag := s.input.Value()
			if len(newTag) == 0 {
				break
			}
			err := s.store.AddTag(context.Background(), s.file.Id, newTag)
			if err == nil {
				s.tags = append(s.tags, newTag)
				s.input.SetValue("")
			} else {
				s.err = err
			}
		default:
			input, rs := s.input.Update(msg)
			s.input = input
			return s, rs
		}
	}
	return s, nil
}

/*
func setTag(id string) func(tag string, st *State) {
	return func(tag string, st *State) {
		file := data.File{
			Id: id,
		}
		log.Printf("trying to update tag for: %s to %s", id, tag)
		_, err := st.db.NamedExec("UPDATE files SET tag=:tag WHERE id = \""+id+"\"", &file)
		if err != nil {
			log.Printf("failed to update tag: %s (%s): %s", id, tag, err.Error())
		}
	}
}
*/
