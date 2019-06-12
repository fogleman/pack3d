package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/fogleman/fauxgl"
	"pack3d/pack3d"
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
	var (
		singleStlSize []fauxgl.Vector
	    done          func()
	    totalVolume   float64
		dimension     []float64
		ntime         int
		)

	rand.Seed(time.Now().UTC().UnixNano())

	model := pack3d.NewModel()
	count := 1
	ok := false


	/*Loading frame size*/
	for _, j := range os.Args[1:4]{
		_dimension, err := strconv.ParseInt(j, 0, 0)
		if err == nil{
			dimension = append(dimension, float64(_dimension))
			continue
		}
	}
	frameSize := fauxgl.V(dimension[0]/2, dimension[1]/2, dimension[2]/2)

	/* Loading stl models */
	for _, arg := range os.Args[4:] {
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
			singleStlSize = append(singleStlSize, size)
		}

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
	model.Deviation = side / 32  //change deviation to change distance between models, set a minimum here

	best := 1e9  //the best score
	/* This loop is to find the best packing stl, thus it will generate mutiple output
Add 'break' in the loop to stop program */
	for {
		model, ntime = model.Pack(annealingIterations, nil, singleStlSize, frameSize)
		if ntime >= 19990{
			model.Reset()
			continue
		}
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