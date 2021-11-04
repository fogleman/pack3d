/*
Instruction: write down this into the command line to use the software.

<pack3d --input_config_json_filename=json_config_file --output_packing_json_filename=export_filename>
For example: <pack3d --input_config_json_filename=input.json>

The frame and spacing's units, in the json file, are in millimeters.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/Authentise/pack3d/pack3d"
	"github.com/fogleman/fauxgl"
)

const (
	bvhDetail           = 8
	annealingIterations = 2000000 // # of trials
)

/* This function returns the current time (it is a timer). */
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
	var jsonFileArg = flag.String("input_config_json_filename", "", "json config file")
	var fileNameArg = flag.String("output_packing_json_filename", "pack3d", "export filename")
	flag.Parse()

	if os.Args[1] == "--version" {
		fmt.Println("Pack3d 1.5.0")
		return
	}

	var config Config
	if jsonFileArg != nil {
		file, err := ioutil.ReadFile(*jsonFileArg)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal([]byte(file), &config)
		if err != nil {
			log.Fatal(err)
		}
	}

	type TransMap struct {
		Filename          string
		Transformation    [4][4]float64
		VolumeWithSpacing float64
	}

	type err_msg struct {
		Error string
	}

	var (
		singleStlSize []fauxgl.Vector
		scaleStl      []fauxgl.Matrix
		done          func()
		totalVolume   float64
		ntime         int
		srcStlNames   []string
		transMaps     []TransMap
	)

	rand.Seed(time.Now().UTC().UnixNano())

	model := pack3d.NewModel()
	scale := 1.0
	scaleMatrix := fauxgl.Matrix{}
	ok := false

	spacing := config.Spacing / 2.0

	// frameSize is the vertex in the first quadrant
	frameSize := fauxgl.V(config.BuildVolume[0]/2.0, config.BuildVolume[1]/2.0, config.BuildVolume[2]/2.0)
	buildVolume := config.BuildVolume[0] * config.BuildVolume[1] * config.BuildVolume[2]
	//fmt.Println(frameSize)

	/* Loading stl models */
	coPackMap := make(map[string][]*Copack)  // object to record co-packed meshes.
	for _, item := range config.Items {

		var mesh *fauxgl.Mesh
		var err error

		if item.Copack == nil {

			done = timed(fmt.Sprintf("loading mesh %s", item.Filename))
			mesh, err = fauxgl.LoadMesh(item.Filename)
			if err != nil {
				panic(err)
			}
			done()

			// Mesh's scaling. If scaling is to be applied, it is
			// done before the computation of the BoundingBox and volume.
			scale = item.Scale
			scaleMatrix = fauxgl.Scale(fauxgl.V(scale, scale, scale))
			if scale != 1.0 {
				done = timed("scaling mesh")
				mesh.Transform(scaleMatrix)
				done()
			}

			// update arrays.
			size := mesh.BoundingBox().Size()
			for i := 0; i < item.Count; i++ {
				singleStlSize = append(singleStlSize, size)
				srcStlNames = append(srcStlNames, item.Filename)
				scaleStl = append(scaleStl, scaleMatrix)
			}

			fmt.Printf("  %d triangles\n", len(mesh.Triangles))
			fmt.Printf("  %g x %g x %g\n", size.X, size.Y, size.Z)

			done = timed("centering mesh")
			mesh.Center()
			done()

			totalVolume += mesh.BoundingBox().Volume()

		} else {

			coPackMap[item.Filename] = item.Copack

			// load and scale the main co-packing mesh.
			done = timed(fmt.Sprintf("loading main co-packing mesh %s", item.Filename))
			mesh, err = fauxgl.LoadMesh(item.Filename)
			if err != nil {
				panic(err)
			}
			done()

			// main co-packing mesh's scaling. If scaling is to be applied, it is
			// done before the computation of the BoundingBox and volume.
			scale = item.Scale
			scaleMatrix = fauxgl.Scale(fauxgl.V(scale, scale, scale))
			if scale != 1.0 {
				done = timed("scaling main co-packing mesh")
				mesh.Transform(scaleMatrix)
				done()
			}

			// load and scale the co-packed meshes.
			for _, cp := range item.Copack {

				done = timed(fmt.Sprintf("loading co-packed mesh %s", cp.Filename))
				coMesh, err := fauxgl.LoadMesh(cp.Filename)
				if err != nil {
					panic(err)
				}
				done()

				// IMPORTANT: cp.Scale is ignored. The main co-packing
				// mesh's scale is applied to all of its co-packed objects.
				if scale != 1.0 {
					done = timed("scaling main co-packing mesh")
					coMesh.Transform(scaleMatrix)
					done()
				}

				// add coMesh to the main mesh.
				mesh.Add(coMesh)
			}

			// update arrays with the main co-packing mesh's data for the json output.
			size := mesh.BoundingBox().Size()
			for i := 0; i < item.Count; i++ {
				singleStlSize = append(singleStlSize, size)
				srcStlNames = append(srcStlNames, item.Filename)
				scaleStl = append(scaleStl, scaleMatrix)
			}

			fmt.Printf("  %d triangles\n", len(mesh.Triangles))
			fmt.Printf("  %g x %g x %g\n", size.X, size.Y, size.Z)

			done = timed("centering co-packed mesh")
			mesh.Center()
			done()

			totalVolume += mesh.BoundingBox().Volume()
		}

		done = timed("building bvh tree")

		model.Add(mesh, bvhDetail, item.Count, spacing)
		ok = true
		done()

		fmt.Println("______________________________________________________")
	}

	if !ok {
		fmt.Println("Usage: pack3d --input_config_json_filename==mesh_config.json --output_packing_json_filename=export.json")
		fmt.Println(" - Packs N copies of each mesh into as small of a volume as possible.")
		fmt.Println(" - Runs forever, looking for the best packing.")
		fmt.Println(" - Results are written to disk whenever a new best is found.")
		return
	}

	side := math.Pow(totalVolume, 1.0/3)
	model.Deviation = side / 32 //it is not the distance between objects. And it seems that it will not reflect the distance.

	/*  Mesh packing loop. This loop is to find the best STL mesh packing.
	Add 'break' in the loop to stop program */
	start := time.Now()
	maxItemNum := len(model.Items)
	var timeLimit float64
	fillVolumeWithSpacing := 0.0
	totalFillVolume := 0.0
	null := fauxgl.Matrix{
		X00: 0, X01: 0, X02: 0, X03: 0,
		X10: 0, X11: 0, X12: 0, X13: 0,
		X20: 0, X21: 0, X22: 0, X23: 0,
		X30: 0, X31: 0, X32: 0, X33: 0,
	}
	timeLimit = 10

	minItemNum := 0
	packItemNum := maxItemNum
	success_model := pack3d.NewModel()

	for {
		model, ntime = model.Pack(annealingIterations, nil, singleStlSize, frameSize, packItemNum)
		/* ntime is the times of trial to find a output solution, if after trying for 100 times
		and no solution is found, then reset the model and try again. Usually if there is a solution,
		ntime will be 1 or 2 for most cases. */
		if ntime >= 100 {
			/* There is a case that even I reset the model for many times, I still can't find a solution,
			In this case, I need to set a threshold (20 second) to stop the software*/
			if time.Since(start).Seconds() <= timeLimit {
				model.Reset()
				continue
			} else {
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

				if minItemNum > maxItemNum {
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
					ioutil.WriteFile(fmt.Sprintf("%s.json", *fileNameArg), err_json, 0644)
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

		if minItemNum > maxItemNum {
			break
		}
		model.Reset()
	}

	done = timed("writing mesh")
	var (
		transMatrix    [4][4]float64
		fillPercentage float64
	)
	transformation := success_model.Transformation()
	for j := 0; j < len(success_model.Items); j++ {
		copack, ok := coPackMap[srcStlNames[j]]
		if !ok {

			t := transformation[j]
			st := t.Mul(scaleStl[j])  // scaled transformation for the j-th mesh.
			fillVolumeWithSpacing = (singleStlSize[j].X + spacing) * (singleStlSize[j].Y + spacing) * (singleStlSize[j].Z + spacing)
			if j < packItemNum {
				totalFillVolume += fillVolumeWithSpacing
				transMatrix = [4][4]float64{{st.X00, st.X01, st.X02, st.X03}, {st.X10, st.X11, st.X12, st.X13}, {st.X20, st.X21, st.X22, st.X23}, {st.X30, st.X31, st.X32, st.X33}}
			} else {
				transMatrix = [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}
			}

			// buildVolume's filling percentage.
			fillPercentage = totalFillVolume / buildVolume

			transMaps = append(transMaps, TransMap{srcStlNames[j], transMatrix, fillVolumeWithSpacing})

		} else {
			
			t := transformation[j]
			st := t.Mul(scaleStl[j])  // scaled transformation for the j-th mesh.
			fillVolumeWithSpacing = (singleStlSize[j].X + spacing) * (singleStlSize[j].Y + spacing) * (singleStlSize[j].Z + spacing)
			if j < packItemNum {
				totalFillVolume += fillVolumeWithSpacing
				transMatrix = [4][4]float64{{st.X00, st.X01, st.X02, st.X03}, {st.X10, st.X11, st.X12, st.X13}, {st.X20, st.X21, st.X22, st.X23}, {st.X30, st.X31, st.X32, st.X33}}
			} else {
				transMatrix = [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}
			}

			// buildVolume's filling percentage.
			fillPercentage = totalFillVolume / buildVolume

			// Add the main co-packing mesh to transMaps.
			transMaps = append(transMaps, TransMap{srcStlNames[j], transMatrix, fillVolumeWithSpacing})

			// Add the co-packed meshes to transMaps.
			for _, cp := range copack {
				// IMPORTANT: The volume of a co-packed object is already included in the volume
				//            of the parent object and the co-packed object's volume is set to 0.
				transMaps = append(transMaps, TransMap{cp.Filename, transMatrix, 0})
			}
		}
	}
	positions_json, err := json.Marshal(transMaps)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("the fill percentage is:", fillPercentage)
	ioutil.WriteFile(fmt.Sprintf("%s.json", *fileNameArg), positions_json, 0644)
	// os.Stdout.Write(positions_json)

	// STL file is no longer created, results returned as JSON for separate packer.
	// Unblock the following line if want to generate the packing STL
	// model.Mesh().SaveSTL(fmt.Sprintf("%s.stl", *fileNameArg))
	// model.TreeMesh().SaveSTL(fmt.Sprintf("out%dtree.stl", int(score*100000)))
	done()
}

type Config struct {
	BuildVolume [3]float64 `json:"build_volume"`
	Spacing     float64    `json:"spacing"`
	Items       []struct {
		Filename string    `json:"filename"`
		Scale    float64   `json:"scale"`
		Count    int       `json:"count"`
		Copack   []*Copack `json:"copack,omitempty"`
	} `json:"items"`
}

type Copack struct {
	Filename     string    `json:"filename"`
	// Scale        float64   `json:"scale"`
	// Transformation [4][4]float64 `json:"transformation"`  // ch32838 initially required this field then the requirements changed.
}
