package tui

import (
	"log"
	"sshtest/config"
	"sshtest/pkg/store"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jmoiron/sqlx"
)

type View interface {
	Enter()
	Exit()
	View() []string
	Update(tea.Msg, *State) (View, tea.Cmd)
}

type State struct {
	viewStack []View
	store     store.Client
	cfg       *config.AppConfig
	done      bool
	drop      bool
}

func NewState(cfg *config.AppConfig, db *sqlx.DB) *State {
	log.Printf("initializing state with DB=%v\n", db)
	return &State{
		cfg:   cfg,
		store: store.NewClient(db),
	}
}

func (s *State) Init() tea.Cmd {
	log.Printf("Has some store.. %+v\n", s.store)
	err := s.store.Migrate()
	if err != nil {
		log.Printf("failed to run migrate: %s", err.Error())
	}

	tagSet := []string{"all"}
	tagSet = append(tagSet, s.cfg.Manage.Buckets...)

	baseView := fileSearchView{
		store: s.store,
		tags:  tagSet,
	}
	s.PushView(&baseView)

	return nil
}

func (s *State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if s.done {
		return s, tea.Quit
	}

	if s.drop {
		s.drop = false
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if len(s.viewStack) == 0 {
				return s, tea.Quit
			}
			s.PopView()
		case "ctrl+c":
			return s, tea.Quit
		}
	}

	if len(s.viewStack) == 0 {
		return s, nil
	}

	top := s.Top()
	newTop, cmd := top.Update(msg, s)
	if newTop != nil {
		s.viewStack[len(s.viewStack)-1] = newTop
	}

	return s, cmd
}

func (s *State) View() string {
	if len(s.viewStack) == 0 {
		return " ::: LOADING ::: "
	}
	s.drop = false
	lines := s.Top().View()
	return strings.Join(lines, "\n")
}

func (s *State) PushView(v View) {
	s.viewStack = append(s.viewStack, v)
	v.Enter()
	s.drop = true
}

func (s *State) PopView() {
	s.Top().Exit()
	s.viewStack = s.viewStack[:len(s.viewStack)-1]
	if len(s.viewStack) == 0 {
		s.done = true
	}
	s.drop = true
}

func (s *State) Top() View {
	return s.viewStack[len(s.viewStack)-1]
}
