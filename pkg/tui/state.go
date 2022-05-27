package tui

import (
	"context"
	"strings"

	"github.com/phy1um/mycloud-go/config"
	"github.com/phy1um/mycloud-go/pkg/store"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type View interface {
	Enter(context.Context)
	Exit(context.Context)
	View() []string
	Update(context.Context, tea.Msg, *State) (View, tea.Cmd)
}

// State manages a stack of views and provides common utility like a back button
// When all views are popped the program ends
type State struct {
	viewStack []View
	store     store.Client
	cfg       *config.AppConfig
	ctx       context.Context
	done      bool
	drop      bool
}

func NewState(ctx context.Context, cfg *config.AppConfig, db *sqlx.DB) *State {
	log.Ctx(ctx).Info().Msg("initializing TUI state")
	return &State{
		ctx:   ctx,
		cfg:   cfg,
		store: store.NewClient(db),
	}
}

func (s *State) Init() tea.Cmd {
	err := s.store.Migrate()
	if err != nil {
		log.Ctx(s.ctx).Error().Stack().Err(err).Msg("failed to run migrate")
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
	newTop, cmd := top.Update(s.ctx, msg, s)
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
	v.Enter(s.ctx)
	s.drop = true
}

func (s *State) PopView() {
	s.Top().Exit(s.ctx)
	s.viewStack = s.viewStack[:len(s.viewStack)-1]
	if len(s.viewStack) == 0 {
		s.done = true
	}
	s.drop = true
}

func (s *State) Top() View {
	return s.viewStack[len(s.viewStack)-1]
}
