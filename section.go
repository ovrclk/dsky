package dsky

type SectionDataStyle uint

const (
	SectionDataStylePane SectionDataStyle = iota
	SectionDataStyleList
)

type SectionDataFormatter interface {
	FormatSectionData(*SectionData) []byte
}

type Section struct {
	ID   string
	Data *SectionData
}

func (s *Section) NewData() *SectionData {
	s.Data = NewSectionData()
	return s.Data
}

func NewSectionData() *SectionData {
	return &SectionData{}
}

type SectionData struct {
	style  SectionDataStyle
	data   map[string][]interface{}
	labels []string
}

func (d *SectionData) Style() SectionDataStyle {
	return d.style
}

func (d *SectionData) AsPane() *SectionData {
	d.style = SectionDataStylePane
	return d
}

func (d *SectionData) AsList() *SectionData {
	d.style = SectionDataStyleList
	return d
}

func (d *SectionData) Data() map[string][]interface{} {
	return d.data
}

func (d *SectionData) Labels() []string {
	return d.labels
}

func (d *SectionData) Add(label string, dv ...interface{}) *SectionData {
	if d.data == nil {
		d.data = make(map[string][]interface{})
	}
	d.data[label] = append(d.data[label], dv...)
	// avoid duplicate labels
	for _, l := range d.labels {
		if l == label {
			return d
		}
	}
	d.labels = append(d.labels, label)
	return d
}
