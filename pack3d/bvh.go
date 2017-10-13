package pack3d

import "github.com/fogleman/fauxgl"

type Tree []fauxgl.Box

func NewTreeForMesh(mesh *fauxgl.Mesh, depth int) Tree {
	mesh = mesh.Copy()
	mesh.Center()
	boxes := make([]fauxgl.Box, len(mesh.Triangles))
	for i, t := range mesh.Triangles {
		boxes[i] = t.BoundingBox()
	}
	root := NewNode(boxes, depth)
	tree := make(Tree, 1<<uint(depth+1)-1)
	root.Flatten(tree, 0)
	return tree
}

func (a Tree) Transform(m fauxgl.Matrix) Tree {
	b := make(Tree, len(a))
	for i, box := range a {
		b[i] = box.Transform(m)
	}
	return b
}

func (a Tree) Intersects(b Tree, t1, t2 fauxgl.Vector) bool {
	return a.intersects(b, t1, t2, 0, 0)
}

func (a Tree) intersects(b Tree, t1, t2 fauxgl.Vector, i, j int) bool {
	if !boxesIntersect(a[i], b[j], t1, t2) {
		return false
	}
	i1 := i*2 + 1
	i2 := i*2 + 2
	j1 := j*2 + 1
	j2 := j*2 + 2
	if i1 >= len(a) && j1 >= len(b) {
		return true
	} else if i1 >= len(a) {
		return a.intersects(b, t1, t2, i, j1) || a.intersects(b, t1, t2, i, j2)
	} else if j1 >= len(b) {
		return a.intersects(b, t1, t2, i1, j) || a.intersects(b, t1, t2, i2, j)
	} else {
		return a.intersects(b, t1, t2, i1, j1) ||
			a.intersects(b, t1, t2, i1, j2) ||
			a.intersects(b, t1, t2, i2, j1) ||
			a.intersects(b, t1, t2, i2, j2)
	}
}

func boxesIntersect(b1, b2 fauxgl.Box, t1, t2 fauxgl.Vector) bool {
	if b1 == fauxgl.EmptyBox || b2 == fauxgl.EmptyBox {
		return false
	}
	return !(b1.Min.X+t1.X > b2.Max.X+t2.X ||
		b1.Max.X+t1.X < b2.Min.X+t2.X ||
		b1.Min.Y+t1.Y > b2.Max.Y+t2.Y ||
		b1.Max.Y+t1.Y < b2.Min.Y+t2.Y ||
		b1.Min.Z+t1.Z > b2.Max.Z+t2.Z ||
		b1.Max.Z+t1.Z < b2.Min.Z+t2.Z)
}

type Node struct {
	Box   fauxgl.Box
	Left  *Node
	Right *Node
}

func NewNode(boxes []fauxgl.Box, depth int) *Node {
	box := fauxgl.BoxForBoxes(boxes).Offset(2.5)
	node := &Node{box, nil, nil}
	node.Split(boxes, depth)
	return node
}

func (a *Node) Flatten(tree Tree, i int) {
	tree[i] = a.Box
	if a.Left != nil {
		a.Left.Flatten(tree, i*2+1)
	}
	if a.Right != nil {
		a.Right.Flatten(tree, i*2+2)
	}
}

func (node *Node) Split(boxes []fauxgl.Box, depth int) {
	if depth == 0 {
		return
	}
	box := node.Box
	best := box.Volume()
	bestAxis := AxisNone
	bestPoint := 0.0
	bestSide := false
	const N = 16
	for s := 0; s < 2; s++ {
		side := s == 1
		for i := 1; i < N; i++ {
			p := float64(i) / N
			x := box.Min.X + (box.Max.X-box.Min.X)*p
			y := box.Min.Y + (box.Max.Y-box.Min.Y)*p
			z := box.Min.Z + (box.Max.Z-box.Min.Z)*p
			sx := partitionScore(boxes, AxisX, x, side)
			if sx < best {
				best = sx
				bestAxis = AxisX
				bestPoint = x
				bestSide = side
			}
			sy := partitionScore(boxes, AxisY, y, side)
			if sy < best {
				best = sy
				bestAxis = AxisY
				bestPoint = y
				bestSide = side
			}
			sz := partitionScore(boxes, AxisZ, z, side)
			if sz < best {
				best = sz
				bestAxis = AxisZ
				bestPoint = z
				bestSide = side
			}
		}
	}
	if bestAxis == AxisNone {
		return
	}
	l, r := partition(boxes, bestAxis, bestPoint, bestSide)
	node.Left = NewNode(l, depth-1)
	node.Right = NewNode(r, depth-1)
}

func partitionBox(box fauxgl.Box, axis Axis, point float64) (left, right bool) {
	switch axis {
	case AxisX:
		left = box.Min.X <= point
		right = box.Max.X >= point
	case AxisY:
		left = box.Min.Y <= point
		right = box.Max.Y >= point
	case AxisZ:
		left = box.Min.Z <= point
		right = box.Max.Z >= point
	}
	return
}

func partitionScore(boxes []fauxgl.Box, axis Axis, point float64, side bool) float64 {
	var major fauxgl.Box
	for _, box := range boxes {
		l, r := partitionBox(box, axis, point)
		if (l && r) || (l && side) || (r && !side) {
			major = major.Extend(box)
		}
	}
	var minor fauxgl.Box
	for _, box := range boxes {
		if !major.ContainsBox(box) {
			minor = minor.Extend(box)
		}
	}
	return major.Volume() + minor.Volume() - major.Intersection(minor).Volume()
}

func partition(boxes []fauxgl.Box, axis Axis, point float64, side bool) (left, right []fauxgl.Box) {
	var major fauxgl.Box
	for _, box := range boxes {
		l, r := partitionBox(box, axis, point)
		if (l && r) || (l && side) || (r && !side) {
			major = major.Extend(box)
		}
	}
	for _, box := range boxes {
		if major.ContainsBox(box) {
			left = append(left, box)
		} else {
			right = append(right, box)
		}
	}
	if !side {
		left, right = right, left
	}
	return
}
