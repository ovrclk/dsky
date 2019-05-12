package dsky

import (
	fc "github.com/fatih/color"
)

var Color = NewColor()

type color struct {
	Success, Notice, Failure, Hi, Normal *fc.Color
}

func NewColor() *color {
	return &color{
		Success: fc.New(fc.FgHiGreen),
		Notice:  fc.New(fc.FgHiYellow),
		Failure: fc.New(fc.FgHiRed),
		Hi:      fc.New(fc.FgHiWhite),
		Normal:  fc.New(fc.FgWhite),
	}
}
