package artist

import (
	"fmt"
	"log"
)

/*

   Geom manages the instances of a Prop.
   Artist creates a graphics hub with a constant number
   of geometries, and the user adds or removes geometries
   by asking Geom for a geometry ID.

*/

type geom struct {
	index int
	count int
	alloc []bool
}

func newGeom(index, count int) *geom {
	geo := &geom{index: index, count: count}
	geo.alloc = make([]bool, count)

	return geo
}

func (g *geom) allocGeom() (int, error) {
	for i := range g.alloc {
		if g.alloc[i] {
			continue
		}

		g.alloc[i] = true
		return i + g.index, nil
	}

	return 0, fmt.Errorf("No more geometries")
}

func (g *geom) freeGeom(id int) {
	if id < 0 || id >= len(g.alloc) {
		log.Printf("id out of range (%d)", id)
		return
	}

	g.alloc[id] = false
}
