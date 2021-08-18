/*
Instruction: write down this into the command line to use the software
<pack3d --json_file=json_config_file --filename=export_filename>
For example: <pack3d --json_file=input.json>
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
	var jsonFileArg = flag.String("json_file", "", "json config file")
	var fileNameArg = flag.String("filename", "pack3d", "export filename")
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
		done          func()
		totalVolume   float64
		ntime         int
		srcStlNames   []string
		transMaps     []TransMap
	)

	rand.Seed(time.Now().UTC().UnixNano())

	model := pack3d.NewModel()
	count := 1
	ok := false

	spacing := config.Spacing / 2.0
	// frameSize is the vertex in the first quadrant
	frameSize := fauxgl.V(config.BuildVolume[0]/2.0, config.BuildVolume[1]/2.0, config.BuildVolume[2]/2.0)
	buildVolume := config.BuildVolume[0] * config.BuildVolume[1] * config.BuildVolume[2]
	//fmt.Println(frameSize)


	/* Loading stl models */
	coPackMap := make(map[string][]*Copack)
	for _, item := range config.Items {
		done = timed(fmt.Sprintf("loading mesh %s", item.Filename))
		var mesh *fauxgl.Mesh
		var err error
		if item.Copack == nil {
			mesh, err = fauxgl.LoadMesh(item.Filename)
			if err != nil {
				panic(err)
			}

			done()

			totalVolume += mesh.BoundingBox().Volume()
			size := mesh.BoundingBox().Size()
			for i := 0; i < count; i++ {
				singleStlSize = append(singleStlSize, size)
				srcStlNames = append(srcStlNames, item.Filename)
			}

			fmt.Printf("  %d triangles\n", len(mesh.Triangles))
			fmt.Printf("  %g x %g x %g\n", size.X, size.Y, size.Z)

			done = timed("centering mesh")
			mesh.Center()
			done()
		} else {
			coPackMap[item.Filename] = item.Copack
			mesh, err = fauxgl.LoadMesh(item.Filename)
			if err != nil {
				panic(err)
			}
			for _, cp := range item.Copack {
				coMesh, err := fauxgl.LoadMesh(cp.Filename)
				if err != nil {
					panic(err)
				}

				coMesh.parent_id = ...
				coMesh.co_packing_transform(fauxgl.Matrix{
					X00: cp.Transformation[0][0],
					X01: cp.Transformation[0][1],
					X02: cp.Transformation[0][2],
					X03: cp.Transformation[0][3],
					X10: cp.Transformation[1][0],
					X11: cp.Transformation[1][1],
					X12: cp.Transformation[1][2],
					X13: cp.Transformation[1][3],
					X20: cp.Transformation[2][0],
					X21: cp.Transformation[2][1],
					X22: cp.Transformation[2][2],
					X23: cp.Transformation[2][3],
					X30: cp.Transformation[3][0],
					X31: cp.Transformation[3][1],
					X32: cp.Transformation[3][2],
					X33: cp.Transformation[3][3],
				})
				mesh.Add(coMesh)
			}
			done()

			totalVolume += mesh.BoundingBox().Volume()
			size := mesh.BoundingBox().Size()
			for i := 0; i < count; i++ {
				singleStlSize = append(singleStlSize, size)
				srcStlNames = append(srcStlNames, item.Filename)
			}

			fmt.Printf("  %d triangles\n", len(mesh.Triangles))
			fmt.Printf("  %g x %g x %g\n", size.X, size.Y, size.Z)

			done = timed("centering mesh")
			mesh.Center()
			done()
		}

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
		X30: 0, X31: 0, X32: 0, X33: 0
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
		for _, item := range config.Items {
		     if item.Copack != nil {
		         /*
		         	step 1: find out the parent's transformation (given by the model.Pack(...)), by using the parent_id or whatever else.

		            step 2: override the item.Transform (given by the packing algorithm) with:
		                    a matrix multiplication of the parent's transform (found in step 1) by the item.co_packing_transform.
		                    The order of the multiplication is meaningful and will be dealt with later. */
		     }
		}
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
			fillVolumeWithSpacing = (singleStlSize[j].X + spacing) * (singleStlSize[j].Y + spacing) * (singleStlSize[j].Z + spacing)
			if j < packItemNum {
				totalFillVolume += fillVolumeWithSpacing
				transMatrix = [4][4]float64{{t.X00, t.X01, t.X02, t.X03}, {t.X10, t.X11, t.X12, t.X13}, {t.X20, t.X21, t.X22, t.X23}, {t.X30, t.X31, t.X32, t.X33}}
			} else {
				transMatrix = [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}
			}
			// It's actually the bounding box filling percentage
			fillPercentage = totalFillVolume / buildVolume
			transMaps = append(transMaps, TransMap{srcStlNames[j], transMatrix, fillVolumeWithSpacing})
		} else {
			t := transformation[j]
			fillVolumeWithSpacing = (singleStlSize[j].X + spacing) * (singleStlSize[j].Y + spacing) * (singleStlSize[j].Z + spacing)
			if j < packItemNum {
				totalFillVolume += fillVolumeWithSpacing
				transMatrix = [4][4]float64{{t.X00, t.X01, t.X02, t.X03}, {t.X10, t.X11, t.X12, t.X13}, {t.X20, t.X21, t.X22, t.X23}, {t.X30, t.X31, t.X32, t.X33}}
			} else {
				transMatrix = [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}
			}
			// It's actually the bounding box filling percentage
			fillPercentage = totalFillVolume / buildVolume
			// transMaps = append(transMaps, TransMap{srcStlNames[j], transMatrix, fillVolumeWithSpacing})

			transMatrix = [4][4]float64{{t.X00, t.X01, t.X02, t.X03}, {t.X10, t.X11, t.X12, t.X13}, {t.X20, t.X21, t.X22, t.X23}, {t.X30, t.X31, t.X32, t.X33}}
			transMaps = append(transMaps, TransMap{srcStlNames[j], transMatrix, 0})
			for _, cp := range copack {
				transMatrix = [4][4]float64{
					{t.X00, t.X01, t.X02, t.X03 + cp.Transformation[0][3]},
					{t.X10, t.X11, t.X12, t.X13 + cp.Transformation[1][3]},
					{t.X20, t.X21, t.X22, t.X23 + cp.Transformation[2][3]},
					{t.X30, t.X31, t.X32, t.X33 + cp.Transformation[3][3]},
				}
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
	// model.Mesh().SaveSTL(fmt.Sprintf("pack3d-%s.stl", *fileNameArg))
	// model.TreeMesh().SaveSTL(fmt.Sprintf("out%dtree.stl", int(score*100000)))
	done()
}

type Config struct {
	BuildVolume [3]float64  `json:"build_volume"`
	Spacing     float64     `json:"spacing"`
	Items       []struct {
		Filename string    `json:"filename"`
		Count    int       `json:"count"`
		Copack  []*Copack  `json:"copack,omitempty"`
	} `json:"items"`
}

type Copack struct {
	Filename       string        `json:"filename"`
	Transformation [4][4]float64 `json:"transformation"`
}
