package hub

import (
	"bufio"
	"chroma-viz/library/templates"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var usage = "Usage: \n" +
	"\t- import [archive|template] <filename>.json\n" +
	"\t- import asset <filename>.png <directory> <name> <image id>\n" +
	"\t- export [archive|template] <filename>.json"

var hubPort int

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

	StartHub(hub, port)

	read := bufio.NewScanner(os.Stdin)
	for ok {
		printMessage("")
		read.Scan()
		input := strings.Split(read.Text(), " ")

		if len(input) < 3 {
			fmt.Println(usage)
			continue
		}

		switch input[0] {
		case "import":
			imported(hub, input[1:])
		case "export":
			exported(hub, input[1:])
		default:
			fmt.Println(usage)
		}
	}
}

func StartHub(hub *DataBase, port int) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Error creating server (%s)", err)
	}

	go hub.AcceptHubConn(ln)
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
		err = hub.ImportTemplate(inputs[1])
	case "asset":
		err = hub.ImportAsset(inputs[1:])
	default:
		fmt.Println(usage)
	}

	if err != nil {
		log.Print(err)
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

	// for _, temp := range archive.Templates {
	// 	if _, ok := hub.Templates[temp.TempID]; ok {
	// 		return fmt.Errorf("Template ID %d already exists", temp.TempID)
	// 	}
	//
	// 	hub.Templates[temp.TempID] = temp
	// 	s := fmt.Sprintf("Loaded Template %d (%s)", temp.TempID, temp.Title)
	// 	printMessage(s)
	// }

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
		log.Fatalf("Couldn't open file (%s)", err)
	}
	defer file.Close()

	buf, err := json.Marshal(hub)
	if err != nil {
		log.Printf("Error encoding hub (%s)", err)
	}

	_, err = file.Write(buf)
	if err != nil {
		log.Printf("Error writing hub (%s)", err)
	}

	s := fmt.Sprintf("Exported hub to %s", fileName)
	printMessage(s)
}

func (hub *DataBase) ImportTemplate(fileName string) error {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var temp templates.Template
	err = json.Unmarshal(buf, &temp)
	if err != nil {
		return err
	}

	// if _, ok := hub.Templates[temp.TempID]; ok {
	// 	printMessage(fmt.Sprintf("Template %d already exists, overwriting", temp.TempID))
	// }

	// hub.Templates[temp.TempID] = &temp
	// s := fmt.Sprintf("Loaded Template %d (%s)", temp.TempID, temp.Title)
	// printMessage(s)

	return nil
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
