package hub

import (
	"bufio"
	"chroma-viz/library/templates"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var usage = `Usage:
    - import [archive|template] <filename>.json
    - import asset <filename>.png <directory> <name> <image id>
    - export [archive|template] <filename>.json
    - clean 
`

var hubPort int

func Logger(message string, args ...any) {
	file, err := os.Open("./chroma_hub.log")
	if err != nil {
		log.Printf("Error opening log file (%s)", err)
	}

	s := log.Prefix() + fmt.Sprintf(message, args...)
	file.Write([]byte(s))
}

func printMessage(s string) {
	fmt.Printf("[Chroma Hub - %d]", hubPort)

	if s == "" {
		fmt.Printf(": ")
	} else {
		fmt.Printf("  %s\n", s)
	}
}

func HubApp(port int) {
	hub := NewDataBase(1_000)
	ok := true
	hubPort = port

	go hub.StartHub(port)

	read := bufio.NewScanner(os.Stdin)
	for ok {
		printMessage("")
		read.Scan()
		input := strings.Split(read.Text(), " ")

		switch input[0] {
		case "import":
			imported(hub, input[1:])
		case "export":
			exported(hub, input[1:])
		case "clean":
			hub.CleanDB()
		default:
			fmt.Println(usage)
		}
	}
}

/*
   Send graphics hub to client using the following grammar

   S -> {'num_temp': num, 'templates': [T]}
   T -> {'id': num, 'num_geo': num, 'geometry': [G]} | T, T
   G -> {'id': num, 'type': string, 'attr': [A]} | G, G
   A -> {'name': string, 'value': string} | A, A

*/

func imported(hub *DataBase, inputs []string) {
	var err error
	switch inputs[0] {
	case "archive":
		err = hub.ImportArchive(inputs[1])
	case "template":
		buf, err := os.ReadFile(inputs[1])
		if err != nil {
			break
		}

		var temp templates.Template
		err = json.Unmarshal(buf, &temp)
		if err != nil {
			break
		}

		err = hub.ImportTemplate(temp)
	case "asset":
		err = hub.ImportAsset(inputs[1:])
	default:
		fmt.Println(usage)
	}

	if err != nil {
		Logger("CLI", "%s", err)
	}
}

func exported(hub *DataBase, inputs []string) {
	switch inputs[0] {
	case "archive":
		hub.ExportArchive(inputs[1])
	case "template":
		hub.ExportTemplate(inputs[1])
	default:
		fmt.Println(usage)
	}
}

func (hub *DataBase) ImportArchive(fileName string) error {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var archive DataBase
	err = json.Unmarshal(buf, &archive)
	if err != nil {
		return err
	}

	for _, temp := range archive.Templates {
		if _, ok := hub.Templates[temp.TempID]; ok {
			return fmt.Errorf("Template ID %d already exists", temp.TempID)
		}

		hub.ImportTemplate(*temp)
		s := fmt.Sprintf("Loaded Template %d (%s)", temp.TempID, temp.Title)
		printMessage(s)
	}

	for id := range archive.Assets {
		hub.Assets[id] = archive.Assets[id]
		hub.Dirs[id] = archive.Dirs[id]
		hub.Names[id] = archive.Names[id]

		s := fmt.Sprintf("Loaded Asset %d at %s/%s", id, archive.Dirs[id], archive.Names[id])
		printMessage(s)
	}

	//log.Printf("Imported Hub")
	return nil
}

func (hub *DataBase) ExportArchive(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		Logger("Couldn't open file (%s)", err)
	}
	defer file.Close()

	buf, err := json.Marshal(hub)
	if err != nil {
		Logger("Error encoding hub (%s)", err)
	}

	_, err = file.Write(buf)
	if err != nil {
		Logger("Error writing hub (%s)", err)
	}

	s := fmt.Sprintf("Exported hub to %s", fileName)
	printMessage(s)
}

func (hub *DataBase) ExportTemplate(fileName string) error {
	return fmt.Errorf("Not implemented")
}

func (hub *DataBase) ImportAsset(args []string) (err error) {
	if len(args) != 4 {
		return fmt.Errorf("Incorrect number of args to import asset")
	}

	image, err := os.ReadFile(args[0])
	if err != nil {
		return
	}

	imageID, err := strconv.Atoi(args[3])
	if err != nil {
		return
	}

	hub.Assets[imageID] = image
	hub.Dirs[imageID] = args[1]
	hub.Names[imageID] = args[2]
	return
}
