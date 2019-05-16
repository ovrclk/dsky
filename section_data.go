package dsky

// ErrInvalidSectionDataID is an error that is
// returned when the SectionData identifier is invalid or missing
type ErrInvalidSectionDataID struct{}

// Error is the error message
func (e ErrInvalidSectionDataID) Error() string {
	return "dsky: invalid or missing SectionData Identifier"
}

// SectionDataMarshaler is the interface that discribes the
// marshaler of section data
type SectionDataMarshaler interface {
	MarshalSectionData(SectionData) ([]byte, error)
}

// SectionDataStyle is the style to render the section data
type SectionDataStyle uint

const (
	//SectionDataStylePane is used to set the section to render as "Pane"
	SectionDataStylePane SectionDataStyle = iota
	//SectionDataStylePane is used to set the section to render as "List"
	SectionDataStyleList
)

// SectionData describes the the data objects for the section
type SectionData interface {
	// Identifier returns the identifier for the section data
	Identifier() string

	// Add adds the items to section with the given id
	Add(id string, items ...interface{}) SectionData

	// Row return the rows (with columns)
	Rows() [][]interface{}

	// Tag returns the interface object associated with the tag
	Tag(tag string) interface{}

	// WithLabel adds a label to the section data which can be later used by the marshaler
	WithLabel(id, label string) SectionData

	// Label returs the string value associated with a label
	Label(id string) string

	// Marshal returns the formated bytes of the section data
	Marshal(SectionDataMarshaler) ([]byte, error)

	// Style returns the section data style
	Style() SectionDataStyle

	// AsPane sets the section data style to Pane (SectionDataStylePane)
	AsPane() SectionData

	// AsList sets the section data style to List (SectionDataStyleList)
	AsList() SectionData

	// WithTag is used to add extra information to the section data,
	// like Raw results while rendering JSON
	WithTag(tag string, msg interface{}) SectionData

	// Data returns a map of children with the field name as the key
	Data() map[string][]interface{}

	// IDs returns all non-hidden ids of the childern
	IDs() []string

	// Hide hides the ids from displaying
	Hide(ids ...string) SectionData
}

// NewSectionData returns a new instance of SectionData
func NewSectionData(id string) SectionData {
	return &sectionData{id: id}
}

type sectionData struct {
	id        string
	style     SectionDataStyle
	data      map[string][]interface{}
	ids       []string
	labels    map[string]string
	tags      map[string]interface{}
	hiddenIDs []string
}

func (d *sectionData) Marshal(m SectionDataMarshaler) ([]byte, error) {
	if len(d.id) == 0 {
		return nil, ErrInvalidSectionDataID{}
	}
	return m.MarshalSectionData(d)
}

func (d *sectionData) Identifier() string {
	return d.id
}

func (d *sectionData) Style() SectionDataStyle {
	return d.style
}

func (d *sectionData) AsPane() SectionData {
	d.style = SectionDataStylePane
	return d
}

func (d *sectionData) WithTag(tag string, msg interface{}) SectionData {
	if d.tags == nil {
		d.tags = make(map[string]interface{})
	}
	d.tags[tag] = msg
	return d
}

func (d *sectionData) AsList() SectionData {
	d.style = SectionDataStyleList
	return d
}

func (d *sectionData) Tag(tag string) interface{} {
	return d.tags[tag]
}

func (d *sectionData) Data() map[string][]interface{} {
	return d.data
}

func (d *sectionData) IDs() []string {
	var retids []string
outLoop:
	for _, id := range d.ids {
		for _, hid := range d.hiddenIDs {
			if id == hid {
				continue outLoop
			}
		}
		retids = append(retids, id)
	}
	return retids
}

func (d *sectionData) WithLabel(id, label string) SectionData {
	if d.labels == nil {
		d.labels = make(map[string]string)
	}
	d.labels[id] = label
	return d
}

func (d *sectionData) Label(id string) (l string) {
	return d.labels[id]
}

func (d *sectionData) Rows() [][]interface{} {
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
				val := d.Data()[secname][rowidx]
				if val == nil {
					rows[rowidx][colidx] = ""
					continue
				}
				rows[rowidx][colidx] = val
			}
		}
	}
	return rows
}
func (d *sectionData) Hide(ids ...string) SectionData {
	d.hiddenIDs = append(d.hiddenIDs, ids...)
	return d
}

func (d *sectionData) Add(id string, items ...interface{}) SectionData {
	if d.data == nil {
		d.data = make(map[string][]interface{})
	}
	for i := 0; i < len(items); i++ {
		switch sd := items[i].(type) {
		case *sectionData:
			if len(sd.Identifier()) == 0 {
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
