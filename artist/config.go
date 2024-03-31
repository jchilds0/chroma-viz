package artist

import (
	"chroma-viz/hub"
	"chroma-viz/library/config"
	"chroma-viz/library/props"
	"chroma-viz/library/tcp"
	"log"
)

func InitialiseArtist() {
    var err error
	conn = make(map[string]*tcp.Connection)

    conf, err = config.ImportConfig("./artist/conf.json")
    if err != nil {
        log.Fatal(err)
    }

	geo := []string{"rect", "text", "circle", "graph", "image", "ticker", "clock"}
	geo_count = []int{10, 10, 10, 10, 10, 10, 10}
	geoms = make(map[int]*geom, len(geo))

	index := 1
	for i, name := range geo {
		geoms[props.StringToProp[name]] = newGeom(index, geo_count[i])
		index += geo_count[i]
	}

	hub.GenerateTemplateHub(geo, geo_count, "artist/artist.json")
	hub.ImportArchive("artist/artist.json")
	go hub.StartHub(conf.HubPort, -1)

    artistHub := tcp.NewConnection("Hub", conf.HubAddr, conf.HubPort)
	artistHub.Connect()

	for _, c := range conf.Connections {
        conn[c.Name] = tcp.NewConnection(c.Name, c.Address, c.Port)
	}

}

func CloseConn() {
	for name, c := range conn {
		if c.IsConnected() {
			c.CloseConn()
			log.Printf("Closed %s\n", name)
		}
	}
}
