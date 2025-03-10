package cover

import (
	"container/heap"
	"io"
	"math"
)

// Tree represents a cover tree.
type Tree[T any] struct {
	root             *Node
	base             float32
	distanceFuncName DistanceFunction
	distanceFnunc    DistanceFunc
	values           values[T]
	indexMap         map[int32]*Point
}

// Insert adds a new point (embedding vector) to the cover tree.
func (t *Tree[T]) Insert(value T, point *Point) int32 {
	point.index = t.values.put(value)
	if t.indexMap == nil {
		t.indexMap = make(map[int32]*Point)
	}
	t.indexMap[point.index] = point //
	if t.root == nil {
		t.root = &Node{point: point, level: 0}
	} else {
		t.insert(t.root, point, 0)
	}
	return point.index
}

// FindPointByIndex returns the point associated with the given index.
func (t *Tree[T]) FindPointByIndex(index int32) *Point {
	if point, exists := t.indexMap[index]; exists {
		return point
	}
	return nil
}

func (t *Tree[T]) EncodeValues(writer io.Writer) error {
	return t.values.Encode(writer)
}

func (t *Tree[T]) DecodeValues(reader io.Reader) error {
	t.values = values[T]{data: make([]T, 0)}
	t.values.ensureType()
	return t.values.Decode(reader)
}

func (t *Tree[T]) EncodeTree(writer io.Writer) error {
	buffer := writers.Get()
	defer writers.Put(buffer)
	buffer.Float32(t.base)

	buffer.String(string(t.distanceFuncName))
	if err := buffer.Coder(t.root); err != nil {
		return err
	}
	_, err := writer.Write(buffer.Bytes())
	return err
}

func (t *Tree[T]) DecodeTree(reader io.Reader) error {
	buffer := readers.Get()
	defer readers.Put(buffer)
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if err = buffer.FromBytes(data); err != nil {
		return err
	}
	buffer.Float32(&t.base)
	var distance string
	buffer.String(&distance)
	t.distanceFuncName = DistanceFunction(distance)
	t.distanceFnunc = t.distanceFuncName.Function()
	t.root = &Node{}
	return buffer.Coder(t.root)
}

func (t *Tree[T]) insert(node *Node, point *Point, level int32) {
	for {
		baseLevel := float32(math.Pow(float64(t.base), float64(level)))
		distance := t.distanceFnunc(point, node.point)

		if distance < baseLevel {
			// Check if the point can be inserted into any child
			inserted := false
			for i := range node.children {
				child := &node.children[i]
				if t.distanceFnunc(point, child.point) < baseLevel {
					node = child
					level--
					inserted = true
					break
				}
			}
			if !inserted {
				// Insert as a new child
				node.children = append(node.children, NewNode(point, level-1, t.base))
				return
			}
		} else {
			// Promote the point to a higher level
			level++
			if level > node.level {
				newRoot := NewNode(point, level, t.base)
				newRoot.children = append(newRoot.children, *t.root)
				t.root = &newRoot
				return
			}
		}
	}
}

// Remove removes a point (embedding vector) from the cover tree.
func (t *Tree[T]) Remove(point *Point) bool {
	if t.root == nil {
		return false
	}
	removed, newRoot := t.remove(t.root, point)
	t.root = newRoot
	if removed {
		var empty T
		t.values.data[point.index] = empty // Remove the value from the slice
		delete(t.indexMap, point.index)    // Remove the point from the map
	}
	return removed
}

func (t *Tree[T]) remove(node *Node, point *Point) (bool, *Node) {
	if node == nil {
		return false, nil
	}
	if t.distanceFnunc(point, node.point) == 0 {
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
func (t *Tree[T]) KNearestNeighbors(point *Point, k int) []*Neighbor {
	if t.root == nil {
		return nil
	}
	h := &Neighbors{}
	heap.Init(h)
	t.kNearestNeighbors(t.root, point, k, h)
	result := make([]*Neighbor, h.Len())
	for i := len(result) - 1; i >= 0; i-- {
		n := heap.Pop(h).(Neighbor)
		result[i] = &n
	}
	return result
}

func (t *Tree[T]) kNearestNeighbors(node *Node, point *Point, k int, h *Neighbors) {
	dist := t.distanceFnunc(point, node.point)
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
func NewTree[T any](base float32, distanceFn DistanceFunction) *Tree[T] {
	return &Tree[T]{base: base, distanceFnunc: distanceFn.Function(), distanceFuncName: distanceFn, values: values[T]{data: make([]T, 0)}}
}
