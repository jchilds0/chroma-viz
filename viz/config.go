package viz

import (
	"chroma-viz/library/config"
	"chroma-viz/library/tcp"
	"fmt"
	"log"
)

var conf *config.Config

func AddConnection(name, conn_type, ip string, port int) error {
	if conn_type == "engine" {
		conn.eng = append(conn.eng, tcp.NewConnection(name, ip, port))
		return nil
	} else if conn_type == "preview" {
		conn.prev = append(conn.prev, tcp.NewConnection(name, ip, port))
		return nil
	}

	return fmt.Errorf("Unknown connection type %s", conn_type)
}

func InitialiseViz(configFile string) {
    var err error
    conf, err = config.ImportConfig(configFile)
    if err != nil {
        log.Fatal(err)
    }

	conn.hub = tcp.NewConnection("Hub", conf.HubAddr, conf.HubPort)
	conn.hub.Connect()

	for _, c := range conf.Connections {
		AddConnection(c.Name, c.Type, c.Address, c.Port)
	}
}

func CloseViz() {
	for _, c := range conn.eng {
		if c == nil {
			continue
		}

		if c.IsConnected() {
			c.CloseConn()
			log.Printf("Closed %s\n", c.Name)
		}
	}

	for _, c := range conn.prev {
		if c == nil {
			continue
		}

		if c.IsConnected() {
			c.CloseConn()
			log.Printf("Closed %s\n", c.Name)
		}
	}
}