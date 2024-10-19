package main

import (
	"chroma-viz/library/hub"
	"flag"
	"fmt"
	"log"

	"github.com/gin-contrib/pprof"
)

var profile = flag.String("profile", "", "write profile to file")
var port = flag.Int("port", 9000, "chroma hub port")
var username = flag.String("u", "", "SQL Database username")
var password = flag.String("p", "", "SQL Database password")
var schema = flag.Bool("c", false, "create database")

func main() {
	flag.Parse()

	db, err := hub.NewDataBase(1_000, *username, *password)
	if err != nil {
		log.Fatal(err)
	}

	if *schema {
		err = createSchema(db)
		if err != nil {
			log.Fatal()
		}
	}

	err = db.SelectDatabase("chroma_hub", *username, *password)
	if err != nil {
		log.Fatal()
	}

	r := db.Router()

	if *profile != "" {
		pprof.Register(r, "debug/pprof")
	}

	r.Run(fmt.Sprintf("localhost:%d", *port))
}

func createSchema(db *hub.DataBase) (err error) {
	var proceed string
	fmt.Printf("Importing schema (this will overwrite any existing chroma_hub database). Proceed (Y/n) ")
	fmt.Scan(&proceed)

	if proceed == "Y" {
		schema := "library/hub/chroma_hub.sql"

		err = db.ImportSchema(schema)
		return
	}

	err = fmt.Errorf("Schema not imported, exiting.")
	return
}
