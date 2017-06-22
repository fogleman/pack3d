package main

import (
	"fmt"
	"os"
	"path"

	. "github.com/fogleman/fauxgl"
)

func separate(filename string) {
	mesh, err := LoadMesh(filename)
	if err != nil {
		panic(err)
	}

	lookup := make(map[Vector][]*Triangle)
	for _, t := range mesh.Triangles {
		lookup[t.V1.Position] = append(lookup[t.V1.Position], t)
		lookup[t.V2.Position] = append(lookup[t.V2.Position], t)
		lookup[t.V3.Position] = append(lookup[t.V3.Position], t)
	}

	var groups [][]*Triangle
	seen := make(map[*Triangle]bool)
	done := false
	for !done {
		done = true
		var q []*Triangle
		for _, t := range mesh.Triangles {
			if !seen[t] {
				q = append(q, t)
				done = false
				break
			}
		}
		var group []*Triangle
		for len(q) > 0 {
			var t *Triangle
			t, q = q[len(q)-1], q[:len(q)-1]
			if seen[t] {
				continue
			}
			group = append(group, t)
			seen[t] = true
			for _, v := range []Vertex{t.V1, t.V2, t.V3} {
				for _, u := range lookup[v.Position] {
					if !seen[u] {
						q = append(q, u)
					}
				}

			}
		}
		if len(group) > 0 {
			groups = append(groups, group)
		}
	}
	for i, group := range groups {
		fmt.Println(len(group))
		mesh := NewTriangleMesh(group)
		ext := path.Ext(filename)
		mesh.SaveSTL(fmt.Sprintf(filename[:len(filename)-len(ext)]+".%d.stl", i))
	}
}

func main() {
	for _, filename := range os.Args[1:] {
		separate(filename)
	}
}
