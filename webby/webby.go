package webby

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

type Webby struct {
	Gobot     *gobot.Gobot
	Sphero    *gobot.Robot
	ApiServer *api.API
}

func NewWebby(spheroFd string) *Webby {
	ss := &Webby{}

	ss.Gobot = gobot.NewGobot()

	adaptor := sphero.NewSpheroAdaptor("sphero", spheroFd)
	driver := sphero.NewSpheroDriver(adaptor, "sphero")

	spheroid := gobot.NewRobot("sphero",
		[]gobot.Connection{adaptor},
		[]gobot.Device{driver},
	)

	ss.Sphero = ss.Gobot.AddRobot(spheroid)

	// Accessible via http://localhost:3000/robots/sphero/commands/turn_blue
	ss.Sphero.AddCommand("blue", func(params map[string]interface{}) interface{} {
		driver.SetRGB(0, 0, 255)
		return "turning blue"
	})

	ss.Sphero.AddCommand("green", func(params map[string]interface{}) interface{} {
		driver.SetRGB(0, 255, 0)
		return "turning green"
	})

	ss.Sphero.AddCommand("red", func(params map[string]interface{}) interface{} {
		driver.SetRGB(255, 0, 0)
		return "turning red"
	})

	ss.Sphero.AddCommand("left", func(params map[string]interface{}) interface{} {
		driver.Roll(75, uint16(270))
		return "moving left"
	})

	ss.Sphero.AddCommand("right", func(params map[string]interface{}) interface{} {
		driver.Roll(75, uint16(90))
		return "moving right"
	})

	ss.Sphero.AddCommand("forward", func(params map[string]interface{}) interface{} {
		driver.Roll(75, uint16(0))
		return "moving forward"
	})

	ss.Sphero.AddCommand("back", func(params map[string]interface{}) interface{} {
		driver.Roll(75, uint16(180))
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

	ss.Sphero.AddCommand("wiggle", func(params map[string]interface{}) interface{} {
		times := []time.Duration{0, 2, 4}

		for _, t := range times {
			gobot.After(t*time.Second, wiggle)
		}

		return "wiggle wiggle wiggle wiggle wiggle wiggle wiggle"
	})

	// Starts the API server on default port 3000
	ss.ApiServer = api.NewAPI(ss.Gobot)
	ss.ApiServer.Start()

	return ss
}

func (ss *Webby) HostFileAtRoute(filePath, route string) {
	demoFileService := func(w http.ResponseWriter, r *http.Request) {
		pageHTML, err := ioutil.ReadFile(filePath)

		if err != nil {
			log.Fatalf("error reading %v", filePath)
		}

		io.WriteString(w, string(pageHTML))
	}

	ss.ApiServer.Get(route, demoFileService)
}

func (ss *Webby) Start() (errs []error) {
	return ss.Gobot.Start()
}
