package pack3d

import "github.com/fogleman/fauxgl"

type Axis uint8

const (
	AxisNone Axis = iota
	AxisX
	AxisY
	AxisZ
)

func (a Axis) Vector() fauxgl.Vector {
	switch a {
	case AxisX:
		return fauxgl.Vector{1, 0, 0}
	case AxisY:
		return fauxgl.Vector{0, 1, 0}
	case AxisZ:
		return fauxgl.Vector{0, 0, 1}
	}
	return fauxgl.Vector{}
}
