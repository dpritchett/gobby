package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

func main() {
	gbot := gobot.NewGobot()

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

	hello.AddCommand("left", func(params map[string]interface{}) interface{} {
		driver.Roll(75, uint16(270))
		return "moving left"
	})

	hello.AddCommand("right", func(params map[string]interface{}) interface{} {
		driver.Roll(75, uint16(90))
		return "moving right"
	})

	hello.AddCommand("forward", func(params map[string]interface{}) interface{} {
		driver.Roll(75, uint16(0))
		return "moving forward"
	})

	hello.AddCommand("back", func(params map[string]interface{}) interface{} {
		driver.Roll(75, uint16(0))
		return "moving back"
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

	// Starts the API server on default port 3000
	apiServer := api.NewAPI(gbot)
	apiServer.Start()

	demoFileService := func(w http.ResponseWriter, r *http.Request) {
		fileName := "./buttons.html"
		pageHTML, err := ioutil.ReadFile(fileName)

		if err != nil {
			log.Fatalf("error reading %v", fileName)
		}

		io.WriteString(w, string(pageHTML))
	}

	apiServer.Get("/demo/", demoFileService)

	gbot.Start()
}
