package main

import (
	"math/rand"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

func colorGen() func() uint8 {
	return func() uint8 {
		return uint8(rand.Intn(255))
	}
}

func newRGB() uint8 {
	return uint8(gobot.Rand(255))
}

func main() {
	rand.Seed(time.Now().Unix())

	gbot := gobot.NewGobot()

	edgeLen := uint8(75)

	adaptor := sphero.NewSpheroAdaptor("sphero", "/dev/rfcomm0")
	driver := sphero.NewSpheroDriver(adaptor, "sphero")

	work := func(edgeLen uint8) {
		//colorizer := colorGen()
		gobot.Every(3*time.Second, func() {
			driver.Roll(edgeLen, uint16(gobot.Rand(360)))
			driver.SetRGB(newRGB(), newRGB(), newRGB())
		})
	}

	worker := func() {
		work(edgeLen)
	}

	robot := gobot.NewRobot("sphero",
		[]gobot.Connection{adaptor},
		[]gobot.Device{driver},
		worker,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
