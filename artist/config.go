package artist

import (
	"chroma-viz/hub"
	"chroma-viz/library/config"
	"chroma-viz/library/props"
	"chroma-viz/library/tcp"
	"log"
)

func InitialiseArtist(fileName string) {
	var err error
	conn = make(map[string]*tcp.Connection)
	chromaHub = hub.NewDataBase(10)

	conf, err = config.ImportConfig(fileName)
	if err != nil {
		log.Fatal(err)
	}

	geo := []string{"rect", "text", "circle", "graph", "image", "ticker", "clock"}
	geo_count := 10
	geoms = make(map[int]*geom, len(geo))

	index := 1
	for _, name := range geo {
		geoms[props.StringToProp[name]] = newGeom(index, geo_count)
		index += geo_count
	}

	chromaHub.AddTemplate(0, "left_to_right", "", "left_to_right")

	var total int
	for i := range geo {
		for j := 0; j < geo_count; j++ {
			chromaHub.AddGeometry(0, total, geo[i])
			total++
		}
	}
	go hub.StartHub(chromaHub, conf.HubPort)

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
