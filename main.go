package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strconv"
	"strings"
)

func main() {
	p := tea.NewProgram(gridModel())
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
	styles     *Styles
}

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

type model struct {
	grid     [9][9]spaces            // board grid user need to fill
	cursorX  int                     // which X square position cursor is pointing at
	cursorY  int                     // which Y square position cursor is pointing at
	selected map[coordinate]struct{} // which square are selected
}

func InputStyle() *Styles {
	s := new(Styles)
	s.BorderColor = "36"
	s.InputField = lipgloss.NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(0).
		Width(3).
		Height(1).
		AlignHorizontal(lipgloss.Center)

	return s
}

func gridModel() model {
	// Render new sudoku board (a squares grid)
	var squares [9][9]spaces
	for i := 0; i < len(squares); i++ {
		for j := 0; j < len(squares[i]); j++ {
			style := InputStyle()
			input := textinput.New()
			input.Prompt = ""
			input.CharLimit = 1

			squares[i][j].value = 0
			squares[i][j].input = input
			squares[i][j].styles = style
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

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "down", "left", "right":
			userInput := msg.String()

			// Clear cursor traces
			m.grid[m.cursorY][m.cursorX].input.Blur()
			m.grid[m.cursorY][m.cursorX].input.TextStyle = noStyle

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

			// New Cursor direction
			cursorDirection := m.grid[m.cursorY][m.cursorX].input.Focus()
			m.grid[m.cursorY][m.cursorX].input.TextStyle = focusedStyle

			return m, cursorDirection

		default:
			if numericInput, err := strconv.Atoi(msg.String()); err == nil && numericInput != 0 {
				var newSpaceValue tea.Cmd
				m.grid[m.cursorY][m.cursorX].value = numericInput
				m.grid[m.cursorY][m.cursorX].input, newSpaceValue = m.grid[m.cursorY][m.cursorX].input.Update(msg)

				return m, newSpaceValue
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	var gridBuilder string
	var headerBuilder strings.Builder
	var footerBuilder strings.Builder

	headerBuilder.WriteString("\nGoSudoku!!\n")

	for columnIndex := 0; columnIndex < len(m.grid); columnIndex++ {
		var row string
		for rowIndex := 0; rowIndex < len(m.grid[columnIndex]); rowIndex++ {
			row = lipgloss.JoinHorizontal(
				lipgloss.Left,
				row,
				m.grid[columnIndex][rowIndex].styles.InputField.Render(
					m.grid[columnIndex][rowIndex].input.View()))

			if rowIndex == 2 || rowIndex == 5 {
				row += "  "
			}
		}
		gridBuilder = lipgloss.JoinVertical(
			lipgloss.Left,
			gridBuilder,
			row,
		)
		if columnIndex == 2 || columnIndex == 5 {
			gridBuilder = lipgloss.JoinVertical(
				lipgloss.Left,
				gridBuilder,
				" ",
			)
		}
	}
	// Debug
	footerBuilder.WriteString(helpStyle.Render(fmt.Sprintf("\n\nx: %v / y: %v", m.cursorX, m.cursorY)))
	footerBuilder.WriteString(helpStyle.Render(fmt.Sprintf("\nValue: %v", m.grid[m.cursorY][m.cursorX].value)))

	footerBuilder.WriteString(helpStyle.Render("\nPress q to quit.\n"))

	return headerBuilder.String() + gridBuilder + footerBuilder.String()
}
