package pack3d

import (
	"math"
	"math/rand"

	"github.com/fogleman/fauxgl"
)

var Rotations []fauxgl.Matrix

func init() {
	for i := 0; i < 4; i++ {
		for s := -1; s <= 1; s += 2 {
			for a := 1; a <= 3; a++ {
				up := AxisZ.Vector()
				m := fauxgl.Rotate(up, float64(i)*fauxgl.Radians(90))
				m = m.RotateTo(up, Axis(a).Vector().MulScalar(float64(s)))
				Rotations = append(Rotations, m)
			}
		}
	}
}

type Undo struct {
	Index       int
	Rotation    int
	Translation fauxgl.Vector
}

type Item struct {
	Mesh        *fauxgl.Mesh
	Trees       []Tree
	Rotation    int
	Translation fauxgl.Vector
}

func (item *Item) Matrix() fauxgl.Matrix {
	return Rotations[item.Rotation].Translate(item.Translation)
}

func (item *Item) Copy() *Item {
	dup := *item
	return &dup
}

type Model struct {
	Items     []*Item
	MinVolume float64
	MaxVolume float64
	Deviation float64
}

func NewModel() *Model {
	return &Model{nil, 0, 0, 1}
}

func (m *Model) Add(mesh *fauxgl.Mesh, detail, count int) {
	tree := NewTreeForMesh(mesh, detail)
	trees := make([]Tree, len(Rotations))
	for i, m := range Rotations {
		trees[i] = tree.Transform(m)
	}
	for i := 0; i < count; i++ {
		m.add(mesh, trees)
	}
}

func (m *Model) add(mesh *fauxgl.Mesh, trees []Tree) {
	index := len(m.Items)
	item := Item{mesh, trees, 0, fauxgl.Vector{}}
	m.Items = append(m.Items, &item)
	d := 1.0
	for !m.ValidChange(index) {
		item.Rotation = rand.Intn(len(Rotations))
		item.Translation = fauxgl.RandomUnitVector().MulScalar(d)
		d *= 1.2
	}
	tree := trees[0]
	m.MinVolume = math.Max(m.MinVolume, tree[0].Volume())
	m.MaxVolume += tree[0].Volume()
}

func (m *Model) Reset() {
	items := m.Items
	m.Items = nil
	m.MinVolume = 0
	m.MaxVolume = 0
	for _, item := range items {
		m.add(item.Mesh, item.Trees)
	}
}

func (m *Model) Pack(iterations int, callback AnnealCallback) *Model {
	e := 0.5
	return Anneal(m, 1e0*e, 1e-4*e, iterations, callback).(*Model)
}

func (m *Model) Meshes() []*fauxgl.Mesh {
	result := make([]*fauxgl.Mesh, len(m.Items))
	for i, item := range m.Items {
		mesh := item.Mesh.Copy()
		mesh.Transform(item.Matrix())
		result[i] = mesh
	}
	return result
}

func (m *Model) Mesh() *fauxgl.Mesh {
	result := fauxgl.NewEmptyMesh()
	for _, mesh := range m.Meshes() {
		result.Add(mesh)
	}
	return result
}

func (m *Model) TreeMeshes() []*fauxgl.Mesh {
	result := make([]*fauxgl.Mesh, len(m.Items))
	for i, item := range m.Items {
		mesh := fauxgl.NewEmptyMesh()
		tree := item.Trees[item.Rotation]
		for _, box := range tree[len(tree)/2:] {
			mesh.Add(fauxgl.NewCubeForBox(box))
		}
		mesh.Transform(fauxgl.Translate(item.Translation))
		result[i] = mesh
	}
	return result
}

func (m *Model) TreeMesh() *fauxgl.Mesh {
	result := fauxgl.NewEmptyMesh()
	for _, mesh := range m.TreeMeshes() {
		result.Add(mesh)
	}
	return result
}

func (m *Model) ValidChange(i int) bool {
	item1 := m.Items[i]
	tree1 := item1.Trees[item1.Rotation]
	for j := 0; j < len(m.Items); j++ {
		if j == i {
			continue
		}
		item2 := m.Items[j]
		tree2 := item2.Trees[item2.Rotation]
		if tree1.Intersects(tree2, item1.Translation, item2.Translation) {
			return false
		}
	}
	return true
}

func (m *Model) BoundingBox() fauxgl.Box {
	box := fauxgl.EmptyBox
	for _, item := range m.Items {
		tree := item.Trees[item.Rotation]
		box = box.Extend(tree[0].Translate(item.Translation))
	}
	return box
}

func (m *Model) Volume() float64 {
	return m.BoundingBox().Volume()
}

func (m *Model) Energy() float64 {
	return m.Volume() / m.MaxVolume
}

func (m *Model) DoMove() Undo {
	i := rand.Intn(len(m.Items))
	item := m.Items[i]
	undo := Undo{i, item.Rotation, item.Translation}
	for {
		if rand.Intn(4) == 0 {
			// rotate
			item.Rotation = rand.Intn(len(Rotations))
		} else {
			// translate
			offset := Axis(rand.Intn(3) + 1).Vector()
			offset = offset.MulScalar(rand.NormFloat64() * m.Deviation)
			item.Translation = item.Translation.Add(offset)
		}
		if m.ValidChange(i) {
			break
		}
		item.Rotation = undo.Rotation
		item.Translation = undo.Translation
	}
	return undo
}

func (m *Model) UndoMove(undo Undo) {
	item := m.Items[undo.Index]
	item.Rotation = undo.Rotation
	item.Translation = undo.Translation
}

func (m *Model) Copy() Annealable {
	items := make([]*Item, len(m.Items))
	for i, item := range m.Items {
		items[i] = item.Copy()
	}
	return &Model{items, m.MinVolume, m.MaxVolume, m.Deviation}
}
