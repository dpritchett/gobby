package main

import (
	"flag"
	"log"

	"github.com/dpritchett/gobby/webby"
)

// pass a different bluetooth FD as needed: webby /dev/rfcomm1
func getSpheroFd() string {
	deviceFd := "/dev/rfcomm0"

	flag.Parse()
	if flag.Arg(0) != "" {
		log.Println(flag.Arg(0))
		deviceFd = flag.Arg(0)
	}

	return deviceFd
}

func main() {
	bot := webby.NewWebby(getSpheroFd())

	bot.AddCommand("green", func() { bot.Driver.SetRGB(0, 255, 0) }, "turning green")
	bot.AddCommand("red", func() { bot.Driver.SetRGB(255, 0, 0) }, "turning red")
	bot.AddCommand("blue", func() { bot.Driver.SetRGB(0, 0, 255) }, "turning blue")

	bot.AddCommand("forward", func() { bot.Roll(75, 0) }, "moving forward")
	bot.AddCommand("right", func() { bot.Roll(75, 90) }, "moving right")
	bot.AddCommand("back", func() { bot.Roll(75, 180) }, "moving back")
	bot.AddCommand("left", func() { bot.Roll(75, 270) }, "moving left")

	bot.AddCommand("wiggle", bot.Wiggle, "wiggle wiggle wiggle wiggle wiggle wiggle wiggle!")

	bot.AddCommand("new_color", bot.SetRandColor, "new color!")

	bot.HostFileAtRoute("./static/buttons.html", "/demo/")
	log.Fatal(bot.Start())
}
