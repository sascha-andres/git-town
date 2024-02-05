package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func TextDisplay(title, text string, inputs TestInput) (bool, error) {
	model := textDisplayModel{
		colors: createColors(),
		status: StatusActive,
		text:   text,
		title:  title,
	}
	program := tea.NewProgram(model)
	// TODO: extract into helper function.
	if len(inputs) > 0 {
		go func() {
			for _, input := range inputs {
				program.Send(input)
			}
		}()
	}
	dialogResult, err := program.Run()
	if err != nil {
		return false, err
	}
	result := dialogResult.(textDisplayModel) //nolint:forcetypeassert
	return result.status == StatusAborted, nil
}

type textDisplayModel struct {
	colors dialogColors
	status status
	text   string
	title  string
}

func (self textDisplayModel) Init() tea.Cmd {
	return nil
}

func (self textDisplayModel) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) { //nolint:ireturn
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type { //nolint:exhaustive
		case tea.KeyEnter:
			self.status = StatusDone
			return self, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			self.status = StatusAborted
			return self, tea.Quit
		}
		switch msg.String() {
		case "o":
			self.status = StatusDone
			return self, tea.Quit
		case "q":
			self.status = StatusAborted
			return self, tea.Quit
		}
	case error:
		panic(msg.Error())
	}
	return self, cmd
}

func (self textDisplayModel) View() string {
	if self.status != StatusActive {
		return ""
	}
	result := strings.Builder{}
	result.WriteRune('\n')
	result.WriteString(self.colors.Title.Styled(self.title))
	result.WriteRune('\n')
	result.WriteString(self.text)
	result.WriteString("\n\n  ")
	// accept
	result.WriteString(self.colors.HelpKey.Styled("o"))
	result.WriteString(self.colors.Help.Styled("/"))
	result.WriteString(self.colors.HelpKey.Styled("enter"))
	result.WriteString(self.colors.Help.Styled(" continue   "))
	// abort
	result.WriteString(self.colors.HelpKey.Styled("q"))
	result.WriteString(self.colors.Help.Styled("/"))
	result.WriteString(self.colors.HelpKey.Styled("esc"))
	result.WriteString(self.colors.Help.Styled("/"))
	result.WriteString(self.colors.HelpKey.Styled("ctrl-c"))
	result.WriteString(self.colors.Help.Styled(" abort"))
	return result.String()
}