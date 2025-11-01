package tui

import (
	"cibo/internal/pipelines"
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type mockFairValuePipeline struct {
	shouldReturnErr bool
	outputToReturn  *pipelines.LynchFairValueOutputs
	wasCalled       bool
	receivedTicker  string
}

func (m *mockFairValuePipeline) RunPipeline(input pipelines.LynchFairValueInputs) (*pipelines.LynchFairValueOutputs, error) {
	m.wasCalled = true
	m.receivedTicker = input.Ticker
	if m.shouldReturnErr {
		return nil, errors.New("mock pipeline error")
	}
	return m.outputToReturn, nil
}

func dispatch(m model, msg tea.Msg) (model, tea.Cmd) {
	newModel, cmd := m.Update(msg)
	return newModel.(model), cmd
}

func processCmd(m model, cmd tea.Cmd) model {
	if cmd == nil {
		return m
	}
	msg := cmd()
	if batchMsg, ok := msg.(tea.BatchMsg); ok {
		for _, subCmd := range batchMsg {
			m = processCmd(m, subCmd)
		}
		return m
	}
	m, _ = dispatch(m, msg)
	return m
}

// Helper to check if any log message contains a specific substring.
func containsLog(logs []LogEntry, substr string) bool {
	for _, log := range logs {
		if strings.Contains(log.Message, substr) {
			return true
		}
	}
	return false
}

// Given a user who fills out the form and submits it successfully,
// verify that the TUI calls the pipeline and logs the success messages.
func TestTUI_HappyPath_Success(t *testing.T) {
	mockPipeline := &mockFairValuePipeline{
		outputToReturn: &pipelines.LynchFairValueOutputs{
			RecordCount: 13,
			FilePath:    "NVDA.parquet",
			Logs:        []string{"Successfully wrote 13 combined records to Parquet file"},
		},
	}
	rootPipelines := &pipelines.Pipelines{LynchFairValue: mockPipeline}
	m := NewModel(rootPipelines, nil)
	var cmd tea.Cmd

	// Simulate typing a ticker.
	for _, char := range "NVDA" {
		m, _ = dispatch(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
	}

	// Verify input text.
	if m.inputs[0].Value() != "NVDA" {
		t.Errorf("Expected input value 'NVDA', got '%s'", m.inputs[0].Value())
	}

	// Simulate navigating to the "Submit" button.
	for i := 0; i < 3; i++ {
		m, _ = dispatch(m, tea.KeyMsg{Type: tea.KeyTab})
	}

	// Verify focus is on the submit button.
	if m.focusIndex != len(m.inputs) {
		t.Errorf("Expected focus on submit button (%d), got %d", len(m.inputs), m.focusIndex)
	}

	// Simulate pressing "Enter" to submit.
	m, cmd = dispatch(m, tea.KeyMsg{Type: tea.KeyEnter})

	// Verify the model is now in a loading state.
	if !m.loading {
		t.Fatal("Expected model.loading to be true, but it was false")
	}

	// Process the command returned by the submission (which runs the pipeline).
	m = processCmd(m, cmd)

	// Verify the pipeline was called correctly.
	if !mockPipeline.wasCalled {
		t.Error("Expected pipeline's RunPipeline method to be called")
	}
	if mockPipeline.receivedTicker != "NVDA" {
		t.Errorf("Expected pipeline to receive ticker 'NVDA', got '%s'", mockPipeline.receivedTicker)
	}

	// Verify the model is no longer loading and is now in the completion state.
	if m.loading {
		t.Error("Expected model.loading to be false after processing")
	}
	if !m.processingComplete {
		t.Error("Expected model.processingComplete to be true after success")
	}

	// Verify the logs contain the expected success messages.
	expectedLogText := "Successfully wrote 13 combined records to Parquet file"
	if !containsLog(m.logs, expectedLogText) {
		t.Errorf("Expected logs to contain '%s', but they did not. Logs: %v", expectedLogText, m.logs)
	}
}

// Given a user who submits the form but the pipeline returns an error,
// verify that the TUI logs the error message correctly.
func TestTUI_PipelineError(t *testing.T) {
	mockPipeline := &mockFairValuePipeline{shouldReturnErr: true}
	rootPipelines := &pipelines.Pipelines{LynchFairValue: mockPipeline}
	m := NewModel(rootPipelines, nil)
	var cmd tea.Cmd

	// Simulate typing a ticker and submitting the form.
	m.inputs[0].SetValue("TSLA")
	m.focusIndex = len(m.inputs) // Move focus to submit button
	m, cmd = dispatch(m, tea.KeyMsg{Type: tea.KeyEnter})

	// Verify loading state.
	if !m.loading {
		t.Fatal("Expected model.loading to be true")
	}

	// Process the command.
	m = processCmd(m, cmd)

	// Verify the loading state is now false.
	if m.loading {
		t.Error("Expected model.loading to be false after error")
	}

	// Verify the model is NOT in the completion state.
	if m.processingComplete {
		t.Error("Expected model.processingComplete to be false after an error")
	}

	// Verify the logs contain the error message.
	expectedErrorText := "Error: mock pipeline error"
	if !containsLog(m.logs, expectedErrorText) {
		t.Errorf("Expected logs to contain '%s', but they did not. Logs: %v", expectedErrorText, m.logs)
	}
}

// Given a set of initial logs passed to the constructor,
// verify that they are present in the model's state immediately.
func TestTUI_InitialLogs(t *testing.T) {
	initialLogs := []string{
		"Using mock API.",
		"Config loaded successfully.",
	}

	m := NewModel(nil, initialLogs)

	if len(m.logs) < len(initialLogs) {
		t.Fatalf("Expected at least %d initial logs, but got %d", len(initialLogs), len(m.logs))
	}

	if !containsLog(m.logs, "Using mock API.") {
		t.Error("Expected logs to contain 'Using mock API.'")
	}
	if !containsLog(m.logs, "Config loaded successfully.") {
		t.Error("Expected logs to contain 'Config loaded successfully.'")
	}
}

// todo test log styles in the Logs pane
