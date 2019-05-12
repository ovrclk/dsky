package dsky

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/huandu/xstrings"
)

var (
	ShellVarPrefix = "akash"
)

type ShellMode struct {
	sections []*Section
	common
}

func NewShellMode(out, errout io.Writer) *ShellMode {
	if out == nil {
		out = os.Stdout
	}
	if errout == nil {
		errout = os.Stderr
	}
	s := &ShellMode{
		sections: make([]*Section, 0),
	}
	s.modeType = ModeTypeShell
	s.out = out
	s.errout = errout
	s.logger = &shellLogger{}
	return s
}

func (i *ShellMode) Printer() Printer {
	return i
}

func (m *ShellMode) When(mtype ModeType, fn runF) Mode {
	if mtype == m.modeType {
		m.runners = append(m.runners, fn)
	}
	return m
}

func (i *ShellMode) NewSection(id string) *Section {
	s := &Section{ID: id}
	i.sections = append(i.sections, s)
	return s
}

func (i *ShellMode) WithSection(s *Section) Printer {
	i.sections = append(i.sections, s)
	return i
}

func (i *ShellMode) Flush() error {
	var buf bytes.Buffer
	for _, sec := range i.sections {
		if sec == nil {
			continue
		}
		d, err := sec.Data.Marshal(i)
		if err != nil {
			return err
		}
		buf.Write(d)
	}
	if _, err := fmt.Fprintln(i.out, buf.String()); err != nil {
		return err
	}
	return nil
}

func (s *ShellMode) MarshalSectionData(sdata *SectionData) ([]byte, error) {
	var buf bytes.Buffer
	data, err := s.marshalSectionData(sdata)
	if err != nil {
		return nil, err
	}

	// vars in associateive arrays to declare with array name as the key
	arrs := make(map[string][]string)
	// vars not part of an array
	nvars := make([]string, 0)

	for _, evar := range data {
		// if this var is an array items, add it to the array declaration
		if len(evar.arrkey) > 0 {
			v := fmt.Sprintf("%s[%s]=%q", evar.name(), evar.arrKey(), evar.val)
			arrs[evar.name()] = append(arrs[evar.name()], v)
			continue
		}
		nvars = append(nvars, fmt.Sprintf("%s=%q", evar.name(), evar.val))
	}

	// render associate array declars first
	for aname, _ := range arrs {
		buf.WriteString(fmt.Sprintf("declare -A %s", aname))
		buf.WriteString("\n")
	}
	// array items
	for _, arr := range arrs {
		for _, v := range arr {
			buf.WriteString(v)
			buf.WriteString("\n")
		}
	}
	// non array vars
	for _, v := range nvars {
		buf.WriteString(v)
		buf.WriteString("\n")
	}

	return buf.Bytes(), nil
}

type envvar struct {
	varname []string // variable name in parts
	val     string   // value
	arrkey  string   // associative array key the variable belongs to
}

func (e envvar) arrKey() string {
	return xstrings.ToSnakeCase(e.arrkey)
}

func (e envvar) name() string {
	name := append([]string{ShellVarPrefix}, e.varname...)
	return xstrings.ToSnakeCase(strings.Join(name, " "))
}

func fmtVarName(v string) string {
	return fmt.Sprintf("%s_%s", ShellVarPrefix, strings.ToLower(v))
}

func (i *ShellMode) marshalSectionData(sectionData *SectionData) ([]envvar, error) {
	var result []envvar
	for rowIdx, row := range sectionData.Rows() {
		for colIdx, cell := range row {
			secname := sectionData.IDs()[colIdx]
			switch item := cell.(type) {
			case string:
				vp := []string{sectionData.Identifier(), strconv.Itoa(rowIdx)}
				result = append(result, envvar{varname: vp, val: item, arrkey: secname})
			case map[string]string:
				var keys []string
				for k, _ := range item {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					v := item[k]
					vp := []string{sectionData.Identifier(), strconv.Itoa(rowIdx), secname}
					result = append(result, envvar{varname: vp, val: v, arrkey: k})
				}
			case *SectionData:
				res, err := i.marshalSectionData(item)
				if err != nil {
					return nil, err
				}
				for _, ev := range res {
					vp := append([]string{sectionData.Identifier(), strconv.Itoa(rowIdx)}, ev.varname...)
					result = append(result, envvar{varname: vp, val: ev.val, arrkey: ev.arrkey})
				}
			}
		}
	}
	return result, nil
}
