package parser

import (
	"github.com/wagoodman/go-partybus"

	"github.com/nextlinux/quill/quill/event"
)

func Exit(e partybus.Event) error {
	return checkEventType(e.Type, event.Exit)
}
