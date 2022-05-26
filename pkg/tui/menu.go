package tui

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type item interface {
	Name() string
	Action(st *State) (View, tea.Cmd)
}

type fnItem struct {
	name string
	do   func(st *State) (View, tea.Cmd)
}

func (f fnItem) Name() string {
	return f.name
}

func (f fnItem) Action(st *State) (View, tea.Cmd) {
	return f.do(st)
}

func menuFromStrings(names []string, fn func(s string, st *State)) menu {
	items := make([]item, 0, len(names))
	for _, name := range names {
		nn := name
		items = append(items, fnItem{
			name: nn,
			do: func(st *State) (View, tea.Cmd) {
				fn(nn, st)
				return nil, nil
			},
		})
	}
	return menu{
		items:      items,
		renderBase: " [%s] %s",
		sel:        "*",
		unsel:      " ",
	}
}

type menu struct {
	items      []item
	cursor     int
	renderBase string
	sel        string
	unsel      string
}

func (m menu) next() menu {
	log.Printf("menu: next()\n")
	return menu{
		items:      m.items,
		cursor:     intmin(m.cursor+1, len(m.items)-1),
		renderBase: m.renderBase,
		sel:        m.sel,
		unsel:      m.unsel,
	}
}

func (m menu) prev() menu {
	log.Printf("menu: prev()\n")
	return menu{
		items:      m.items,
		cursor:     intmax(m.cursor-1, 0),
		renderBase: m.renderBase,
		sel:        m.sel,
		unsel:      m.unsel,
	}
}

func (m menu) action(st *State) (View, tea.Cmd) {
	item := m.items[m.cursor]
	return item.Action(st)
}

func (m menu) render() string {
	log.Printf("rendering menu, base=\"%s\"\n", m.renderBase)
	li := make([]string, len(m.items))
	for i, item := range m.items {
		sel := m.sel
		if i != m.cursor {
			sel = m.unsel
		}
		if item == nil {
			log.Printf("nil item in menu @ %d\n", i)
			continue
		}
		li = append(li, fmt.Sprintf(m.renderBase, sel, item.Name()))
	}
	return strings.Join(li, "\n")
}
