package cover

// Neighbor represents a neighbor of a point.
type Neighbor struct {
	Point    *Point
	Distance float32
}

// Neighbors is a slice nearest of Neighbors.
type Neighbors []Neighbor

// Len Implement the heap.Interface for Neighbors.
func (h Neighbors) Len() int { return len(h) }

// Less Implement the heap.Interface for Neighbors.
func (h Neighbors) Less(i, j int) bool { return h[i].Distance > h[j].Distance }

// Swap Implement the heap.Interface for Neighbors.
func (h Neighbors) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Push Implement the heap.Interface for Neighbors.
func (h *Neighbors) Push(x interface{}) {
	*h = append(*h, x.(Neighbor))
}

// Pop Implement the heap.Interface for Neighbors.
func (h *Neighbors) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
