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

var hub = NewDataBase()
var usage = "Usage: {import, export} {archive, template} <filename>.json"
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
	ok := true
	hubPort = port

	StartHub(port, -1)

	read := bufio.NewScanner(os.Stdin)
	for ok {
		printMessage("")
		read.Scan()
		input := strings.Split(read.Text(), " ")

		if len(input) != 3 {
			fmt.Println(usage)
			continue
		}

		switch input[0] {
		case "import":
			Import(input[1], input[2])
		case "export":
			Export(input[1], input[2])
		default:
			fmt.Println(usage)
		}
	}
}

func StartHub(port, count int) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Error creating server (%s)", err)
	}

	go hub.SendHub(ln)
}

/*
   Send graphics hub to client using the following grammar

   S -> {'num_temp': num, 'templates': [T]}
   T -> {'id': num, 'num_geo': num, 'geometry': [G]} | T, T
   G -> {'id': num, 'type': string, 'attr': [A]} | G, G
   A -> {'name': string, 'value': string} | A, A

*/

func Import(typed, file string) {
	var err error
	switch typed {
	case "archive":
		err = ImportArchive(file)
	case "template":
		err = ImportTemplate(file)
	default:
		fmt.Println(usage)
	}

	if err != nil {
		log.Print(err)
	}
}

func Export(typed, file string) {
	switch typed {
	case "archive":
		ExportArchive(file)
	case "template":
		ExportTemplate(file)
	default:
		fmt.Println(usage)
	}
}

func ImportArchive(fileName string) error {
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

		hub.Templates[temp.TempID] = temp
		s := fmt.Sprintf("Loaded Template %d (%s)", temp.TempID, temp.Title)
		printMessage(s)
	}

	//log.Printf("Imported Hub")
	return nil
}

func ExportArchive(fileName string) {
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

func ImportTemplate(fileName string) error {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var temp templates.Template
	err = json.Unmarshal(buf, &temp)
	if err != nil {
		return err
	}

	if _, ok := hub.Templates[temp.TempID]; ok {
		return fmt.Errorf("Template ID %d already exists", temp.TempID)
	}

	hub.Templates[temp.TempID] = &temp
	s := fmt.Sprintf("Loaded Template %d (%s)", temp.TempID, temp.Title)
	printMessage(s)

	return nil
}

func ExportTemplate(fileName string) error {
	return fmt.Errorf("Not implemented")
}
