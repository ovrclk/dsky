package dsky

import (
	"fmt"
	"io"
)

type ModeType string

const (
	ModeTypeInteractive ModeType = "interactive"
	ModeTypeShell                = "shell"
	ModeTypeJSON                 = "json"
)

type runF func() error

type ErrInvalidModeType struct{}

func (e ErrInvalidModeType) Error() string {
	return fmt.Sprintf("dsky: invalid mode type")
}

type Mode interface {
	// Type must return the type of Mode
	Type() ModeType

	// Printer must return the printer for the mode
	Printer() Printer

	// When registers the events to run for the
	// current Mode when Run is invoked. It returns the current Mode
	When(ModeType, runF) Mode

	// Run runs the functions
	Run() error

	// Ask returns an Asker
	Ask() Asker

	// IsInteractive returns true if the current ModeType is ModeTypeInteractive
	IsInteractive() bool
}

func NewMode(m ModeType, stdout, errout io.Writer) (Mode, error) {
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

type common struct {
	out      io.Writer
	errout   io.Writer
	modeType ModeType
	runners  []runF
	logger   Logger
	asker    Asker
}

func (c common) Log() Logger {
	return c.logger
}

func (m common) Run() error {
	for i := 0; i < len(m.runners); i++ {
		var fn runF
		fn, m.runners = m.runners[0], m.runners[1:]
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (m common) Type() ModeType {
	return m.modeType
}

func (m common) Ask() Asker {
	if m.asker == nil {
		m.asker = NewInteractiveAsker(m.Type())
	}
	return m.asker
}

func (m common) IsInteractive() bool {
	return m.modeType == ModeTypeInteractive
}
