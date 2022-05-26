package file

import (
	"sshtest/internal/data"
	"sshtest/internal/tui"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type fileManageView struct {
	title        lipgloss.Style
	body         lipgloss.Style
	inputs       []textinput.Model
	fieldNames   []string
	currentTag   string
	possibleTags []string
	file         data.File
	cursor       int
	savePrompt   bool
	dirty        bool
}

func NewFileManageView()

func (f fileManageView) Enter() {}

func (f fileManageView) Exit() {}

func (f fileManageView) Update(msg tea.Msg, st *tui.State) (tui.View, tea.Cmd) {
	return f, nil
}

func (f fileManageView) View() string {
	sb := strings.Builder{}
	body := strings.Builder{}
	sb.WriteString(f.title.Render("Manage File"))

	body.WriteString("id = ")
	body.WriteString(f.file.Id)
	body.WriteByte('\n')

}
