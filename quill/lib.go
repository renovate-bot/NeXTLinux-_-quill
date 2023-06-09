package quill

import (
	"github.com/wagoodman/go-partybus"

	"github.com/nextlinux/gologger"
	"github.com/nextlinux/quill/internal/bus"
	"github.com/nextlinux/quill/internal/log"
)

// SetLogger sets the logger object used for all logging calls.
func SetLogger(logger logger.Logger) {
	log.Set(logger)
}

// SetBus sets the event bus for all library bus publish events onto (in-library subscriptions are not allowed).
func SetBus(b *partybus.Bus) {
	bus.SetPublisher(b)
}
