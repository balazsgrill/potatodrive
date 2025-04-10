package main

import (
	"github.com/balazsgrill/potatodrive/bindings"
	"github.com/balazsgrill/potatodrive/ui"
	"github.com/balazsgrill/potatodrive/ui/configui"
	"github.com/rs/zerolog/log"
)

func main() {
	uicontext := ui.NewUIContext("0.0.0")

	configProvider := bindings.NewRegistryConfigWriter(log.Logger, "SOFTWARE\\PotatoDrive")

	configui.Create(uicontext, configProvider).Run()
}
