package main

import (
	"pandora/pkg/view"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(view.MenuReady, view.MenuExit)
}
