package main

import (
	"chroma-viz/library/hub"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
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
var gen = flag.String("g", "", "generate templates")

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

	if *gen != "" {
		input := strings.Split(*gen, ",")
		generate(db, input)
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

func generate(db *hub.DataBase, inputs []string) {
	if len(inputs) != 2 {
		return
	}

	db.CleanDB()
	numTemp, err := strconv.Atoi(inputs[0])
	if err != nil {
		fmt.Println(usage)
		return
	}

	numGeo, err := strconv.Atoi(inputs[1])
	if err != nil {
		fmt.Println(usage)
		return
	}

	start := time.Now()

	var i int64
	for i = 1; i <= int64(numTemp); i++ {
		randomTemplate(db, i, numGeo)
	}

	end := time.Now()
	elapsed := end.Sub(start)
	hub.Logger("Generated Random Hub in %s\n", elapsed)
}
