package tui

import (
	"cibo/internal/pipelines"
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// --- Mock Implementations for Testing

// mockFairValuePipeline is a mock implementation of the FairValuePipeline interface.
// It allows us to control the pipeline's behavior and verify its usage during tests.
type mockFairValuePipeline struct {
	shouldReturnErr bool
	outputToReturn  *pipelines.LynchFairValueOutputs

	// Fields to inspect the mock's state after the test
	wasCalled      bool
	receivedTicker string
}

// RunPipeline implements the FairValuePipeline interface for our mock.
func (m *mockFairValuePipeline) RunPipeline(input pipelines.LynchFairValueInputs) (*pipelines.LynchFairValueOutputs, error) {
	m.wasCalled = true
	m.receivedTicker = input.Ticker

	if m.shouldReturnErr {
		return nil, errors.New("mock pipeline error")
	}
	return m.outputToReturn, nil
}

// Helper that calls the Update method and performs the necessary type assertion,
// reducing boilerplate in the test when Update calls have to be run
func dispatch(m model, msg tea.Msg) (model, tea.Cmd) {
	newModel, cmd := m.Update(msg)
	return newModel.(model), cmd
}

// Given a user who fills out the fair value form and submits it successfully,
// verify that the TUI correctly calls the pipeline and displays the success message.
func TestTUI_HappyPath_Success(t *testing.T) {
	mockPipeline := &mockFairValuePipeline{
		outputToReturn: &pipelines.LynchFairValueOutputs{
			RecordCount: 13,
			FileName:    "NVDA.parquet",
		},
	}

	rootPipelines := &pipelines.Pipelines{
		LynchFairValue: mockPipeline,
	}

	m := NewModel(rootPipelines)

	// Simulate the user typing a ticker into the first input.
	for _, char := range "NVDA" {
		m, _ = dispatch(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
	}

	// Verify the input text is correct.
	if m.inputs[0].Value() != "NVDA" {
		t.Errorf("Expected input value to be 'NVDA', but got '%s'", m.inputs[0].Value())
	}

	// Simulate navigating to the "Submit" button.
	m, _ = dispatch(m, tea.KeyMsg{Type: tea.KeyTab}) // Focus moves to input 1
	m, _ = dispatch(m, tea.KeyMsg{Type: tea.KeyTab}) // Focus moves to input 2
	m, _ = dispatch(m, tea.KeyMsg{Type: tea.KeyTab}) // Focus moves to Submit button

	// Verify focus is on the submit button.
	if m.focusIndex != len(m.inputs) {
		t.Errorf("Expected focusIndex to be on the submit button (%d), but got %d", len(m.inputs), m.focusIndex)
	}

	// Simulate pressing "Enter" to submit the form.
	m, _ = dispatch(m, tea.KeyMsg{Type: tea.KeyEnter})

	// Verify the model is now in a loading state.
	if m.loadingMessage == "" {
		t.Fatal("Expected model to be in a loading state, but loadingMessage was empty")
	}

	// todo this is kinda cursed, handle batching instead of grey box testing
	resultMsg := m.processDataCmd()

	if !mockPipeline.wasCalled {
		t.Error("Expected the pipeline's RunPipeline method to be called, but it was not.")
	}
	if mockPipeline.receivedTicker != "NVDA" {
		t.Errorf("Expected pipeline to receive ticker 'NVDA', but got '%s'", mockPipeline.receivedTicker)
	}

	m, _ = dispatch(m, resultMsg)

	if m.successMessage == "" {
		t.Fatal("Expected a success message, but it was empty")
	}
	expectedSuccessText := "Success! Wrote 13 records to NVDA.parquet"
	if !strings.Contains(m.successMessage, expectedSuccessText) {
		t.Errorf("Expected success message to contain '%s', but got '%s'", expectedSuccessText, m.successMessage)
	}
}
