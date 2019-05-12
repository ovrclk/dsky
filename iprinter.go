package dsky

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/gosuri/uitable"
)

type InteractivePrinter struct {
	out      io.Writer
	errout   io.Writer
	sections []*Section
}

func NewInteractivePrinter(out, errout io.Writer) *InteractivePrinter {
	if out == nil {
		out = os.Stdout
	}
	if errout == nil {
		errout = os.Stderr
	}
	return &InteractivePrinter{
		out:      out,
		errout:   errout,
		sections: make([]*Section, 0),
	}
}

func (i *InteractivePrinter) NewSection(id string) *Section {
	s := &Section{ID: id}
	i.sections = append(i.sections, s)
	return s
}

func (i *InteractivePrinter) WithSection(s *Section) *InteractivePrinter {
	i.sections = append(i.sections, s)
	return i
}

func (i *InteractivePrinter) Flush() error {
	var buf bytes.Buffer
	for _, sec := range i.sections {
		buf.Write(i.FormatSectionData(sec.Data))
	}
	fmt.Fprint(i.out, buf.String())
	return nil
}

func (i *InteractivePrinter) FormatSectionData(dv *SectionData) []byte {
	switch dv.Style() {
	case SectionDataStylePane:
		return i.formatSDPane(dv)
	case SectionDataStyleList:
		return i.formatDataViewList(dv)
	}
	return nil
}

func (i *InteractivePrinter) formatSDPane(dv *SectionData) []byte {
	t := uitable.New()
	t.Wrap = true
	for _, label := range dv.Labels() {
		items := dv.Data()[label]
		if len(items) == 0 {
			return nil
		}
		cellItems := []interface{}{label}
		vt := uitable.New()
		vt.Wrap = true
		for _, value := range items {
			var buf bytes.Buffer
			switch v := value.(type) {
			case *SectionData:
				buf.Write(i.FormatSectionData(v))
			case map[string]interface{}:
				mt := uitable.New()
				mt.Wrap = true
				for key, val := range v {
					switch mv := val.(type) {
					case *SectionData:
						mt.AddRow(key, string(i.FormatSectionData(mv)))
					default:
						mt.AddRow(key, fmt.Sprintf("%v", mv))
					}
				}
				buf.WriteString(mt.String())
			default:
				buf.WriteString(fmt.Sprintf("%v", v))
			}
			vt.AddRow(buf.String())
		}
		cellItems = append(cellItems, vt.String())
		t.AddRow(cellItems...)
	}
	return t.Bytes()
}

func (i *InteractivePrinter) formatDataViewList(dv *SectionData) []byte {
	return nil
}
