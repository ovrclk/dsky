package dsky

type SectionDataStyle uint

const (
	SectionDataStylePane SectionDataStyle = iota
	SectionDataStyleList
)

type ErrInvalidSectionDataID struct{}

func (e ErrInvalidSectionDataID) Error() string {
	return "dsky: invalid or missing SectionData Identifier"
}

type SectionDataMarshaler interface {
	MarshalSectionData(*SectionData) ([]byte, error)
}

type Section struct {
	ID    string
	Data  *SectionData
	Label string
}

func (s *Section) NewData() *SectionData {
	s.Data = NewSectionData(s.ID)
	return s.Data
}

func (s *Section) WithLabel(l string) *Section {
	s.Label = l
	return s
}

func NewSectionData(id string) *SectionData {
	return &SectionData{id: id}
}

type SectionData struct {
	id     string
	style  SectionDataStyle
	data   map[string][]interface{}
	ids    []string
	labels map[string]string
	tags   map[string]interface{}
}

func (d *SectionData) Marshal(m SectionDataMarshaler) ([]byte, error) {
	if len(d.id) == 0 {
		return nil, ErrInvalidSectionDataID{}
	}
	return m.MarshalSectionData(d)
}

func (d *SectionData) Identifier() string {
	return d.id
}

func (d *SectionData) Style() SectionDataStyle {
	return d.style
}

func (d *SectionData) AsPane() *SectionData {
	d.style = SectionDataStylePane
	return d
}

func (d *SectionData) WithTag(tag string, msg interface{}) *SectionData {
	if d.tags == nil {
		d.tags = make(map[string]interface{})
	}
	d.tags[tag] = msg
	return d
}

func (d *SectionData) AsList() *SectionData {
	d.style = SectionDataStyleList
	return d
}

func (d *SectionData) Tag(tag string) interface{} {
	return d.tags[tag]
}

func (d *SectionData) Data() map[string][]interface{} {
	return d.data
}

func (d *SectionData) IDs() []string {
	return d.ids
}

func (d *SectionData) WithLabel(id, label string) *SectionData {
	if d.labels == nil {
		d.labels = make(map[string]string)
	}
	d.labels[id] = label
	return d
}

func (d *SectionData) Label(id string) (l string) {
	return d.labels[id]
}

func (d *SectionData) Rows() [][]interface{} {
	var rowc int // record count
	for _, id := range d.IDs() {
		if c := len(d.Data()[id]); c > rowc {
			rowc = c
		}
	}
	rows := make([][]interface{}, rowc)
	for rowidx := 0; rowidx < rowc; rowidx++ {
		rows[rowidx] = make([]interface{}, len(d.IDs()))
		for colidx := 0; colidx < len(d.IDs()); colidx++ {
			secname := d.IDs()[colidx]
			if len(d.Data()[secname]) > rowidx {
				rows[rowidx][colidx] = d.Data()[secname][rowidx]
			}
		}
	}
	return rows
}

func (d *SectionData) Add(id string, items ...interface{}) *SectionData {
	if d.data == nil {
		d.data = make(map[string][]interface{})
	}
	for i := 0; i < len(items); i++ {
		switch sd := items[i].(type) {
		case *SectionData:
			if len(sd.id) == 0 {
				sd.id = id
			}
			items[i] = sd
		}
	}
	d.data[id] = append(d.data[id], items...)
	// avoid duplicate ids
	for _, l := range d.ids {
		if l == id {
			return d
		}
	}
	d.ids = append(d.ids, id)
	return d
}
