package tui

import (
	"fmt"
	"io/ioutil"
	"log"
	"sshtest/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jmoiron/sqlx"
)

type View interface {
	Enter()
	Exit()
	View() string
	Update(tea.Msg, *State) (View, tea.Cmd)
}

type State struct {
	viewStack []View
	db        *sqlx.DB
	cfg       *config.AppConfig
}

func NewState(cfg *config.AppConfig) *State {
	return &State{
		cfg: cfg,
	}
}

func (s *State) Init() tea.Cmd {
	db, err := sqlx.Open("sqlite3", s.cfg.DBFile)

	if err != nil {
		log.Println("failed to connect to database")
		return tea.Quit
	}

	s.db = db

	dir, err := ioutil.ReadDir(s.cfg.FilePath)
	if err != nil {
		fmt.Println("failed to read path")
	}

	baseView := fileView{
		files:  dir,
		cursor: 0,
	}

	s.PushView(baseView)

	return nil
}

func (s *State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if len(s.viewStack) == 0 {
				return s, tea.Quit
			}
			s.PopView()
		case "ctrl+c":
			return s, tea.Quit
		}
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
	return s.Top().View()
}

func (s *State) PushView(v View) {
	s.viewStack = append(s.viewStack, v)
	v.Enter()
}

func (s *State) PopView() {
	s.Top().Exit()
	s.viewStack = s.viewStack[:len(s.viewStack)-1]
}

func (s *State) Top() View {
	return s.viewStack[len(s.viewStack)-1]
}