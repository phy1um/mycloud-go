package styles

import "github.com/charmbracelet/lipgloss"

type Box struct {
	Width   int
	Height  int
	padding int
}

func (b Box) Style() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(b.Width).
		Height(b.Height).
		Padding(b.padding)
}

func Title(from lipgloss.Style) lipgloss.Style {
	return from.Copy().
		Border(lipgloss.NormalBorder()).
		MarginLeft(5).
		MarginRight(5).
		Padding(0, 1).
		Bold(true)
}

func Body(from lipgloss.Style) lipgloss.Style {
	return from.Copy()
}

func Highlight(from lipgloss.Style) lipgloss.Style {
	return from.Copy().
		BorderForeground(lipgloss.Color("228"))
}
