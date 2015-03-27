package webby

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

type Webby struct {
	Gobot     *gobot.Gobot
	Sphero    *gobot.Robot
	APIServer *api.API
	Driver    *sphero.SpheroDriver
}

func NewWebby(spheroFd string) *Webby {
	ss := &Webby{}

	ss.Gobot = gobot.NewGobot()

	adaptor := sphero.NewSpheroAdaptor("sphero", spheroFd)
	ss.Driver = sphero.NewSpheroDriver(adaptor, "sphero")

	spheroid := gobot.NewRobot("sphero",
		[]gobot.Connection{adaptor},
		[]gobot.Device{ss.Driver},
	)

	ss.Sphero = ss.Gobot.AddRobot(spheroid)

	// Starts the API server on default port 3000
	ss.APIServer = api.NewAPI(ss.Gobot)

	return ss
}

func (ss *Webby) AddCommand(name string, impl func(), desc string) {
	ss.Sphero.AddCommand(name, func(params map[string]interface{}) interface{} {
		impl()
		return desc
	})
}

func (ss *Webby) HostFileAtRoute(filePath, route string) {
	demoFileService := func(w http.ResponseWriter, r *http.Request) {
		pageHTML, err := ioutil.ReadFile(filePath)

		if err != nil {
			log.Fatalf("error reading %v", filePath)
		}

		io.WriteString(w, string(pageHTML))
	}

	ss.APIServer.Get(route, demoFileService)
}

func (ss *Webby) Start() (errs []error) {
	ss.APIServer.Start()
	return ss.Gobot.Start()
}

func (ss *Webby) SetRandColor() {
	rgb := func() uint8 {
		return uint8(gobot.Rand(255))
	}

	ss.Driver.SetRGB(rgb(), rgb(), rgb())
}

func (ss *Webby) Wiggle() {
	ss.SetRandColor()
	ss.Roll(50+ss.Rand(50), ss.Rand(360))
}

// exposing some useful sphero utils to use for new commands
func (ss *Webby) Rand(n int) int {
	return gobot.Rand(n)
}

func (ss *Webby) Roll(distance, heading int) {
	ss.Driver.Roll(uint8(distance), uint16(heading))
}
