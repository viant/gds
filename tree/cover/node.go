package cover

import "math"

// Node represents a node in the cover tree.
type Node struct {
	level     int
	baseLevel float32
	children  []Node
	point     *Point
}

func NewNode(point *Point, level int, base float32) Node {
	return Node{
		level:     level,
		baseLevel: float32(math.Pow(float64(base), float64(level))),
		point:     point,
	}
}
