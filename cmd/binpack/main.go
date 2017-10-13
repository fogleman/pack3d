package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/fogleman/fauxgl"
	"github.com/fogleman/pack3d/binpack"
)

const (
	SizeX = 165
	SizeY = 165
	SizeZ = 320
)

var Rotations []fauxgl.Matrix

func init() {
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			m := fauxgl.Rotate(fauxgl.Vector{0, 0, 1}, float64(i)*math.Pi/2)
			switch j {
			case 1:
				m = m.Rotate(fauxgl.Vector{1, 0, 0}, math.Pi/2)
			case 2:
				m = m.Rotate(fauxgl.Vector{0, 1, 0}, math.Pi/2)
			}
			Rotations = append(Rotations, m)
		}
	}
}

func timed(name string) func() {
	if len(name) > 0 {
		fmt.Printf("%s... ", name)
	}
	start := time.Now()
	return func() {
		fmt.Println(time.Since(start))
	}
}

func main() {
	const S = 100
	const P = 2.5

	var items []binpack.Item
	var meshes []*fauxgl.Mesh

	var done func()

	score := 1
	ok := false
	for _, arg := range os.Args[1:] {
		_score, err := strconv.ParseInt(arg, 0, 0)
		if err == nil {
			score = int(_score)
			continue
		}

		done = timed("loading mesh")
		mesh, err := fauxgl.LoadMesh(arg)
		if err != nil {
			panic(err)
		}
		done()

		i := len(meshes)
		meshes = append(meshes, mesh)
		box := mesh.BoundingBox()
		for j, m := range Rotations {
			id := i*len(Rotations) + j
			s := box.Transform(m).Size()
			sx := int(math.Ceil((s.X + P*2) * S))
			sy := int(math.Ceil((s.Y + P*2) * S))
			sz := int(math.Ceil((s.Z + P*2) * S))
			items = append(items, binpack.Item{id, score, binpack.Vector{sx, sy, sz}})
		}
		ok = true
	}

	if !ok {
		fmt.Println("Usage: binpack N1 mesh1.stl N2 mesh2.stl ...")
		fmt.Println(" - Packs as many items into the volume as possible.")
		fmt.Println(" - N specifies how many items the mesh contains.")
		fmt.Println(" - Provide multiple pack3d meshes for best results.")
		return
	}

	done = timed("bin packing")
	box := binpack.Box{binpack.Vector{}, binpack.Vector{SizeX * S, SizeY * S, SizeZ * S}}
	result := binpack.Pack(items, box)
	done()

	fmt.Printf("packed %d items\n", result.Score)

	done = timed("building result")
	mesh := fauxgl.NewEmptyMesh()
	for _, placement := range result.Placements {
		p := placement.Position
		v := fauxgl.Vector{float64(p.X)/S + P, float64(p.Y)/S + P, float64(p.Z)/S + P}
		i := placement.Item.ID / len(Rotations)
		j := placement.Item.ID % len(Rotations)
		m := meshes[i].Copy()
		m.Transform(Rotations[j])
		m.MoveTo(v, fauxgl.Vector{})
		mesh.Add(m)
	}
	mesh.MoveTo(fauxgl.Vector{}, fauxgl.Vector{})
	done()

	done = timed("writing stl file")
	mesh.SaveSTL("binpack.stl")
	done()

	// for _, p := range result.Placements {
	// 	fmt.Println(p)
	// }
	// fmt.Println(result.Score)
}
