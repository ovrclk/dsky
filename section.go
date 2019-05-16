package dsky

// Section represent a data section in the printer
type Section interface {
	// WithID set the section id with the provided string and returns the section
	WithID(id string) Section

	// ID returns the id of the section
	ID() string

	// NewData create a new SectionData and set the section's data with the newly created data
	NewData() SectionData

	// WithData attaches the data to the section and returns the section
	WithData(SectionData) Section

	// Data returns the section's data
	Data() SectionData

	// WithLabel sets the section's label with the given string
	WithLabel(label string) Section

	// Label returns the section's label
	Label() string
}

// NewSection creates and returns a new instance of a section
func NewSection(id string) Section {
	return &section{id: id}
}

type section struct {
	id    string
	data  SectionData
	label string
}

func (s *section) WithID(string) Section {
	return s
}

func (s *section) ID() string {
	return s.id
}

func (s *section) NewData() SectionData {
	s.data = NewSectionData(s.id)
	return s.data
}

func (s *section) WithData(data SectionData) Section {
	s.data = data
	return s
}

func (s *section) Data() SectionData {
	return s.data
}

func (s *section) WithLabel(l string) Section {
	s.label = l
	return s
}

func (s *section) Label() string {
	return s.label
}
