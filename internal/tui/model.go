package tui

import (
	"cibo/internal/pipelines"
	"cibo/internal/web"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
The goal of this module is to act as the primary control point for the Terminal UI
layer, powered by the Bubble Tea library. This library functions via a model (state),
update, view system that should be familiar to anyone whose worked with Redux or a similar
UI library built on the ELM architecture, with the "Update()" cycle acting as a state reducer

Docs on Bubble Tea can be found here: github.com/charmbracelet/bubbletea
*/

const maxLogMessages = 20

const (
	InfoLog LogType = iota
	SuccessLog
	ErrorLog
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

type processSuccessMsg struct {
	recordCount int
	filePath    string
	logs        []string
}

type processErrorMsg struct {
	err error
}

type webUILaunchedMsg struct {
	url string
}

type LogType int

type LogEntry struct {
	Type    LogType
	Message string
}

// Defines the state space of the Bubbletea TUI
type model struct {
	pipelines          *pipelines.Pipelines
	focusIndex         int
	inputs             []textinput.Model
	spinner            spinner.Model
	loading            bool
	err                error
	processingComplete bool
	resultFilePath     string
	launchUIPrompt     textinput.Model
	logs               []LogEntry
	width              int
	height             int
}

func (m model) reset() model {
	for i := range m.inputs {
		m.inputs[i].Reset()
	}
	m.launchUIPrompt.Reset()
	m.focusIndex = 0
	m.inputs[0].Focus()
	m.loading = false
	m.processingComplete = false
	m.resultFilePath = ""
	m.err = nil

	return m
}

// Defines the initial state of the TUI
func NewModel(pipelines *pipelines.Pipelines, initialLogs []string) model {
	m := model{
		inputs:    make([]textinput.Model, 3),
		logs:      make([]LogEntry, 0),
		pipelines: pipelines,
	}

	for _, msg := range initialLogs {
		m.logInfo(msg)
	}
	m.logInfo("Welcome! Enter a stock ticker to begin analysis.")

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
			t.Width = 5
		case 1:
			t.Placeholder = "YYYY-MM-DD"
			t.Prompt = "Start Date:   "
			t.CharLimit = 10
			t.Width = 10
		case 2:
			t.Placeholder = "YYYY-MM-DD"
			t.Prompt = "End Date:     "
			t.CharLimit = 10
			t.Width = 10
		}
		m.inputs[i] = t
	}

	launchPrompt := textinput.New()
	launchPrompt.Prompt = "Launch the web UI to view the chart? (y/n) "
	launchPrompt.Cursor.Style = cursorStyle
	launchPrompt.CharLimit = 1
	launchPrompt.Width = 3
	m.launchUIPrompt = launchPrompt

	return m
}

func (m *model) log(entry LogEntry) {
	m.logs = append(m.logs, entry)
	if len(m.logs) > maxLogMessages {
		m.logs = m.logs[1:]
	}
}

func (m *model) logInfo(msg string) {
	m.log(LogEntry{Type: InfoLog, Message: msg})
}

func (m *model) logSuccess(msg string) {
	m.log(LogEntry{Type: SuccessLog, Message: msg})
}

func (m *model) logError(msg string) {
	m.log(LogEntry{Type: ErrorLog, Message: msg})
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) launchWebUICmd() tea.Msg {
	listener, url, err := web.PrepareListener()
	if err != nil {
		return processErrorMsg{err: err}
	}
	web.StartNonBlocking(listener, m.resultFilePath)
	return webUILaunchedMsg{url: url}
}

func (m model) processDataCmd() tea.Msg {
	ticker := m.inputs[0].Value()
	startDate := m.inputs[1].Value()
	endDate := m.inputs[2].Value()
	m.logInfo(fmt.Sprintf("Fetching data for %s...", ticker))
	lynchFairValueInputs := pipelines.LynchFairValueInputs{
		Ticker:    ticker,
		StartDate: startDate,
		EndDate:   endDate,
	}

	lynchFairValueOutputs, err := m.pipelines.LynchFairValue.RunPipeline(lynchFairValueInputs)
	if err != nil {
		return processErrorMsg{err: err}
	}

	return processSuccessMsg{
		recordCount: lynchFairValueOutputs.RecordCount,
		filePath:    lynchFairValueOutputs.FilePath,
		logs:        lynchFairValueOutputs.Logs,
	}
}

// --- Bubbletea Update ---
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case processSuccessMsg:
		m.loading = false
		m.processingComplete = true
		m.resultFilePath = msg.filePath
		for _, log := range msg.logs {
			m.logSuccess(log)
		}
		cmd = m.launchUIPrompt.Focus()
		return m, cmd

	case processErrorMsg:
		m.loading = false
		m.err = msg.err
		m.logError(fmt.Sprintf("Error: %v", msg.err))

		return m, nil

	case webUILaunchedMsg:
		m.logInfo(fmt.Sprintf("Web server running at %s. Open browser to:", msg.url))
		return m, nil

	case tea.KeyMsg:
		if m.processingComplete {
			if msg.String() == "enter" {
				answer := strings.ToLower(m.launchUIPrompt.Value())
				if answer == "y" {
					cmds = append(cmds, m.launchWebUICmd)
				}
				m = m.reset()
				cmds = append(cmds, textinput.Blink)
				return m, tea.Batch(cmds...)
			}
			m.launchUIPrompt, cmd = m.launchUIPrompt.Update(msg)
			return m, cmd
		}

		// --- Handle main form ---
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.loading = true
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
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds = append(cmds, m.inputs[i].Focus())
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}
			}
			return m, tea.Batch(cmds...)
		}
	}

	for i := range m.inputs {
		if m.inputs[i].Focused() {
			m.inputs[i], cmd = m.inputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// --- Bubbletea View ---
func (m model) View() string {
	paneWidth := m.width / 2

	// --- Left Pane Form
	var form strings.Builder
	for i := range m.inputs {
		form.WriteString(m.inputs[i].View() + "\n")
	}
	form.WriteString("\n")

	if m.loading {
		form.WriteString(m.spinner.View() + " Processing...")
	} else if m.processingComplete {
		form.WriteString(m.launchUIPrompt.View())
	} else {
		button := "[ Submit ]"
		if m.focusIndex == len(m.inputs) {
			button = focusedStyle.Render("[ Submit ]")
		}
		form.WriteString(button)
	}
	form.WriteString("\n\n" + helpStyle.Render("tab: next field • q: quit"))

	formPaneStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63"))

	// --- Right Pane Logs
	logHeader := lipgloss.NewStyle().
		Padding(0, 1).
		Bold(true).
		Background(lipgloss.Color("63")).
		Render("Logs")

	availableHeight := max(m.height-5, 0)

	start := 0
	if len(m.logs) > availableHeight {
		start = len(m.logs) - availableHeight
	}
	var styledLogs []string
	for _, logEntry := range m.logs[start:] {
		switch logEntry.Type {
		case SuccessLog:
			styledLogs = append(styledLogs, successMessageStyle.Render(logEntry.Message))
		case ErrorLog:
			styledLogs = append(styledLogs, errorMessageStyle.Render(logEntry.Message))
		default:
			styledLogs = append(styledLogs, logEntry.Message)
		}
	}
	logContent := strings.Join(styledLogs, "\n")

	logPaneStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		formPaneStyle.
			Width(paneWidth-10).
			Height(m.height-10).
			Render(form.String()),
		logPaneStyle.Width(m.width-paneWidth-10).
			Height(m.height-10).
			Render(logHeader+"\n"+logContent),
	)
}
