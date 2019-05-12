package dsky

import (
	"bytes"
	"strings"
)

// TitleUnderliner is the underline character for the title
var TitleUnderliner = "="

// Title is a UI component that renders a title. Title implements Component interface.
type Title struct {
	text        string
	uliner      string
	isUnderLine bool
	isCaps      bool
}

func NewTitle(text string) *Title {
	return &Title{text: text, uliner: TitleUnderliner}
}

func (t *Title) WithUnderliner(u string) *Title {
	t.isUnderLine = true
	t.uliner = u
	return t
}

func (t *Title) H1() *Title {
	return t.WithUnderliner("=")
}

func (t *Title) H2() *Title {
	return t.WithUnderliner("-")
}

func (t *Title) H3() *Title {
	t.isCaps = true
	return t
}

// String returns the formated string of the title
func (t *Title) Bytes() []byte {
	var buf bytes.Buffer

	switch {
	case t.isUnderLine:
		buf.WriteString(t.text + "\n")
		for i := 0; i < len(t.text); i++ {
			buf.Write([]byte(t.uliner))
		}
		buf.WriteString("\n")
	case t.isCaps:
		buf.WriteString(strings.ToUpper(t.text))
	}

	return buf.Bytes()
}

func (t *Title) String() string {
	return string(t.Bytes())
}
