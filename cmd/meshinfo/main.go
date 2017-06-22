package main

import (
	"fmt"
	"os"

	. "github.com/fogleman/fauxgl"
)

func main() {
	for _, path := range os.Args[1:] {
		fmt.Println(path)
		mesh, err := LoadMesh(path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		box := mesh.BoundingBox()
		size := box.Size()
		volume := size.X * size.Y * size.Z
		center := box.Anchor(V(0.5, 0.5, 0.5))
		fmt.Printf("  triangles = %d\n", len(mesh.Triangles))
		fmt.Printf("  x range   = %g to %g\n", box.Min.X, box.Max.X)
		fmt.Printf("  y range   = %g to %g\n", box.Min.Y, box.Max.Y)
		fmt.Printf("  z range   = %g to %g\n", box.Min.Z, box.Max.Z)
		fmt.Printf("  center    = %g, %g, %g\n", center.X, center.Y, center.Z)
		fmt.Printf("  size      = %g x %g x %g\n", size.X, size.Y, size.Z)
		fmt.Printf("  volume    = %g\n", volume)
		fmt.Println()
	}
}
