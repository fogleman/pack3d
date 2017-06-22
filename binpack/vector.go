package binpack

type Vector struct {
	X, Y, Z int
}

func (a Vector) Add(b Vector) Vector {
	return Vector{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func (a Vector) Sub(b Vector) Vector {
	return Vector{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func (a Vector) Negate() Vector {
	return Vector{-a.X, -a.Y, -a.Z}
}

func (a Vector) Min(b Vector) Vector {
	if b.X < a.X {
		a.X = b.X
	}
	if b.Y < a.Y {
		a.Y = b.Y
	}
	if b.Z < a.Z {
		a.Z = b.Z
	}
	return a
}

func (a Vector) Max(b Vector) Vector {
	if b.X > a.X {
		a.X = b.X
	}
	if b.Y > a.Y {
		a.Y = b.Y
	}
	if b.Z > a.Z {
		a.Z = b.Z
	}
	return a
}

func (a Vector) Sort() Vector {
	if a.X > a.Z {
		a.X, a.Z = a.Z, a.X
	}
	if a.X > a.Y {
		a.X, a.Y = a.Y, a.X
	}
	if a.Y > a.Z {
		a.Y, a.Z = a.Z, a.Y
	}
	return a
}

func (a Vector) Fits(b Vector) bool {
	return a.X >= b.X && a.Y >= b.Y && a.Z >= b.Z
}

func (a Vector) GreaterThanOrEqual(b Vector) bool {
	return a.X >= b.X && a.Y >= b.Y && a.Z >= b.Z
}

func (a Vector) LessThanOrEqual(b Vector) bool {
	return a.X <= b.X && a.Y <= b.Y && a.Z <= b.Z
}
