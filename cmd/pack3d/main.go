/* Instruction: write down this into the command line to use the software
<pack3d {your_frame_size_X,your_frame_size_Y,your_frame_size_Z} spacing your_output_name STL_numbers STL_path STL_numbers STL_path .....>
For example: <pack3d {100,100,100} 5 Pikacu 1 /home/corner.stl 2 /home/Pika.stl>*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	type TransMap struct {
		Filename         string
		Transformation   [4][4]float64
	}

	type err_msg struct {
		Error            string
	}

	var (
		singleStlSize []fauxgl.Vector
	    done          func()
	    totalVolume   float64
		dimension     []float64
		ntime         int
		srcStlNames   []string
		transMaps     []TransMap
		)

	rand.Seed(time.Now().UTC().UnixNano())

	model := pack3d.NewModel()
	count := 1
	ok := false

	/*Loading frame size*/
	for _, j := range os.Args[1:5]{
		_dimension, err := strconv.ParseInt(j, 0, 0)
		if err == nil{
			dimension = append(dimension, float64(_dimension))
			continue
		}
	}
	frameSize := fauxgl.V(dimension[0]/2, dimension[1]/2, dimension[2]/2)
	spacing := dimension[3]/2.0

	/* Loading stl models */
	for _, arg := range os.Args[6:] {
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
			srcStlNames = append(srcStlNames, arg)
		}

		fmt.Printf("  %d triangles\n", len(mesh.Triangles))
		fmt.Printf("  %g x %g x %g\n", size.X, size.Y, size.Z)

		done = timed("centering mesh")
		mesh.Center()
		done()

		done = timed("building bvh tree")

		model.Add(mesh, bvhDetail, count, spacing)
		ok = true
		done()
	}

	if !ok {
		fmt.Println("Usage: pack3d frame_size output_fname N1 mesh1.stl N2 mesh2.stl ...")
		fmt.Println(" - Packs N copies of each mesh into as small of a volume as possible.")
		fmt.Println(" - Runs forever, looking for the best packing.")
		fmt.Println(" - Results are written to disk whenever a new best is found.")
		return
	}

	side := math.Pow(totalVolume, 1.0/3)
	model.Deviation = side / 32  //it is not the distance between objects. And it seems that it will not reflect the distance.

	best := 1e9  //the best score
	/* This loop is to find the best packing stl, thus it will generate mutiple output
Add 'break' in the loop to stop program */
	start := time.Now()
	for {
		model, ntime = model.Pack(annealingIterations, nil, singleStlSize, frameSize)
		/* ntime is the times of trial to find a output solution, if after trying for 19990 times
		and no solution is found, then reset the model and try again. Usually if there is a solution,
		ntime will be 1 or 2 for most cases. */
		if ntime >= 19990{
			/* There is a case that even I reset the model for many times, I still can't find a solution,
			In this case, I need to set a threshold (20 second) to stop the software*/
			if time.Since(start).Seconds() <= 20{
				model.Reset()
				continue
			}else{
				err_content := err_msg{"Cannot get a result, please decrease your numbers of STLs or enlarge the frame sizes"}
				fmt.Println(err_content.Error)
				//err_json, err := json.Marshal(err_content)
				_, err := json.Marshal(err_content)
				if err != nil{
					fmt.Println("error:", err)
				}

				/*Do the following lines if want to return a json file including the error content*/
				//ioutil.WriteFile(fmt.Sprintf("%s.json", os.Args[5]), err_json, 0644)
				//break

				os.Exit(508)
			}

		}
		score := model.Energy()  // score < 1, the smaller the better
		if score < best {
			best = score
			done = timed("writing mesh")
			transformation := model.Transformation()
			for j:=0; j<len(transformation); j++{
				t := transformation[j]
				transMatrix := [4][4]float64{{t.X00, t.X01, t.X02, t.X03},{t.X10, t.X11, t.X12, t.X13},{t.X20, t.X21, t.X22, t.X23},{t.X30, t.X31, t.X32, t.X33}}
				content := TransMap{srcStlNames[j], transMatrix}
				transMaps = append(transMaps, content)
			}
			positions_json, err := json.Marshal(transMaps)
			if err != nil {
				fmt.Println("error:", err)
			}
			ioutil.WriteFile(fmt.Sprintf("%s.json", os.Args[5]), positions_json, 0644)
			//os.Stdout.Write(positions_json)

			/*Add the following line if want to generate the packing STL*/
			//model.Mesh().SaveSTL(fmt.Sprintf("pack3d-%s.stl", os.Args[5]))

			// model.TreeMesh().SaveSTL(fmt.Sprintf("out%dtree.stl", int(score*100000)))
			done()
			os.Exit(201)
		}
		model.Reset()
	}
}