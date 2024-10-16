package main

import (
	"chroma-viz/library/hub"
	"flag"
	"fmt"
)

var usage = `Usage:
    - import [archive|template|assets] <filename>
    - export [archive|assets] <filename>
    - export [template] <id> <filename>
    - generate <num. templates> <num. geometries>
    - clean 
`

func printMessage(port int, s string) {
	fmt.Printf("[Chroma Hub - %d]", port)

	if s == "" {
		fmt.Printf(": ")
	} else {
		fmt.Printf("  %s\n", s)
	}
}

var port = flag.Int("port", 9000, "chroma hub port")
var username = flag.String("u", "", "SQL Database username")
var password = flag.String("p", "", "SQL Database password")
var schema = flag.Bool("c", false, "create database")

func main() {
	flag.Parse()

	db, err := hub.NewDataBase(1_000, *username, *password)
	if err != nil {
		printMessage(*port, err.Error())
		return
	}

	if *schema {
		err = createSchema(db)
		if err != nil {
			printMessage(*port, err.Error())
			return
		}
	}

	err = db.SelectDatabase("chroma_hub", *username, *password)
	if err != nil {
		printMessage(*port, err.Error())
		return
	}

	db.StartRestAPI(*port)
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
