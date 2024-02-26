package viz

import (
	"chroma-viz/shows"
	"chroma-viz/tcp"
	"chroma-viz/templates"
	"fmt"
	"log"
	"strconv"

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
type ShowTree struct {
    treeView  *gtk.TreeView
    treeList  *gtk.ListStore
    // pages     map[int]*shows.Page
    show      *shows.Show
    columns   [NUMCOL]bool
}

func NewShowTree(pageToEditor func(*shows.Page)) *ShowTree {
    var err error
    showTree := &ShowTree{}

    showTree.treeView, err = gtk.TreeViewNew()
    if err != nil {
        log.Fatalf("Error creating show (%s)", err)
    }

    showTree.show = shows.NewShow()
    showTree.treeView.SetReorderable(true)
    showTree.columns = [NUMCOL]bool{true, true, true}

    // create tree columns
    var column *gtk.TreeViewColumn
    for key := range showTree.columns {
        if showTree.columns[key] == false {
            continue
        }

        switch key {
        case TITLE:
            title, err := gtk.CellRendererTextNew()
            if err != nil {
                log.Fatalf("Error creating show (%s)", err)
            }

            title.SetProperty("editable", true)
            title.Connect("edited", 
                func(cell *gtk.CellRendererText, path string, text string) {
                    iter, err := showTree.treeList.GetIterFromString(path)
                    if err != nil {
                        log.Printf("Error editing page (%s)", err)
                        return
                    }

                    id, err := showTree.treeList.GetValue(iter, PAGENUM)
                    if err != nil {
                        log.Printf("Error editing page (%s)", err)
                        return
                    }

                    val, err := id.GoValue()
                    if err != nil {
                        log.Printf("Error editing page (%s)", err)
                        return
                    }

                    pageNum, err := strconv.Atoi(val.(string))
                    if err != nil {
                        log.Printf("Error editing page (%s)", err)
                        return
                    }

                    if _, ok := showTree.show.Pages[pageNum]; !ok {
                        log.Print("Error getting page")
                        return
                    }

                    showTree.show.Pages[pageNum].Title = text
                    showTree.treeList.SetValue(iter, TITLE, text)
            })

            column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[key], title, "text", 1)
            if err != nil {
                log.Fatalf("Error creating show (%s)", err)
            }

        default:
            cell, err := gtk.CellRendererTextNew()
            if err != nil {
                log.Fatalf("Error creating show (%s)", err)
            }

            column, err = gtk.TreeViewColumnNewWithAttribute(KEYTITLE[key], cell, "text", key)
            if err != nil {
                log.Fatalf("Error creating show (%s)", err)
            }

        }

        showTree.treeView.AppendColumn(column)
    }

    showTree.treeList, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
    if err != nil {
        log.Fatalf("Error creating show (%s)", err)
    }

    showTree.treeView.SetModel(showTree.treeList)

    // send page to editor on double click
    showTree.treeView.Connect("row-activated", 
        func(tree *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) { 
            iter, err := showTree.treeList.GetIter(path)
            if err != nil {
                log.Printf("Error sending page to editor (%s)", err)
                return
            }

            id, err := showTree.treeList.GetValue(iter, PAGENUM)
            if err != nil {
                log.Printf("Error sending page to editor (%s)", err)
                return
            }

            val, err := id.GoValue()
            if err != nil {
                log.Printf("Error sending page to editor (%s)", err)
                return
            }

            pageNum, err := strconv.Atoi(val.(string))
            if err != nil {
                log.Printf("Error sending page to editor (%s)", err)
                return
            }

            pageToEditor(showTree.show.Pages[pageNum])
            SendPreview(showTree.show.Pages[pageNum], tcp.ANIMATE_ON)
        })

    return showTree 
}

func (showTree *ShowTree) NewShowPage(page *shows.Page) {
    if page == nil {
        log.Print("Missing template")
        return
    }

    showTree.treeList.Set(
        showTree.treeList.Append(), 
        []int{PAGENUM, TITLE, TEMPLATEID}, 
        []interface{}{page.PageNum, page.Title, page.TemplateID},
    )
}

func (showTree *ShowTree) ImportShow(temps *TempTree, filename string) {
    showTree.show.ImportShow(temps.Temps, filename)

    for _, page := range showTree.show.Pages {
        showTree.NewShowPage(page)
    }
}

func (showTree *ShowTree) ImportPage(title string, temp *templates.Template) error {
    if temp == nil {
        return fmt.Errorf("Missing template")
    }

    page := showTree.show.AddPage(title, temp)
    showTree.NewShowPage(page)

    return nil 
}

