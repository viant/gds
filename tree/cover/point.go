package cover

// Point represents a point in a vector space.
type Point struct {
	index     int32
	Magnitude float32
	Vector    []float32
}

func (p *Point) HasValue() bool {
	return p.index != -1
}
