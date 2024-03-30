package templates

import "log"

type Temps struct {
	Temps map[int]*Template
}

func NewTemps() *Temps {
	temps := &Temps{}
	temps.Temps = make(map[int]*Template)

	return temps
}

func (temp *Temps) SetTemplate(id, layer, num_geo int, title string) {
	if _, ok := temp.Temps[id]; ok {
		log.Printf("Template %d already exists", id)
	}

	temp.Temps[id] = NewTemplate(title, id, layer, num_geo)
}
