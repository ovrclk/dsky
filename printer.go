package dsky

type Printer interface {
	AddSection(*Section)
	Flush() error
	SectionDataFormatter
}
