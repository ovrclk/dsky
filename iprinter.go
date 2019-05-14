package dsky

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/gosuri/uitable"
)

type InteractiveMode struct {
	sections []*Section
	common
}

func NewInteractiveMode(out, errout io.Writer) *InteractiveMode {
	if out == nil {
		out = os.Stdout
	}
	if errout == nil {
		errout = os.Stderr
	}
	m := &InteractiveMode{
		sections: make([]*Section, 0),
	}
	m.modeType = ModeTypeInteractive
	m.out = out
	m.errout = errout
	m.logger = NewInteractiveLogger(errout)
	return m
}

func (m *InteractiveMode) When(mtype ModeType, fn runF) Mode {
	if mtype == m.modeType {
		m.runners = append(m.runners, fn)
	}
	return m
}

func (i *InteractiveMode) Printer() Printer {
	return i
}

func (i *InteractiveMode) Log() Logger {
	return i.logger
}

func (i *InteractiveMode) NewSection(id string) *Section {
	s := &Section{ID: id}
	i.sections = append(i.sections, s)
	return s
}

func (i *InteractiveMode) WithSection(s *Section) Printer {
	i.sections = append(i.sections, s)
	return i
}

func (i *InteractiveMode) Flush() error {
	var buf bytes.Buffer
	for _, sec := range i.sections {
		if sec == nil {
			continue
		}
		if len(sec.ID) == 0 {
			return errors.New("dksy: section needs a title")
		}

		title := sec.ID
		if len(sec.Label) > 0 {
			title = sec.Label
		}
		buf.WriteString("\n")
		buf.WriteString(NewTitle(title).H1().String())
		buf.WriteString("\n")
		d, err := sec.Data.Marshal(i)
		if err != nil {
			return err
		}
		buf.Write(d)
		buf.WriteString("\n")
	}
	if _, err := fmt.Fprintln(i.out, buf.String()); err != nil {
		return err
	}
	return nil
}

func (i *InteractiveMode) MarshalSectionData(dv *SectionData) ([]byte, error) {
	return i.marshalSectionData(0, dv)
}

func (i *InteractiveMode) marshalSectionData(depth int, dv *SectionData) ([]byte, error) {
	if dv == nil {
		return nil, nil
	}
	switch dv.Style() {
	case SectionDataStylePane:
		return i.formatSDPane(depth, dv)
	case SectionDataStyleList:
		return i.formatSDList(depth, dv)
	default:
		return nil, fmt.Errorf("dsky: invalid section data style")
	}
	return nil, nil
}

func (i *InteractiveMode) formatSDPane(depth int, sectionData *SectionData) ([]byte, error) {
	wrapper := uitable.New()
	wrapper.Wrap = true
	// for each ID, create a row in the wrapper table
	for _, id := range sectionData.IDs() {
		// fetch the items for the id
		items := sectionData.Data()[id]
		if len(items) == 0 {
			return nil, nil
		}
		// use the label for the id as row name, if any
		label := id
		if l := sectionData.Label(id); len(l) > 0 {
			label = l
		}
		label = fmt.Sprintf("%s: ", label)
		// row items with the label as caption
		ritems := []interface{}{label}
		for _, v := range items {
			s, err := i.parsesd(v, depth)
			if err != nil {
				return nil, err
			}

			ritems = append(ritems, s)
		}
		// add the row to the wrapper table
		wrapper.AddRow(ritems...)
	}
	return wrapper.Bytes(), nil
}

func (i *InteractiveMode) formatSDList(depth int, sectionData *SectionData) ([]byte, error) {
	wrapper := uitable.New()
	wrapper.Wrap = true
	// create the header column with ids as the captions
	var headers []interface{}
	//var ids []interface{}

	var lc int // linecount
	// set the headers and the line count
	for _, id := range sectionData.IDs() {
		// use the label for the id as row name, if any
		label := id
		if l := sectionData.Label(id); len(l) > 0 {
			label = l
		}
		tl := NewTitle(label)
		switch depth {
		case 0:
			label = tl.H2().String()
		default:
			label = tl.H3().String()
		}
		headers = append(headers, label)
		if clc := len(sectionData.Data()[id]); clc > lc {
			lc = clc
		}
	}
	wrapper.AddRow(headers...)
	for rowIdx := 0; rowIdx < len(sectionData.Rows()); rowIdx++ {
		rowItems := make([]interface{}, 0)
		for cellIdx := 0; cellIdx < len(headers); cellIdx++ {
			if len(sectionData.Rows()[rowIdx]) > cellIdx {
				d, err := i.parsesd(sectionData.Rows()[rowIdx][cellIdx], depth)
				if err != nil {
					return nil, err
				}
				rowItems = append(rowItems, d)
				continue
			}
			// insert empty row
			rowItems = append(rowItems, "")
		}
		wrapper.AddRow(rowItems...)
	}
	return wrapper.Bytes(), nil
}

func (i *InteractiveMode) parsesd(v interface{}, depth int) (string, error) {
	var buf bytes.Buffer
	switch item := v.(type) {
	case *SectionData:
		d, err := i.marshalSectionData(depth+1, item)
		if err != nil {
			return "", err
		}
		buf.Write(d)
	case map[string]string:
		var lines []string
		var keys []string
		for k, _ := range item {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := item[k]
			lines = append(lines, fmt.Sprintf("%s: %s", k, v))
		}

		buf.WriteString(strings.Join(lines, " | "))
	default:
		buf.WriteString(fmt.Sprintf("%v", item))
	}
	return buf.String(), nil
}
