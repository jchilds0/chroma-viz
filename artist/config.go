package artist

import (
	"chroma-viz/hub"
	"chroma-viz/library/config"
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

	hubConn = tcp.NewConnection("Hub", conf.HubAddr, conf.HubPort)
	hubConn.Connect()

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
