package cover

import "github.com/viant/vec/search"

// DistanceFunc is a function that calculates the distance between two points.
type DistanceFunc func(p1, p2 *Point) float32

// CosineDistance calculates the cosine distance between two points.
func CosineDistance(p1, p2 *Point) float32 {
	v1 := search.Float32s(p1.Vector)
	if p1.Magnitude == 0 {
		p1.Magnitude = v1.Magnitude()
	}
	v2 := search.Float32s(p2.Vector)
	if p2.Magnitude == 0 {
		p2.Magnitude = v2.Magnitude()
	}
	return v1.CosineDistanceWithMagnitude(p2.Vector, p1.Magnitude, p2.Magnitude)
}
