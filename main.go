package main

import (
	"chroma-viz/artist"
	"chroma-viz/viz"
	"encoding/json"
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var profile = flag.String("profile", "", "write profile to file")
var mode = flag.String("mode", "", "chroma mode (artist | viz)")

func main() {
    flag.Parse()
    if *profile != "" {
        f, err := os.Create(*profile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    var app *gtk.Application
    var err error
    if *mode == "artist" {
        artist.InitConnections()
        defer artist.CloseConn()

        artist.AddConnection("Preview", "127.0.0.1", 6100)

        app, err = gtk.ApplicationNew("app.chroma.artist", glib.APPLICATION_FLAGS_NONE)
        if err != nil {
            log.Print(err)
        }

        app.Connect("activate", artist.ArtistGui)
    } else if *mode == "viz" {
        defer viz.CloseConn()

        buf, err := os.ReadFile("conf.json")
        if err != nil {
            log.Fatalf("Error reading config file (%s)", err)
        }

        conf := NewConfig()
        err = json.Unmarshal(buf, conf)
        if err != nil {
            log.Fatalf("Error parsing config file (%s)", err)
        }

        SendConnToViz(conf)

        app, err = gtk.ApplicationNew("app.chroma.viz", glib.APPLICATION_FLAGS_NONE)
        if err != nil {
            log.Print(err)
        }

        app.Connect("activate", viz.VizGui)
    } else {
        flag.PrintDefaults()
        return
    }

    app.Run([]string{})
}

type Conn struct {
    Name    string
    Address string
    Port    int
    Type    string
}

type Config struct {
    HubAddr       string
    HubPort       int
    Connections   []Conn
}

func NewConfig() *Config {
    conf := &Config{}
    return conf
}

func SendConnToViz(conf *Config) {
    viz.AddGraphicsHub(conf.HubAddr, conf.HubPort)

    for _, c := range conf.Connections {
        viz.AddConnection(c.Name, c.Type, c.Address, c.Port)
    }
}
