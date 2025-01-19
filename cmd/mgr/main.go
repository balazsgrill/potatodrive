package main

import (
	"github.com/go-ole/go-ole"
	"github.com/integrii/flaggy"
)

func main() {
	err := ole.RoInitialize(1)
	if err != nil {
		panic(err)
	}
	listcmd := flaggy.NewSubcommand("list")
	unregcmd := flaggy.NewSubcommand("unreg")
	var unregid string
	unregcmd.AddPositionalValue(&unregid, "ID", 1, true, "The id of the sync root to unregister")
	flaggy.AttachSubcommand(listcmd, 1)
	flaggy.AttachSubcommand(unregcmd, 1)
	flaggy.Parse()

	if listcmd.Used {
		list()
	}
	if unregcmd.Used {
		unreg(unregid)
	}
}
