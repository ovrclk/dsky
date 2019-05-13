package main

import (
	"flag"
	"fmt"

	"github.com/ovrclk/dsky"
)

type group struct {
	Seq       int
	name      string
	req       map[string]string
	resources map[string]string
}

type bid struct{ group, price, provider string }

var (
	mode   string
	groups = []group{
		{1, "west", map[string]string{"region": "us-west"}, map[string]string{"cpu": "200", "memory": "2Gb", "price": "100", "count": "2"}},
		{2, "east", map[string]string{"region": "us-east"}, map[string]string{"cpu": "800", "memory": "4Gb", "price": "80", "count": "5"}},
	}

	bids = []bid{
		{"1", "9", "9859d7fe1b4a0052b0c62a0b55b1881eccff14c5cfbe71d1b2393a265f6c3c92"},
		{"1", "12", "3b58798de9e03a23517d52bd2f9e1c0be5ded0589196ce480972757397f82260"},
	}
)

func main() {
	flag.StringVar(&mode, "m", string(dsky.ModeTypeInteractive), "mode")
	flag.Parse()

	mode, err := dsky.NewMode(dsky.ModeType(mode), nil, nil)
	printer := mode.Printer()

	if err != nil {
		panic(err)
	}

	log := printer.Log().WithModule("broadcast")
	log.Warn("requesting deployment for group(s): westcoast")
	log.Info("request accepted, deployment created with id: 81d79c80c4c7eb202cfd4846bb8e5328110cb299e9864674836b9fec6b536285")
	log.WithModule("keys").Error("Unable to select a default key.\nToo many keys are stored locally to pick a default, a key is selected as the default only when there is a single key present.\nFound 3 keys instead of 1")
	log.WithAction(dsky.LogActionDone).Info("request deployment for group(s): westcoast")

	gd := printer.NewSection("Groups").WithLabel("Deployment Status").NewData().AsList()
	for _, g := range groups {
		gd.Add("Seq", g.Seq)
		gd.Add("name", g.name)
	}

	data := printer.NewSection("Deployment").WithLabel("Deployment Status").NewData().AsPane().
		Add("DeployID", "f258a119d288989fa50471ebf7e8635d8fe93412077d77b89e5adc359202e144").WithLabel("DeployID", "Deployment ID").
		Add("Services", map[string]string{
			"web": "http://example.com",
			"app": "http://example.com",
		})

	gp := dsky.NewSectionData("").AsList()
	for _, g := range groups {
		gp.
			Add("Sequence", g.Seq).
			Add("Name", g.name).
			Add("Requirements", g.req).
			Add("Resources", g.resources)
	}
	data.Add("Groups", gp)

	fp := dsky.NewSectionData("").AsList()
	for _, b := range bids {
		fp.
			Add("Group", b.group).WithLabel("Group", "Group (Sequence)").
			Add("Price", b.price).
			Add("Provider", b.provider)
	}
	data.Add("Bids", fp)

	data.WithTag("raw", groups)

	printer.Flush()

	mode.When(dsky.ModeTypeInteractive, func() error {
		fmt.Println("(this will only show in interactive mode)")
		return nil
	}).Run()
}
