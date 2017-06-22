package main

import (
	"fmt"
	"math"
	"os"
	"path"
	"strconv"

	. "github.com/fogleman/fauxgl"
	"github.com/fogleman/pack3d/pack3d"
)

func main() {
	detail := 8
	for _, arg := range os.Args[1:] {
		_detail, err := strconv.ParseInt(arg, 0, 0)
		if err == nil {
			detail = int(_detail)
			continue
		}
		mesh, err := LoadMesh(arg)
		if err != nil {
			panic(err)
		}
		tree := pack3d.NewTreeForMesh(mesh, detail)
		mesh = NewEmptyMesh()
		n := int(math.Pow(2, float64(detail)))
		for _, box := range tree[len(tree)-n:] {
			mesh.Add(NewCubeForBox(box))
		}
		ext := path.Ext(arg)
		mesh.SaveSTL(fmt.Sprintf(arg[:len(arg)-len(ext)]+".bvh.%d.stl", detail))
	}
}
