// Package app implements the root Bubble Tea interactive TUI model.
package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sujeevanferos/tercode/internal/agent/runtime"
	"github.com/sujeevanferos/tercode/internal/agent/session"
	"github.com/sujeevanferos/tercode/internal/events"
	"github.com/sujeevanferos/tercode/internal/tui/commands"
	"github.com/sujeevanferos/tercode/internal/tui/theme"
	"github.com/sujeevanferos/tercode/internal/workspace"
)

type responseMsg struct {
	text string
	err  error
}

// Model is the root Bubble Tea model for Tercode.
type Model struct {
	runtime   *runtime.Runtime
	workspace *workspace.Workspace
	events    *events.Bus
	session   *session.Session
	theme     *theme.Theme
	commands  *commands.Registry

	input   textinput.Model
	spinner spinner.Model

	history []string
	busy    bool
	width   int
	height  int
}

// New creates an initialized TUI Model.
func New(rt *runtime.Runtime, ws *workspace.Workspace, bus *events.Bus, s *session.Session) Model {
	ti := textinput.New()
	ti.Placeholder = "Ask Tercode anything, or type / for commands..."
	ti.Focus()
	ti.CharLimit = 2048
	ti.Width = 60

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))

	return Model{
		runtime:   rt,
		workspace: ws,
		events:    bus,
		session:   s,
		theme:     theme.Default(),
		commands:  commands.New(),
		input:     ti,
		spinner:   sp,
		history: []string{
			"Welcome to Tercode — Terminal Code",
			"Type a request to begin, or /help for commands.",
			"------------------------------------------------",
		},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = msg.Width - 10

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			val := strings.TrimSpace(m.input.Value())
			if val == "" {
				return m, nil
			}
			m.input.SetValue("")

			// Check for slash command
			if output, handled := m.commands.Execute(val); handled {
				m.history = append(m.history, "> "+val, output)
				return m, nil
			}

			// Agent request
			m.history = append(m.history, "> "+val)
			m.busy = true

			prompt := val
			sessionID := m.session.ID
			rt := m.runtime

			cmds = append(cmds, m.spinner.Tick, func() tea.Msg {
				resp, err := rt.Run(context.Background(), sessionID, prompt)
				return responseMsg{text: resp, err: err}
			})
			return m, tea.Batch(cmds...)
		}

	case responseMsg:
		m.busy = false
		if msg.err != nil {
			m.history = append(m.history, "Error: "+msg.err.Error())
		} else {
			m.history = append(m.history, msg.text)
		}
		return m, nil

	case spinner.TickMsg:
		if m.busy {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var sb strings.Builder

	// Header
	headerText := fmt.Sprintf(" Tercode | Workspace: %s | Mode: %s ", m.workspace.Name, m.session.Mode)
	sb.WriteString(m.theme.Header.Render(headerText) + "\n\n")

	// History scroll
	maxLines := 15
	if m.height > 10 {
		maxLines = m.height - 8
	}
	start := 0
	if len(m.history) > maxLines {
		start = len(m.history) - maxLines
	}
	for _, line := range m.history[start:] {
		if strings.HasPrefix(line, "> ") {
			sb.WriteString(m.theme.Prompt.Render(line) + "\n")
		} else if strings.HasPrefix(line, "Error:") {
			sb.WriteString(m.theme.Error.Render(line) + "\n")
		} else {
			sb.WriteString(line + "\n")
		}
	}

	// Status / Spinner
	if m.busy {
		sb.WriteString("\n" + m.spinner.View() + " Agent is working...\n")
	} else {
		sb.WriteString("\n")
	}

	// Prompt Input Box
	sb.WriteString(m.theme.Border.Render(m.input.View()) + "\n")

	return sb.String()
}
