package artist

import (
	"chroma-viz/hub"
	"chroma-viz/library"
	"log"
)

func InitialiseArtist(fileName string) {
	var err error
	conn = make(map[string]*library.Connection)
	chromaHub, err = hub.NewDataBase(10)
	if err != nil {
		log.Fatal(err)
	}

	conf, err = library.ImportConfig(fileName)
	if err != nil {
		log.Fatal(err)
	}

	hubConn = library.NewConnection("Hub", conf.HubAddr, conf.HubPort)
	hubConn.Connect()

	for _, c := range conf.Connections {
		conn[c.Name] = library.NewConnection(c.Name, c.Address, c.Port)
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
