package gui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const NUMCOL = 3

const (
    PAGENUM = iota
    TITLE
    TEMPLATEID 
)

var KEYTITLE = map[int]string{
    PAGENUM: "Page Num",
    TITLE: "Title",
    TEMPLATEID: "Template ID",
}

type Page struct {
    Box         *gtk.ListBoxRow
    pageNum     int
    title       string
    templateID  int
    props       map[string]Property
    propList    []string
}

func NewPage(pageNum int, title string, temp *Template) *Page {
    page := &Page{pageNum: pageNum, title: title, templateID: temp.templateID}
    page.props = make(map[string]Property)
    page.propList = []string{"x Pos", "y Pos", "Width", "Height", "Title", "Subtitle"}

    for key, prop := range temp.props {
        page.props[key] = prop.Copy()
    }

    return page
}

func (page *Page) pageToListRow() *gtk.ListBoxRow {
    row1, _ := gtk.ListBoxRowNew()
    row1.Add(textToBuffer(strconv.Itoa(page.pageNum)))
    row1.Add(textToBuffer(page.title))

    return row1
}

type ShowTree struct {
    *gtk.TreeView
    treeList  *gtk.ListStore
    pages     map[int]*Page
    numPages  int
    edit      *Editor
    columns   [NUMCOL]bool
}

func NewShow(edit *Editor, prev *Connection) *ShowTree {
    show := &ShowTree{}
    show.TreeView, _ = gtk.TreeViewNew()
    show.pages = make(map[int]*Page)
    show.columns = [NUMCOL]bool{true, true, true}

    /* Columns */
    var column *gtk.TreeViewColumn
    for key := range show.columns {
        if show.columns[key] == false {
            continue
        }

        switch key {
        case TITLE:
            title, _ := gtk.CellRendererTextNew()
            title.SetProperty("editable", true)
            title.Connect("edited", 
                func(cell *gtk.CellRendererText, path string, text string) {
                    iter, _ := show.treeList.GetIterFromString(path)
                    id, _ := show.treeList.GetValue(iter, PAGENUM)
                    val, _ := id.GoValue()
                    pageNum, _ := strconv.Atoi(val.(string))
                    show.pages[pageNum].title = text
                    show.treeList.SetValue(iter, TITLE, text)
            })

            column, _ = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[key], title, "text", 1)
        default:
            cell, _ := gtk.CellRendererTextNew()
            column, _ = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[key], cell, "text", key)
        }

        show.AppendColumn(column)
    }

    show.treeList, _ = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)

    show.SetModel(show.treeList)
    show.edit = edit

    // TODO: remove reference to show from outside scope
    show.Connect("row-activated", 
        func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) { 
            iter, _ := show.treeList.GetIter(path)
            id, _ := show.treeList.GetValue(iter, PAGENUM)
            val, _ := id.GoValue()
            pageNum, _ := strconv.Atoi(val.(string))

            show.edit.SetPage(show.pages[pageNum])
            prev.SendPage(show.pages[pageNum], ANIMATE_ON)
        })

    return show 
}

func (show *ShowTree) NewShowPage(temp *Template) *Page {
    show.numPages++
    show.pages[show.numPages] = NewPage(show.numPages, temp.title, temp)
    page := show.pages[show.numPages]
    show.treeList.Set(
        show.treeList.Append(), 
        []int{PAGENUM, TITLE, TEMPLATEID}, 
        []interface{}{page.pageNum, page.title, page.templateID})
    return page
}

func (show *ShowTree) ImportShow(temp *TempTree, filename string) {
    pageReg, err := regexp.Compile("temp [0-9]*; title .*;")
    if err != nil {
        log.Print(err)
        return
    }

    file, err := os.Open(filename)
    if err != nil {
        log.Print(err)
        return
    }

    scanner := bufio.NewScanner(file)

    var page *Page
    for scanner.Scan() {
        line := scanner.Text()
        if pageReg.Match(scanner.Bytes()) {
            split := strings.Split(line, ";")

            tempID := parse_int_value(split[0], "temp")
            page = show.NewShowPage(temp.temps[tempID])
            page.title = strings.TrimLeft(split[1], " title ")
        } else if page != nil {
            Decode(page, line)
        }

    }

}

func (show *ShowTree) ExportShow(filename string) {
    file, err := os.Create(filename)
    if err != nil {
        log.Print(err)
    }
    defer file.Close()

    for _, page := range show.pages {
        pageString := fmt.Sprintf("temp %d; title %s;\n", page.templateID, page.title)
        file.Write([]byte(pageString))

        for name, prop := range page.props {
            file.Write(prop.Encode(name))
        }
    }
}
