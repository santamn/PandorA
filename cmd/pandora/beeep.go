package main

import "github.com/gen2brain/beeep"

const (
	pathToAppIcon = "icon/PandorAmicro.png"
)

func alert(text string) error {
	return beeep.Alert("PandorA Error", text, pathToAppIcon)
}

func notify(text string) error {
	return beeep.Notify("PandorA", text, pathToAppIcon)
}
