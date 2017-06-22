package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	. "github.com/fogleman/fauxgl"
)

func main() {
	path := os.Args[1]
	_count, _ := strconv.ParseInt(os.Args[2], 0, 0)
	count := int(_count)

	mesh, err := LoadMesh(path)
	if err != nil {
		panic(err)
	}

	mesh.MoveTo(Vector{}, Vector{})

	meshes := make([]*Mesh, 0, count)
	n := len(mesh.Triangles) / count
	for i := 0; i < len(mesh.Triangles); i += n {
		m := NewTriangleMesh(mesh.Triangles[i : i+n])
		meshes = append(meshes, m)
	}

	sort.Slice(meshes, func(i, j int) bool {
		a := meshes[i].BoundingBox().Min
		b := meshes[j].BoundingBox().Min
		a = Vector{a.Z, a.X, a.Y}
		b = Vector{b.Z, b.X, b.Y}
		return a.Less(b)
	})

	result := NewEmptyMesh()
	for _, mesh := range meshes {
		result.Add(mesh)
		fmt.Println(mesh.BoundingBox())
	}
	result.SaveSTL("out.stl")
}
