/*
Package bublog to show your logs in TUI and duplicate them to other places.
*/
package bublog

import (
	"errors"
	"io"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	LogView           *viewport.Model
	AdditionalWriters []io.Writer
	LogContent        string
}

func New(additionalWriters []io.Writer) *Model {
	return &Model{&viewport.Model{}, additionalWriters, ""}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (viewport.Model, tea.Cmd) {
	res, cmd := m.LogView.Update(msg)
	return res, cmd
}

func (m *Model) View() string {
	return m.LogView.View()
}

// SetSize accepts up to two arguments. First fill be Width, second - Height.
func (m *Model) SetSize(sizes ...int) {
	for i, size := range sizes {
		switch i {
		case 0:
			m.LogView.Width = size
		case 1:
			m.LogView.Height = size
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
	m.LogContent += string(p)
	m.LogView.SetContent(m.fitLogs())
	m.LogView.GotoBottom()
	return len(p), errors.Join(errs...)
}

func (m *Model) fitLogs() string {
	fitTo := m.LogView.Width - m.LogView.Style.GetHorizontalFrameSize()
	return lipgloss.NewStyle().Width(fitTo).Render(m.LogContent)
}
