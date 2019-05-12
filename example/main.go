package main

import (
	"github.com/ovrclk/dsky"
)

func main() {
	p := dsky.NewInteractivePrinter(nil, nil)
	dv := p.NewSection("").NewData().AsPane().
		Add("name(s)", "foo").
		Add("name(s)", "fomo")

	m := map[string]interface{}{
		"container": "akash-node",
		"namespace": "foo-bar",
	}
	details := dsky.NewSectionData().
		Add("app", "akash").
		Add("release", "us-west").
		Add("readiness", map[string]interface{}{
			"ttl":   "0sec",
			"state": "active",
		}).Add("info", m)
	dv.Add("details", details)

	p.Flush()
}
