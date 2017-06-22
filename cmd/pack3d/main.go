package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/fogleman/fauxgl"
	"github.com/fogleman/pack3d/pack3d"
)

const (
	bvhDetail           = 8
	annealingIterations = 2000000
)

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
	var done func()

	rand.Seed(time.Now().UTC().UnixNano())

	model := pack3d.NewModel()
	count := 1
	ok := false
	var totalVolume float64
	for _, arg := range os.Args[1:] {
		_count, err := strconv.ParseInt(arg, 0, 0)
		if err == nil {
			count = int(_count)
			continue
		}

		done = timed(fmt.Sprintf("loading mesh %s", arg))
		mesh, err := fauxgl.LoadMesh(arg)
		if err != nil {
			panic(err)
		}
		done()

		totalVolume += mesh.BoundingBox().Volume()
		size := mesh.BoundingBox().Size()
		fmt.Printf("  %d triangles\n", len(mesh.Triangles))
		fmt.Printf("  %g x %g x %g\n", size.X, size.Y, size.Z)

		done = timed("centering mesh")
		mesh.Center()
		done()

		done = timed("building bvh tree")
		model.Add(mesh, bvhDetail, count)
		ok = true
		done()
	}

	if !ok {
		fmt.Println("Usage: pack3d N1 mesh1.stl N2 mesh2.stl ...")
		fmt.Println(" - Packs N copies of each mesh into as small of a volume as possible.")
		fmt.Println(" - Runs forever, looking for the best packing.")
		fmt.Println(" - Results are written to disk whenever a new best is found.")
		return
	}

	side := math.Pow(totalVolume, 1.0/3)
	model.Deviation = side / 32

	best := 1e9
	for {
		model = model.Pack(annealingIterations, nil)
		score := model.Energy()
		if score < best {
			best = score
			done = timed("writing mesh")
			model.Mesh().SaveSTL(fmt.Sprintf("pack3d-%.3f.stl", score))
			// model.TreeMesh().SaveSTL(fmt.Sprintf("out%dtree.stl", int(score*100000)))
			done()
		}
		model.Reset()
	}
}
