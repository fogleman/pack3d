package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/fogleman/fauxgl" //fauxgl is a go library
	"github.com/fogleman/pack3d/pack3d"
)

const (
	bvhDetail           = 8
	annealingIterations = 2000000 // # of trials
)


/* This function returns current time (it's a timer) */ 
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
	var stlsize []fauxgl.Vector
	//var boundsize fauxgl.Vector
	var done func()

	boundsize := fauxgl.V(50.0, 50.0, 50.0)

	rand.Seed(time.Now().UTC().UnixNano())

	model := pack3d.NewModel()
	count := 1
	ok := false
	var totalVolume float64

	/* Loading stl models */
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
		for i:=0; i<count; i++{
			stlsize = append(stlsize, size)
		}
		//fmt.Println(reflect.TypeOf(stlsize))

		// fmt.Println(" My name is Minglun ")
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

	//fmt.Println(len(stlsize))
	side := math.Pow(totalVolume, 1.0/3)
	model.Deviation = side / 32  //change deviation to change distance between models, set a minimum here

	best := 1e9  //the best score
	/* This loop is to find the best packing stl, thus it will generate mutiple output 
Add 'break' in the loop to stop program */
	for {
		model = model.Pack(annealingIterations, nil, stlsize, boundsize)
		score := model.Energy()  // score < 1, the smaller the better
		if score < best {
			best = score
			done = timed("writing mesh")
			model.Mesh().SaveSTL(fmt.Sprintf("pack3d-%.3f.stl", score))  // calling the mesh function in model
			// model.TreeMesh().SaveSTL(fmt.Sprintf("out%dtree.stl", int(score*100000)))
			done()
		}
		model.Reset()
	}
}
