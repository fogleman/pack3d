package main

import (
	"fmt"
	"time"

	. "github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"
)

const (
	N      = 48
	scale  = 4     // optional supersampling
	width  = 1200  // output width in pixels
	height = 1600  // output height in pixels
	fovy   = 18    // vertical field of view in degrees
	near   = 100   // near clipping plane
	far    = 10000 // far clipping plane
)

var (
	eye    = V(1000, 1000, 160)           // camera position
	center = V(165/2.0, 165/2.0, 160)     // view center position
	up     = V(0, 0, 1)                   // up vector
	light  = V(0.75, 0.25, 1).Normalize() // light direction
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

	done = timed("loading mesh")
	fuse, err := LoadMesh("fuse2.stl")
	if err != nil {
		panic(err)
	}
	done()

	// load a mesh
	done = timed("loading mesh")
	mesh, err := LoadMesh("giraffe48.stl")
	if err != nil {
		panic(err)
	}
	done()

	t := len(mesh.Triangles) / N

	// create a rendering context
	context := NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(HexColor("#FFFFFF"))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// render
	shader := NewPhongShader(matrix, light, eye)
	context.Shader = shader

	done = timed("rendering fuse")
	shader.ObjectColor = HexColor("#2A2C2B")
	context.DrawMesh(fuse)
	done()

	for i := 0; i <= N; i++ {
		if i > 0 {
			j := N - i + 1
			done = timed("rendering boats")
			shader.ObjectColor = HexColor("#7E827A")
			context.DrawTriangles(mesh.Triangles[t*(j-1) : t*j])
			done()
		}

		// downsample image for antialiasing
		done = timed("downsampling image")
		image := context.Image()
		image = resize.Resize(width, height, image, resize.Bilinear)
		done()

		// save image
		done = timed("writing output")
		SavePNG(fmt.Sprintf("out%03d.png", i), image)
		done()
	}
}
