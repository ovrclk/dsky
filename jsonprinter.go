package dsky

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/huandu/xstrings"
)

type JSONMode struct {
	sections []Section
	common
}

func NewJSONMode(out, errout io.Writer) *JSONMode {
	if out == nil {
		out = os.Stdout
	}
	if errout == nil {
		errout = os.Stderr
	}
	m := &JSONMode{
		sections: make([]Section, 0),
	}
	m.modeType = ModeTypeJSON
	m.out = out
	m.errout = errout
	m.logger = &jsonLogger{}
	return m
}

func (m *JSONMode) When(mtype ModeType, fn runF) Mode {
	if mtype == m.modeType {
		m.runners = append(m.runners, fn)
	}
	return m
}

func (i *JSONMode) Printer() Printer {
	return i
}

func (i *JSONMode) NewSection(id string) Section {
	s := NewSection(id)
	i.sections = append(i.sections, s)
	return s
}

func (i *JSONMode) WithSection(s Section) Printer {
	i.sections = append(i.sections, s)
	return i
}

func (i *JSONMode) Flush() error {
	var buf bytes.Buffer
	for _, sec := range i.sections {
		if sec == nil {
			continue
		}
		b, err := sec.Data().Marshal(i)
		if err != nil {
			return err
		}
		buf.Write(b)
		if raw := sec.Data().Tag("raw"); raw != nil {
			dat := map[string]interface{}{"raw": raw}
			d, err := json.Marshal(dat)
			if err != nil {
				return err
			}
			buf.Write(d)
		}

	}
	fmt.Fprintln(i.out, buf.String())
	return nil
}

func (i *JSONMode) MarshalSectionData(sectionData SectionData) ([]byte, error) {
	d, err := i.marshalSectionData(sectionData)
	if err != nil {
		return nil, err
	}
	id := xstrings.ToSnakeCase(sectionData.Identifier())
	return json.Marshal(map[string]interface{}{id: d})
}

func (i *JSONMode) marshalSectionData(sectionData SectionData) (interface{}, error) {
	var recc int // record count
	// get the records count
	for _, id := range sectionData.IDs() {
		if c := len(sectionData.Data()[id]); c > recc {
			recc = c
		}
	}
	recs := make([]map[string]interface{}, recc)
	for rowidx, row := range sectionData.Rows() {
		recs[rowidx] = make(map[string]interface{})
		for colidx, secdata := range row {
			if v, ok := secdata.(SectionData); ok {
				d, err := i.marshalSectionData(v)
				if err != nil {
					return nil, err
				}
				secdata = d
			}
			secname := xstrings.ToSnakeCase(sectionData.IDs()[colidx])
			recs[rowidx][secname] = secdata
		}
	}
	return recs, nil
}
