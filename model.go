/*
Package bublog to show your logs in TUI and duplicate them to other places.
*/
package bublog

import (
	"errors"
	"io"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	StickToBottom key.Binding
	ScrollDown    key.Binding
	ScrollUp      key.Binding
}

type Model struct {
	AdditionalWriters []io.Writer
	Viewer            *TextViewer
	KeyMap            *KeyMap
}

var _ tea.Model = &Model{}

func NewModel(initialText string, additionalWriters ...io.Writer) *Model {
	v := NewTextViewer([]rune(initialText))
	v.SwitchStickToBottom()
	return &Model{additionalWriters, v, DefaultKeyMap()}
}

func DefaultKeyMap() *KeyMap {
	return &KeyMap{
		StickToBottom: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "Switch stick to bottom"),
		),
		ScrollDown: key.NewBinding(
			key.WithKeys("j", tea.KeyDown.String()),
			key.WithHelp("j/↓", "Scroll down"),
		),
		ScrollUp: key.NewBinding(
			key.WithKeys("k", tea.KeyUp.String()),
			key.WithHelp("k/↑", "Scroll up"),
		),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.ScrollDown):
			m.Viewer.ScrollDown()
			return m, nil
		case key.Matches(msg, m.KeyMap.ScrollUp):
			m.Viewer.ScrollUp()
			return m, nil
		case key.Matches(msg, m.KeyMap.StickToBottom):
			m.Viewer.SwitchStickToBottom()
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	toShow := make([]rune, 0, m.Viewer.Height*m.Viewer.Width)
	for _, line := range m.Viewer.View() {
		toShow = append(toShow, line...)
	}
	return string(toShow)
}

// SetSize accepts up to two arguments. First fill be Width, second - Height.
func (m *Model) SetSize(sizes ...int) {
	for i, size := range sizes {
		switch i {
		case 0:
			m.Viewer.SetWidth(size)
		case 1:
			m.Viewer.SetHeight(size)
		}
	}
}

// Write allows using this model as an argument to logging tools.
func (m *Model) Write(p []byte) (n int, err error) {
	errs := make([]error, 0, len(m.AdditionalWriters))
	for _, w := range m.AdditionalWriters {
		_, wErr := w.Write(p)
		errs = append(errs, wErr)
	}
	m.Viewer.AppendToText(string(p))
	return len(p), errors.Join(errs...)
}
