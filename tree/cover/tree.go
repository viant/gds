package cover

import (
	"container/heap"
	"math"
)

// Tree represents a cover tree.
type Tree[T any] struct {
	root       *Node
	pointSize  int //point vector size
	base       float32
	distanceFn DistanceFunc
	values     values[T]
}

// Insert adds a new point (embedding vector) to the cover tree.
func (t *Tree[T]) Insert(point *Point, value T) {
	point.index = t.values.put(value)
	if t.root == nil {
		t.root = &Node{point: point, level: 0}
	} else {
		t.insert(t.root, point, 0)
	}
}

func (t *Tree[T]) insert(node *Node, point *Point, level int) {
	baseLevel := float32(math.Pow(float64(t.base), float64(level)))
	if t.distanceFn(point, node.point) < baseLevel {
		for i := range node.children {
			child := &node.children[i]
			t.insert(child, point, level-1)
			return
		}
		node.children = append(node.children, NewNode(point, level-1, t.base))
		return
	}
	if t.distanceFn(point, node.point) < node.baseLevel {
		node.children = append(node.children, NewNode(point, level-1, t.base))
		return
	}
	for level > node.level && t.distanceFn(point, node.point) >= baseLevel {
		level++
	}
	t.insert(t.root, point, level)
}

// Remove removes a point (embedding vector) from the cover tree.
func (t *Tree[T]) Remove(point *Point) bool {
	if t.root == nil {
		return false
	}
	removed, newRoot := t.remove(t.root, point)
	t.root = newRoot
	return removed
}

func (t *Tree[T]) remove(node *Node, point *Point) (bool, *Node) {
	if node == nil {
		return false, nil
	}
	if t.distanceFn(point, node.point) == 0 {
		if len(node.children) == 0 {
			return true, nil
		}

		// Promote one of the children to be the new node
		newNode := &node.children[0]
		for _, child := range node.children[1:] {
			t.insert(newNode, child.point, child.level)
		}
		return true, newNode
	}
	for i := range node.children {
		child := &node.children[i]
		removed, newChild := t.remove(child, point)
		if removed {
			if newChild == nil {
				node.children = append(node.children[:i], node.children[i+1:]...)
			} else {
				node.children[i] = *newChild
			}
			return true, node
		}
	}
	return false, node
}

func (t *Tree[T]) Value(point *Point) T {
	var r T
	if point == nil || !point.HasValue() {
		return r
	}
	return t.values.value(point.index)
}

func (t *Tree[T]) Values(points []*Point) []T {
	var result = make([]T, 0, len(points))
	for i, point := range points {
		if point == nil || point.index < 0 {
			continue
		}
		result[i] = t.values.value(point.index)
	}
	return result
}

// KNearestNeighbors finds the k nearest neighbors of the given point (embedding vector) in the cover tree.
func (t *Tree[T]) KNearestNeighbors(point *Point, k int) []*Point {
	if t.root == nil {
		return nil
	}
	h := &Neighbors{}
	heap.Init(h)
	t.kNearestNeighbors(t.root, point, k, h)
	result := make([]*Point, h.Len())
	for i := len(result) - 1; i >= 0; i-- {
		result[i] = heap.Pop(h).(Neighbor).Point
	}
	return result
}

func (t *Tree[T]) kNearestNeighbors(node *Node, point *Point, k int, h *Neighbors) {
	dist := t.distanceFn(point, node.point)
	if h.Len() < k {
		heap.Push(h, Neighbor{Point: node.point, Distance: dist})
	} else if dist < (*h)[0].Distance {
		heap.Pop(h)
		heap.Push(h, Neighbor{Point: node.point, Distance: dist})
	}
	for i := range node.children {
		t.kNearestNeighbors(&node.children[i], point, k, h)
	}
}

// NewTree initializes and returns a new Tree.
func NewTree[T any](base float32, distanceFn DistanceFunc) *Tree[T] {
	return &Tree[T]{base: base, distanceFn: distanceFn, values: values[T]{data: make([]T, 0)}}
}
