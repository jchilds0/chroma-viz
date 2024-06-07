package main

import (
	"bufio"
	"chroma-viz/library/util"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	PREVIEW = iota
	TOP_LEFT
	LOWER_FRAME
	TICKER
)

type TempTree struct {
	treeView     *gtk.TreeView
	treeList     *gtk.ListStore
	sendTemplate func(int)
}

func NewTempTree(templateToShow func(int)) *TempTree {
	var err error
	temp := &TempTree{sendTemplate: templateToShow}

	temp.treeView, err = gtk.TreeViewNew()
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	// create tree columns
	cell, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute("Name", cell, "text", 0)
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	temp.treeView.AppendColumn(column)
	column, err = gtk.TreeViewColumnNewWithAttribute("Template ID", cell, "text", 1)
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	temp.treeView.AppendColumn(column)

	temp.treeList, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_INT)
	if err != nil {
		log.Fatalf("Error creating temp list (%s)", err)
	}

	temp.treeView.SetModel(temp.treeList)

	// send template to show on double click
	temp.treeView.Connect("row-activated",
		func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {
			iter, err := temp.treeList.GetIter(path)
			if err != nil {
				log.Fatalf("Error sending template to show (%s)", err)
			}

			model := &temp.treeList.TreeModel
			tempID, err := util.ModelGetValue[int](model, iter, 1)
			if err != nil {
				log.Fatalf("Error sending template to show (%s)", err)
			}

			temp.sendTemplate(tempID)
		})

	return temp
}

func (temp *TempTree) ImportTemplates(hub net.Conn) {
	if hub == nil {
		log.Print("Chroma Hub is disconnected")
		return
	}

	s := fmt.Sprintf("ver 0 1 tempids;")
	hub.Write([]byte(s))

	buf := bufio.NewReader(hub)
	for {
		s, err := buf.ReadString(';')
		if err != nil {
			return
		}

		s = strings.TrimSuffix(s, ";")
		if s == "EOF" {
			return
		}

		data := strings.Split(s, " ")

		tempID, err := strconv.Atoi(data[0])
		if err != nil {
			log.Printf("Error reading template id (%s)", err)
			continue
		}
		title := strings.Join(data[1:], " ")

		err = temp.treeList.Set(
			temp.treeList.Append(),
			[]int{0, 1},
			[]interface{}{title, tempID},
		)
		if err != nil {
			log.Printf("Error adding template to gtk treestore (%s)", err)
			return
		}
	}
}
