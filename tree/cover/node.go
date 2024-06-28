package cover

import (
	"github.com/viant/bintly"
	"math"
)

// Node represents a node in the cover tree.
type Node struct {
	level     int32
	baseLevel float32
	point     *Point
	children  []Node
}

func (n *Node) EncodeBinary(stream *bintly.Writer) error {
	stream.Int32(n.level)
	stream.Float32(n.baseLevel)
	if err := stream.Coder(n.point); err != nil {
		return err
	}
	stream.Int32(int32(len(n.children)))
	for i := range n.children {
		if err := stream.Coder(&n.children[i]); err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) DecodeBinary(stream *bintly.Reader) error {
	stream.Int32(&n.level)
	stream.Float32(&n.baseLevel)
	n.point = &Point{}
	if err := stream.Coder(n.point); err != nil {
		return err
	}
	var size int32
	stream.Int32(&size)
	n.children = make([]Node, 0)
	for i := int32(0); i < size; i++ {
		child := Node{}
		if err := stream.Coder(&child); err != nil {
			return err
		}
		n.children = append(n.children, child)
	}
	return nil
}

func NewNode(point *Point, level int32, base float32) Node {
	return Node{
		level:     level,
		baseLevel: float32(math.Pow(float64(base), float64(level))),
		point:     point,
	}
}
