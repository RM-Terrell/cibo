package tui

import (
	"fmt"
	"stats_engine/api"
	"stats_engine/io"
	"stats_engine/parse"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xitongsys/parquet-go-source/local"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	successMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)
	errorMessageStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)

// --- Custom Messages for data processing pipeline ---

type processSuccessMsg struct {
	recordCount int
	fileName    string
}

type processErrorMsg struct {
	err error
}

// --- Bubbletea Model ---

type model struct {
	focusIndex     int
	inputs         []textinput.Model
	spinner        spinner.Model
	loadingMessage string
	apiClient      *api.Client
	successMessage string
	err            error
}

// Defines the initial state of the TUI
func NewModel(client *api.Client) model {
	m := model{
		inputs:    make([]textinput.Model, 3),
		apiClient: client,
	}

	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "AAPL"
			t.Focus()
			t.Prompt = "Stock Ticker: "
			t.CharLimit = 5
		case 1:
			t.Placeholder = "YYYY-MM-DD"
			t.Prompt = "Start Date:   "
			t.CharLimit = 10
		case 2:
			t.Placeholder = "YYYY-MM-DD"
			t.Prompt = "End Date:     "
			t.CharLimit = 10
		}
		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// --- Functionality Commands ---
// Run the entire stock price data pipeline.
func (m model) processDataCmd() tea.Msg {
	jsonData, err := m.apiClient.FetchDailyStockData(m.inputs[0].Value())
	if err != nil {
		return processErrorMsg{err: fmt.Errorf("API fetch failed: %w", err)}
	}

	records, err := parse.ParseToFlat(jsonData, true)
	if err != nil {
		return processErrorMsg{err: fmt.Errorf("parsing failed: %w", err)}
	}

	fileName := fmt.Sprintf("%s.parquet", m.inputs[0].Value())
	fw, err := local.NewLocalFileWriter(fileName)
	if err != nil {
		return processErrorMsg{err: fmt.Errorf("failed to create file '%s': %w", fileName, err)}
	}
	defer fw.Close()

	if err := io.WriteToParquet(records, fw); err != nil {
		return processErrorMsg{err: fmt.Errorf("failed to write parquet data: %w", err)}
	}

	return processSuccessMsg{recordCount: len(records), fileName: fileName}
}

// --- Bubbletea Update ---
// Update handles messages and updates the model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If we are in a loading state, we only listen for spinner ticks and results.
	if m.loadingMessage != "" {
		switch msg := msg.(type) {
		case spinner.TickMsg:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		case processSuccessMsg:
			m.loadingMessage = ""
			m.successMessage = fmt.Sprintf("Success! Wrote %d records to %s", msg.recordCount, msg.fileName)
			return m, tea.Quit // Quit after success
		case processErrorMsg:
			m.loadingMessage = ""
			m.err = msg.err
			return m, tea.Quit // Quit after error
		default:
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		// Intended to handle navigation between inputs and submitting the form
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Enter keystroke on the "Submit" button
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.loadingMessage = "Processing data..."

				return m, tea.Batch(m.spinner.Tick, m.processDataCmd)
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		if i == m.focusIndex {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}
	}
	return tea.Batch(cmds...)
}

// --- Bubbletea View ---
func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nAn error occurred: %s\n\n", errorMessageStyle.Render(m.err.Error()))
	}

	if m.successMessage != "" {
		return fmt.Sprintf("\n%s\n\n", successMessageStyle.Render(m.successMessage))
	}

	if m.loadingMessage != "" {
		return fmt.Sprintf("\n   %s %s... press q to quit\n\n", m.spinner.View(), m.loadingMessage)
	}

	var b strings.Builder
	fmt.Fprintln(&b, "Enter stock information for analysis.")
	fmt.Fprintln(&b)
	for i := range m.inputs {
		fmt.Fprintln(&b, m.inputs[i].View())
	}
	button := "[ Submit ]"
	if m.focusIndex == len(m.inputs) {
		button = focusedStyle.Render("[ Submit ]")
	}
	fmt.Fprintf(&b, "\n%s\n\n", button)
	fmt.Fprintln(&b, helpStyle.Render("tab: next field • q: quit"))
	return b.String()
}
