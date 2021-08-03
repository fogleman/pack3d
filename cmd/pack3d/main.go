/*
Instruction: write down this into the command line to use the software
<pack3d {your_frame_size_X,your_frame_size_Y,your_frame_size_Z} spacing your_output_name STL_numbers STL_path STL_numbers STL_path .....>
For example: <pack3d {100,100,100} 5 Pikacu 1 /home/corner.stl 2 /home/Pika.stl>
The unit of frame and spacing is millimeters.
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
	"strconv"
	"strings"
	"time"

	"github.com/Authentise/pack3d/pack3d"
	"github.com/fogleman/fauxgl"
	"github.com/kr/pretty"
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
	var jsonFileArg = flag.String("json_file", "", "json config file")
	var buildVolumeArg = flag.String("build_volume", "100,100,100", "build volume")
	var spacingArg = flag.Float64("spacing", 2, "spacing")
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
		dimension     []float64
		ntime         int
		srcStlNames   []string
		transMaps     []TransMap
	)

	rand.Seed(time.Now().UTC().UnixNano())

	model := pack3d.NewModel()
	count := 1
	ok := false

	// Loading build_volume size
	for _, j := range strings.Split(*buildVolumeArg, ",") {
		_dimension, err := strconv.ParseFloat(j, 64)
		if err == nil {
			dimension = append(dimension, float64(_dimension))
			continue
		}
	}
	spacing := *spacingArg / 2.0
	// frameSize is the vertex in the first quadrant
	frameSize := fauxgl.V(dimension[0]/2.0, dimension[1]/2.0, dimension[2]/2.0)
	buildVolume := dimension[0] * dimension[1] * dimension[2]
	//fmt.Println(frameSize)

	meshMap := make(map[string]*fauxgl.Mesh)
	coPrintMap := make(map[string]*Coprint)
	/* Loading stl models */
	for _, item := range config.Items {
		done = timed(fmt.Sprintf("loading mesh %s", item.Filename))
		var mesh *fauxgl.Mesh
		var err error
		if item.Coprint == nil {
			mesh, err = fauxgl.LoadMesh(item.Filename)
			if err != nil {
				panic(err)
			}

			meshMap[item.Filename] = mesh.Copy()
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
			coPrintMap[item.Filename] = item.Coprint
			mesh, err = fauxgl.LoadMesh(item.Filename)
			if err != nil {
				panic(err)
			}

			size := mesh.BoundingBox().Size()

			meshMap[item.Filename] = mesh.Copy()

			coMesh, err := fauxgl.LoadMesh(item.Coprint.Filename)
			if err != nil {
				panic(err)
			}

			size = coMesh.BoundingBox().Size()

			coMesh.Transform(fauxgl.Matrix{
				X00: item.Coprint.Transformation[0][0],
				X01: item.Coprint.Transformation[0][1],
				X02: item.Coprint.Transformation[0][2],
				X03: item.Coprint.Transformation[0][3],
				X10: item.Coprint.Transformation[1][0],
				X11: item.Coprint.Transformation[1][1],
				X12: item.Coprint.Transformation[1][2],
				X13: item.Coprint.Transformation[1][3],
				X20: item.Coprint.Transformation[2][0],
				X21: item.Coprint.Transformation[2][1],
				X22: item.Coprint.Transformation[2][2],
				X23: item.Coprint.Transformation[2][3],
				X30: item.Coprint.Transformation[3][0],
				X31: item.Coprint.Transformation[3][1],
				X32: item.Coprint.Transformation[3][2],
				X33: item.Coprint.Transformation[3][3],
			})
			meshMap[item.Coprint.Filename] = coMesh.Copy()

			size = coMesh.BoundingBox().Size()

			mesh.Add(coMesh)

			done()

			totalVolume += mesh.BoundingBox().Volume()
			size = mesh.BoundingBox().Size()
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

	// for j := 0; j < len(model.Items); j++ {
	// 	fmt.Println("")
	// 	fmt.Println("##########")
	// 	pretty.Println(model.Items[j].Rotation, model.Items[j].Translation)
	// 	fmt.Println("##########")
	// 	fmt.Println("")
	// }

	if !ok {
		fmt.Println("Usage: pack3d frame_size output_fname N1 mesh1.stl N2 mesh2.stl ...")
		fmt.Println(" - Packs N copies of each mesh into as small of a volume as possible.")
		fmt.Println(" - Runs forever, looking for the best packing.")
		fmt.Println(" - Results are written to disk whenever a new best is found.")
		return
	}

	side := math.Pow(totalVolume, 1.0/3)
	model.Deviation = side / 32 //it is not the distance between objects. And it seems that it will not reflect the distance.

	/* This loop is to find the best packing stl, thus it will generate mutiple output
	Add 'break' in the loop to stop program */
	start := time.Now()
	maxItemNum := len(model.Items)
	var timeLimit float64
	fillVolumeWithSpacing := 0.0
	totalFillVolume := 0.0
	null := fauxgl.Matrix{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
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
		coprint, ok := coPrintMap[srcStlNames[j]]
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
			fmt.Println("")
			fmt.Println("##########")
			pretty.Println(transMatrix, success_model.Items[j].Rotation, success_model.Items[j].Translation)
			fmt.Println("##########")
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
			transMaps = append(transMaps, TransMap{srcStlNames[j], transMatrix, fillVolumeWithSpacing})

			m := meshMap[srcStlNames[j]]
			m.Transform(fauxgl.Translate(success_model.Items[j].Translation))
			t = pack3d.Rotations[success_model.Items[j].Rotation].Translate(success_model.Items[j].Translation)
			transMatrix = [4][4]float64{{t.X00, t.X01, t.X02, t.X03}, {t.X10, t.X11, t.X12, t.X13}, {t.X20, t.X21, t.X22, t.X23}, {t.X30, t.X31, t.X32, t.X33}}
			transMaps = append(transMaps, TransMap{srcStlNames[j], transMatrix, 0})

			m = meshMap[coprint.Filename]
			m.Transform(fauxgl.Translate(success_model.Items[j].Translation))
			t = pack3d.Rotations[success_model.Items[j].Rotation].Translate(success_model.Items[j].Translation)
			transMatrix = [4][4]float64{{t.X00, t.X01, t.X02, t.X03}, {t.X10, t.X11, t.X12, t.X13}, {t.X20, t.X21, t.X22, t.X23}, {t.X30, t.X31, t.X32, t.X33}}
			transMaps = append(transMaps, TransMap{coprint.Filename, transMatrix, 0})

			fmt.Println("")
			fmt.Println("##########")
			pretty.Println(transMatrix, success_model.Items[j].Rotation, success_model.Items[j].Translation)
			fmt.Println("##########")
		}
	}
	positions_json, err := json.Marshal(transMaps)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("the fill percentage is:", fillPercentage)
	ioutil.WriteFile(fmt.Sprintf("%s.json", *fileNameArg), positions_json, 0644)
	// os.Stdout.Write(positions_json)
	// model.Mesh().SaveSTL(fmt.Sprintf("pack3d-%s.stl", *fileNameArg))

	// STL file is no longer created
	// Unblock the following line if want to generate the packing STL
	/*
		model.Mesh().SaveSTL(fmt.Sprintf("pack3d-%s.stl", os.Args[5]))
		model.TreeMesh().SaveSTL(fmt.Sprintf("out%dtree.stl", int(score*100000)))
	*/
	done()
}

type Config struct {
	Items []struct {
		Filename string   `json:"filename"`
		Count    int      `json:"count"`
		Coprint  *Coprint `json:"coprint,omitempty"`
	} `json:"items"`
}

type Coprint struct {
	Filename       string        `json:"filename"`
	Transformation [4][4]float64 `json:"transformation"`
}
