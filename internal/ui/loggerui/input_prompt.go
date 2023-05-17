package loggerui

import (
	"fmt"

	"github.com/wagoodman/go-partybus"

	"github.com/nextlinux/quill/internal/ui/tui/bubbles/prompt"
	"github.com/nextlinux/quill/quill/event/parser"
)

func (u *UI) handleInputPrompt(e partybus.Event) error {
	writer, err := parser.InputPrompt(e)
	if err != nil {
		return fmt.Errorf("unable to parse event: %+v", err)
	}

	model := prompt.New(writer)
	_, err = model.RunPrompt()
	return err
}
