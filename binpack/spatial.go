package binpack

type SpatialHash struct {
	CellSize int
	Cells    map[SpatialKey][]*SpatialValue
}

type SpatialKey struct {
	X, Y, Z int
}

type SpatialValue struct {
	Min, Max Vector
	Result   Result
}

func NewSpatialHash(cellSize int) *SpatialHash {
	cells := make(map[SpatialKey][]*SpatialValue)
	return &SpatialHash{cellSize, cells}
}

func (h *SpatialHash) KeyForVector(v Vector) SpatialKey {
	x := v.X / h.CellSize
	y := v.Y / h.CellSize
	z := v.Z / h.CellSize
	return SpatialKey{x, y, z}
}

func (h *SpatialHash) Add(min, max Vector, result Result) {
	value := &SpatialValue{min, max, result}
	k1 := h.KeyForVector(min)
	k2 := h.KeyForVector(max)
	for x := k1.X; x <= k2.X; x++ {
		for y := k1.Y; y <= k2.Y; y++ {
			for z := k1.Z; z <= k2.Z; z++ {
				k := SpatialKey{x, y, z}
				h.Cells[k] = append(h.Cells[k], value)
			}
		}
	}
}

func (h *SpatialHash) Get(v Vector) (Result, bool) {
	k := h.KeyForVector(v)
	for _, value := range h.Cells[k] {
		if v.GreaterThanOrEqual(value.Min) && v.LessThanOrEqual(value.Max) {
			return value.Result, true
		}
	}
	return Result{}, false
}
