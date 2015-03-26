package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

func main() {
	gbot := gobot.NewGobot()

	// Starts the API server on default port 3000
	api.NewAPI(gbot).Start()

	adaptor := sphero.NewSpheroAdaptor("sphero", "/dev/rfcomm0")
	driver := sphero.NewSpheroDriver(adaptor, "sphero")

	spheroid := gobot.NewRobot("sphero",
		[]gobot.Connection{adaptor},
		[]gobot.Device{driver},
	)

	hello := gbot.AddRobot(spheroid)

	// Accessible via http://localhost:3000/robots/sphero/commands/turn_blue
	hello.AddCommand("blue", func(params map[string]interface{}) interface{} {
		driver.SetRGB(0, 0, 255)
		return "turning blue"
	})

	hello.AddCommand("green", func(params map[string]interface{}) interface{} {
		driver.SetRGB(0, 255, 0)
		return "turning green"
	})

	hello.AddCommand("red", func(params map[string]interface{}) interface{} {
		driver.SetRGB(255, 0, 0)
		return "turning red"
	})

	rgb := func() uint8 {
		return uint8(gobot.Rand(255))
	}

	setRandColor := func() {
		driver.SetRGB(rgb(), rgb(), rgb())
	}

	wiggle := func() {
		setRandColor()
		driver.Roll(75, uint16(gobot.Rand(360)))
	}

	hello.AddCommand("wiggle", func(params map[string]interface{}) interface{} {
		wiggle()

		gobot.After(2*time.Second, func() {
			wiggle()
		})

		gobot.After(4*time.Second, func() {
			wiggle()
		})

		return "wiggle wiggle wiggle wiggle wiggle wiggle wiggle"
	})

	gbot.Start()
}
