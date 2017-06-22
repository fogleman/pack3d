package binpack

type Axis int

const (
	AxisX Axis = iota
	AxisY
	AxisZ
)

type Item struct {
	ID    int
	Score int
	Size  Vector
}

type Box struct {
	Origin Vector
	Size   Vector
}

func (box Box) Cut(axis Axis, offset int) (Box, Box) {
	o1 := box.Origin
	o2 := box.Origin
	s1 := box.Size
	s2 := box.Size
	switch axis {
	case AxisX:
		s1.X = offset
		s2.X -= offset
		o2.X += offset
	case AxisY:
		s1.Y = offset
		s2.Y -= offset
		o2.Y += offset
	case AxisZ:
		s1.Z = offset
		s2.Z -= offset
		o2.Z += offset
	}
	return Box{o1, s1}, Box{o2, s2}
}

func (box Box) Cuts(a1, a2, a3 Axis, s1, s2, s3 int) (Box, Box, Box) {
	b := box
	b, b1 := b.Cut(a1, s1)
	b, b2 := b.Cut(a2, s2)
	_, b3 := b.Cut(a3, s3)
	return b1, b2, b3
}

type Placement struct {
	Item     Item
	Position Vector
}

type Result struct {
	Score      int
	Placements []Placement
}

func MakeResult(r0, r1, r2 Result, item Item, position Vector) Result {
	r3 := Result{item.Score, []Placement{{item, position}}}
	score := r0.Score + r1.Score + r2.Score + r3.Score
	n := len(r0.Placements) + len(r1.Placements) + len(r2.Placements) + len(r3.Placements)
	placements := make([]Placement, 0, n)
	placements = append(placements, r0.Placements...)
	placements = append(placements, r1.Placements...)
	placements = append(placements, r2.Placements...)
	placements = append(placements, r3.Placements...)
	return Result{score, placements}
}

func (result Result) Translate(offset Vector) Result {
	placements := make([]Placement, len(result.Placements))
	for i, p := range result.Placements {
		p.Position = p.Position.Add(offset)
		placements[i] = p
	}
	return Result{result.Score, placements}
}

func Pack(items []Item, box Box) Result {
	hash := NewSpatialHash(1000)
	minVolume := items[0].Size.Sort()
	for _, item := range items {
		minVolume = minVolume.Min(item.Size.Sort())
	}
	return pack(items, box, hash, minVolume)
}

func pack(items []Item, box Box, hash *SpatialHash, minVolume Vector) Result {
	bs := box.Size
	if !bs.Sort().Fits(minVolume) {
		return Result{}
	}
	if result, ok := hash.Get(bs); ok {
		return result
	}
	best := Result{}
	for _, item := range items {
		s := item.Size
		if s.X > bs.X || s.Y > bs.Y || s.Z > bs.Z {
			continue
		}
		var b [6][3]Box
		b[0][0], b[0][1], b[0][2] = box.Cuts(AxisX, AxisY, AxisZ, s.X, s.Y, s.Z)
		b[1][0], b[1][1], b[1][2] = box.Cuts(AxisX, AxisZ, AxisY, s.X, s.Z, s.Y)
		b[2][0], b[2][1], b[2][2] = box.Cuts(AxisY, AxisX, AxisZ, s.Y, s.X, s.Z)
		b[3][0], b[3][1], b[3][2] = box.Cuts(AxisY, AxisZ, AxisX, s.Y, s.Z, s.X)
		b[4][0], b[4][1], b[4][2] = box.Cuts(AxisZ, AxisX, AxisY, s.Z, s.X, s.Y)
		b[5][0], b[5][1], b[5][2] = box.Cuts(AxisZ, AxisY, AxisX, s.Z, s.Y, s.X)
		for i := 0; i < 6; i++ {
			var r [3]Result
			score := item.Score
			for j := 0; j < 3; j++ {
				r[j] = pack(items, b[i][j], hash, minVolume)
				score += r[j].Score
			}
			if score > best.Score {
				for j := 0; j < 3; j++ {
					r[j] = r[j].Translate(b[i][j].Origin)
				}
				best = MakeResult(r[0], r[1], r[2], item, box.Origin)
			}
		}
	}
	best = best.Translate(box.Origin.Negate())
	var size Vector
	for _, p := range best.Placements {
		size = size.Max(p.Position.Add(p.Item.Size))
	}
	hash.Add(size, bs, best)
	return best
}
