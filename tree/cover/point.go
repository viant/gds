package cover

import "github.com/viant/bintly"

// Point represents a point in a vector space.
type Point struct {
	index     int32
	Magnitude float32
	Vector    []float32
}

func (p *Point) EncodeBinary(stream *bintly.Writer) error {
	stream.Int32(p.index)
	stream.Float32(p.Magnitude)
	stream.Float32s(p.Vector)
	return nil
}

func (p *Point) DecodeBinary(stream *bintly.Reader) error {
	stream.Int32(&p.index)
	stream.Float32(&p.Magnitude)
	stream.Float32s(&p.Vector)
	return nil

}

func (p *Point) HasValue() bool {
	return p.index != -1
}

func NewPoint(vector ...float32) *Point {
	p := &Point{Vector: vector}
	return p
}
