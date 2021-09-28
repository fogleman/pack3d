/*
Instruction: write down this into the command line to use the software
<pack3d {your_frame_size_X,your_frame_size_Y,your_frame_size_Z} spacing your_output_name STL_numbers STL_path STL_numbers STL_path .....>
For example: <pack3d {100,100,100} 5 Pikacu 1 /home/corner.stl 2 /home/Pika.stl>
The unit of frame and spacing is millimeters.
*/

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

	"github.com/Authentise/pack3d/pack3d"
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

	if os.Args[1] == "--version" {
		fmt.Println("Pack3d 1.5.0")
		return
	}

	type TransMap struct {
		Filename         string
		Transformation   [4][4]float64
		VolumeWithSpacing   float64
	}

	type err_msg struct {
		Error            string
	}

	var (
		singleStlSize []fauxgl.Vector
		scaleStl      []fauxgl.Matrix
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
	scale := 1.0
	scaleMatrix := fauxgl.Matrix{}
	ok := false

	// Loading build_volume size
	// os.Args[1:3] --> build dimensions.
	// os.Args[4]   --> minimum distance.
	// os.Args[5]   --> output path.
	for _, j := range os.Args[1:5]{
		_dimension, err := strconv.ParseFloat(j, 64)
		if err == nil{
			dimension = append(dimension, float64(_dimension))
			continue
		}
	}
	spacing := dimension[3]/2.0
	// frameSize is the vertex in the first quadrant
	frameSize := fauxgl.V(dimension[0]/2.0, dimension[1]/2.0, dimension[2]/2.0)
	buildVolume := dimension[0] * dimension[1] * dimension[2]
	//fmt.Println(frameSize)

	/* Loading stl models */
	// os.Args[6], os.Args[6+3], ... --> objects counts.
	// os.Args[7], os.Args[7+3], ... --> objects scales.
	// os.Args[8], os.Args[8+3], ... --> objects filenames/path.
	for _, arg := range os.Args[6:] {

		// object's count.
		_count, err := strconv.ParseInt(arg, 0, 0)
		if err == nil {
			count = int(_count)
			continue
		}

		// object's scale.
		_scale, err := strconv.ParseFloat(arg, 64)
		if err == nil {
			scale = float64(_scale)
			continue
		}

		// object's filename.
		done = timed(fmt.Sprintf("loading mesh %s", arg))
		mesh, err := fauxgl.LoadMesh(arg)
		if err != nil {
			panic(err)
		}
		done()

		// Notice that the scaling is applied before the computation of the mesh's BoundingBox and volume.
		done = timed("scaling mesh")
		scaleMatrix = fauxgl.Scale(fauxgl.V(scale, scale, scale))
		mesh.Transform(scaleMatrix)
		done()

		totalVolume += mesh.BoundingBox().Volume()
		size := mesh.BoundingBox().Size()
		for i:=0; i<count; i++{
			singleStlSize = append(singleStlSize, size)
			srcStlNames = append(srcStlNames, arg)
			scaleStl = append(scaleStl, scaleMatrix)
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

		fmt.Println("______________________________________________________")
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

	/* This loop is to find the best packing stl, thus it will generate mutiple output
		Add 'break' in the loop to stop program */
	start := time.Now()
	maxItemNum := len(model.Items)
	var timeLimit float64
	fillVolumeWithSpacing := 0.0
	totalFillVolume := 0.0
	null := fauxgl.Matrix{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ,0}
	timeLimit = 10

	minItemNum := 0
	packItemNum := maxItemNum
	success_model := pack3d.NewModel()

	for {
		model, ntime = model.Pack(annealingIterations, nil, singleStlSize, frameSize, packItemNum)
		/* ntime is the times of trial to find a output solution, if after trying for 100 times
		and no solution is found, then reset the model and try again. Usually if there is a solution,
		ntime will be 1 or 2 for most cases. */
		if ntime >= 100{
			/* There is a case that even I reset the model for many times, I still can't find a solution,
			In this case, I need to set a threshold (20 second) to stop the software*/
			if time.Since(start).Seconds() <= timeLimit{
				model.Reset()
				continue
			}else{
				// Linear search
				//packItemNum -= 1

				// Binary search
				fmt.Println("Failed")
				fmt.Println("packing#, max#, min# is: ", packItemNum, maxItemNum, minItemNum)
				fmt.Println("-----------------------------------")
				maxItemNum = packItemNum - 1
				packItemNum = int(math.Ceil(float64((maxItemNum + minItemNum) / 2)))

				model.Reset()
				model.Transformation()[packItemNum] = null
				start = time.Now()

				if minItemNum > maxItemNum{
					break
				}

				continue

				//TODO: Unblock the following lines if want to return a json file including the error content
				/*
				err_content := err_msg{"Cannot get a result, please decrease your numbers of STLs or enlarge the frame sizes"}
				fmt.Println(err_content.Error)
				err_json, err := json.Marshal(err_content)
				_, err := json.Marshal(err_content)
				if err != nil{
				fmt.Println("error:", err)
				}
				ioutil.WriteFile(fmt.Sprintf("%s.json", os.Args[5]), err_json, 0644)
				break
				*/
			}
		}

		// Binary search
		fmt.Println("Succeeded")
		fmt.Println("packing#, max#, min# is: ", packItemNum, maxItemNum, minItemNum)
		fmt.Println("-----------------------------------------")
		minItemNum = packItemNum + 1
		packItemNum = int(math.Ceil(float64((maxItemNum + minItemNum) / 2)))
		success_model = model
		start = time.Now()

		if minItemNum > maxItemNum{
			break
		}
		model.Reset()
	}

	done = timed("writing mesh")
	var (transMatrix [4][4]float64
		fillPercentage float64)
	transformation := success_model.Transformation()
	for j:=0; j<len(success_model.Items); j++{
		t := transformation[j]
		st := t.Mul(scaleStl[j])  // scaled transformation for the j-th mesh.
		fillVolumeWithSpacing = (singleStlSize[j].X + spacing) * (singleStlSize[j].Y + spacing) * (singleStlSize[j].Z + spacing)
		if j<packItemNum {
			totalFillVolume += fillVolumeWithSpacing
			transMatrix = [4][4]float64{{st.X00, st.X01, st.X02, st.X03}, {st.X10, st.X11, st.X12, st.X13}, {st.X20, st.X21, st.X22, st.X23}, {st.X30, st.X31, st.X32, st.X33}}
		}else{
			transMatrix = [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}
		}
		// It's actually the bounding box filling percentage
		fillPercentage = totalFillVolume/buildVolume
		content := TransMap{srcStlNames[j], transMatrix, fillVolumeWithSpacing}
		transMaps = append(transMaps, content)
	}
	positions_json, err := json.Marshal(transMaps)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("the fill percentage is:", fillPercentage)
	ioutil.WriteFile(fmt.Sprintf("%s.json", os.Args[5]), positions_json, 0644)
	//os.Stdout.Write(positions_json)

	//TODO: Unblock the following line if want to generate the packing STL
	/*
		model.Mesh().SaveSTL(fmt.Sprintf("%s.stl", os.Args[5]))
		model.TreeMesh().SaveSTL(fmt.Sprintf("out%dtree.stl", int(score*100000)))
	*/
	done()
}