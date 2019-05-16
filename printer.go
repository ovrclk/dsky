package dsky

import (
	"io"
)

type Printer interface {
	NewSection(id string) Section
	WithSection(Section) Printer
	Flush() error
	Log() Logger
}

func NewPrinter(m ModeType, stdout, errout io.Writer) (Printer, error) {
	switch m {
	case ModeTypeInteractive:
		return NewInteractiveMode(stdout, errout), nil
	case ModeTypeJSON:
		return NewJSONMode(stdout, errout), nil
	case ModeTypeShell:
		return NewShellMode(stdout, errout), nil
	default:
		return nil, ErrInvalidModeType{}
	}
}
