package main

import (
	"fmt"
	"os"
	"time"

	. "github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"
)

const (
	scale  = 4    // optional supersampling
	width  = 2048 // output width in pixels
	height = 2048 // output height in pixels
	fovy   = 35   // vertical field of view in degrees
	near   = 1    // near clipping plane
	far    = 1000 // far clipping plane
)

var (
	eye        = V(100, 200, 100)             // camera position
	center     = V(0, 0, 0)                   // view center position
	up         = V(0, 0, 1)                   // up vector
	light      = V(0.75, 1, 0.25).Normalize() // light direction
	color      = HexColor("#468966")          // object color
	background = HexColor("#FFF8E3")          // background color
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

func swap(v Vector, n int) Vector {
	switch n {
	case 0:
		return Vector{v.X, v.Y, v.Z}
	case 1:
		return Vector{v.X, v.Z, v.Y}
	case 2:
		return Vector{v.Y, v.X, v.Z}
	case 3:
		return Vector{v.Y, v.Z, v.X}
	case 4:
		return Vector{v.Z, v.X, v.Y}
	case 5:
		return Vector{v.Z, v.Y, v.X}
	}
	return v
}

func main() {
	var done func()

	for _, path := range os.Args[1:] {
		for i := 0; i < 6; i++ {

			// load a mesh
			done = timed("loading mesh")
			mesh, err := LoadMesh(path)
			if err != nil {
				panic(err)
			}
			done()
			// mesh.Transform(Rotate(up, Radians(180)))

			// fit mesh in a bi-unit cube centered at the origin
			done = timed("transforming mesh")
			mesh.MoveTo(Vector{}, Vector{0.5, 0.5, 0.5})
			done()

			_eye := swap(eye, i)
			_center := swap(center, i)
			_up := swap(up, i)
			_light := swap(light, i)

			// create a rendering context
			context := NewContext(width*scale, height*scale)
			context.ClearColorBufferWith(background)

			// create transformation matrix and light direction
			aspect := float64(width) / float64(height)
			matrix := LookAt(_eye, _center, _up).Perspective(fovy, aspect, near, far)

			// render
			shader := NewPhongShader(matrix, _light, _eye)
			shader.ObjectColor = color
			context.Shader = shader
			done = timed("rendering mesh")
			context.DrawMesh(mesh)
			done()

			context.Shader = NewSolidColorShader(matrix, Black)
			context.LineWidth = scale * 3
			context.DrawMesh(NewCubeOutlineForBox(mesh.BoundingBox()))

			// downsample image for antialiasing
			done = timed("downsampling image")
			image := context.Image()
			image = resize.Resize(width, height, image, resize.Bilinear)
			done()

			// save image
			done = timed("writing output")
			SavePNG(fmt.Sprintf("%s.%d.png", path, i), image)
			done()
		}
	}
}
