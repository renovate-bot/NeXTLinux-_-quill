package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wagoodman/go-partybus"

	"github.com/nextlinux/quill/internal/log"
	"github.com/nextlinux/quill/internal/ui/tui/bubbles/prompt"
	"github.com/nextlinux/quill/quill/event/parser"
)

func (m *UI) handleInputPrompt(e partybus.Event) (tea.Model, tea.Cmd) {
	writer, err := parser.InputPrompt(e)
	if err != nil {
		log.Warnf("unable to parse event: %+v", err)
		return m, nil
	}

	teaModel := prompt.New(writer)

	m.uiElements = append(m.uiElements, teaModel)

	return m, teaModel.Init()
}
