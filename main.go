package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle.Copy()
)

type coordinate struct {
	positionX int
	positionY int
}

type spaces struct {
	value      int
	input      textinput.Model
	coordinate coordinate
}

type model struct {
	grid     [9][9]spaces            // board grid user need to fill
	cursorX  int                     // which X square position cursor is pointing at
	cursorY  int                     // which Y square position cursor is pointing at
	selected map[coordinate]struct{} // which square are selected
}

func initialModel() model {
	// Render new sudoku board (a squares grid)
	var squares [9][9]spaces
	for i := 0; i < len(squares); i++ {
		for j := 0; j < len(squares[i]); j++ {
			// Setup space input
			input := textinput.New()
			input.Placeholder = "0"
			input.CharLimit = 1
			input.Prompt = "|"

			// Declare default values
			squares[i][j].value = 0
			squares[i][j].input = input
		}
	}

	return model{
		grid:     squares,
		selected: make(map[coordinate]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		// Keys to exit program
		case "ctrl+c", "q":
			return m, tea.Quit

		// Handle movement keys in board grid
		case "up", "down", "left", "right":
			userInput := msg.String()

			if userInput == "up" && m.cursorY > 0 {
				m.cursorY--
			}

			if userInput == "down" && m.cursorY < len(m.grid)-1 {
				m.cursorY++
			}

			if userInput == "left" && m.cursorX > 0 {
				m.cursorX--
			}

			if userInput == "right" && m.cursorX < len(m.grid)-1 {
				m.cursorX++
			}

			// Catch coordinates and focus on current square
			commands := make([]tea.Cmd, len(m.grid))
			for i := 0; i <= len(m.grid)-1; i++ {
				for j := 0; j <= len(m.grid[i])-1; j++ {
					if i == m.cursorY && j == m.cursorX {
						// Set focused state
						commands[i] = m.grid[i][j].input.Focus()
						m.grid[i][j].input.PromptStyle = focusedStyle
						m.grid[i][j].input.TextStyle = focusedStyle
						continue
					}

					// Remove focused state
					m.grid[i][j].input.Blur()
					m.grid[i][j].input.PromptStyle = noStyle
					m.grid[i][j].input.TextStyle = noStyle
				}
			}
			return m, tea.Batch(commands...)
		}
	}

	// Handle character input
	commands := m.updateInputs(msg)

	return m, commands
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	var commands = make([]tea.Cmd, len(m.grid))

	// Update the input value
	for i := 0; i <= len(m.grid)-1; i++ {
		for j := 0; j <= len(m.grid[i])-1; j++ {
			m.grid[i][j].input, commands[i] = m.grid[i][j].input.Update(msg)
			//TODO: Upadate the "value" variable on square
		}
	}

	return tea.Batch(commands...)
}

func (m model) View() string {
	// Build spaces strings
	var builder strings.Builder

	// Iterate over grid
	for columnIndex := 0; columnIndex < len(m.grid); columnIndex++ {
		for rowIndex := 0; rowIndex < len(m.grid[columnIndex]); rowIndex++ {
			// Build text input on space
			builder.WriteString(m.grid[columnIndex][rowIndex].input.View())
		}
		builder.WriteString("\n")
	}

	// The footer
	builder.WriteString(fmt.Sprintf("\nx: %v / y: %v", m.cursorX, m.cursorY))
	builder.WriteString(helpStyle.Render("\nPress q to quit.\n"))

	// Send the UI for rendering
	return builder.String()
}
