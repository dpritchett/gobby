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
	bot.HostFileAtRoute("./static/buttons.html", "/demo/")
	log.Fatal(bot.Start())
}
