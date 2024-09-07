package main

import (
	"bufio"
	"chroma-viz/library/hub"
	"chroma-viz/library/templates"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var usage = `Usage:
    - import [archive|template] <filename>.json
    - import asset <filename>.png <directory> <name> <image id>
    - export archive <filename>.json
    - export template <template id> <filename>.json
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

	ok := true
	go db.StartHub(*port)

	read := bufio.NewScanner(os.Stdin)
	for ok {
		printMessage(*port, "")
		read.Scan()
		input := strings.Split(read.Text(), " ")

		switch input[0] {
		case "import":
			imported(db, input[1:])
		case "export":
			exported(db, input[1:])
		case "generate":
			generate(db, input[1:])
		case "clean":
			db.CleanDB()
		default:
			fmt.Println(usage)
		}
	}
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

func imported(db *hub.DataBase, inputs []string) {
	var err error

	if len(inputs) == 0 {
		fmt.Println(usage)
		return
	}

	switch inputs[0] {
	case "archive":
		if len(inputs) != 2 {
			fmt.Println(usage)
			return
		}

		err = db.ImportArchive(inputs[1])
	case "template":
		if len(inputs) != 2 {
			fmt.Println(usage)
			return
		}

		var temp templates.Template
		temp, err = templates.NewTemplateFromFile(inputs[1])
		if err != nil {
			break
		}

		err = db.ImportTemplate(temp)
	case "asset":
		if len(inputs) != 5 {
			fmt.Println(usage)
			return
		}

		err = db.ImportAsset(inputs[1:])
	default:
		fmt.Println(usage)
	}

	if err != nil {
		hub.Logger("CLI: %s", err)
	}
}

func exported(db *hub.DataBase, inputs []string) {
	switch inputs[0] {
	case "archive":
		if len(inputs) != 2 {
			fmt.Println(usage)
			return
		}

		db.ExportArchive(inputs[1])
	case "template":
		if len(inputs) != 3 {
			fmt.Println(usage)
			return
		}

		tempID, err := strconv.ParseInt(inputs[1], 10, 64)
		if err != nil {
			hub.Logger("CLI: %s", err)
			return
		}

		err = db.ExportTemplate(inputs[2], tempID)
		if err != nil {
			hub.Logger("CLI: %s", err)
			return
		}

	default:
		fmt.Println(usage)
	}
}

func generate(db *hub.DataBase, inputs []string) {
	if len(inputs) != 2 {
		fmt.Println(usage)
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
