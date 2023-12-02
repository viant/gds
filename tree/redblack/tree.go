package redblack

import (
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"sync"
)

// Color defines color
type Color bool

const (
	Red   Color = false
	Black Color = true
)

// Ordered is a custom interface for types that are either integer, float, or string.
// These types must implement the Get() method.
type Ordered[T constraints.Ordered] interface {
	Get() T
}

// Callback is a function type that takes an Ordered type as an argument.
type Callback[T constraints.Ordered, O Ordered[T]] func(*O)

type Node[T constraints.Ordered, O Ordered[T]] struct {
	Color               Color
	Value               O
	Left, Right, Parent *Node[T, O]
}

type Tree[T constraints.Ordered, O Ordered[T]] struct {
	mux  sync.RWMutex
	size int
	Root *Node[T, O]
}

func NewNode[T constraints.Ordered, O Ordered[T]](value O) *Node[T, O] {
	return &Node[T, O]{Value: value, Color: Red}
}

func (t *Tree[T, O]) rotateLeft(x *Node[T, O]) {

	y := x.Right
	x.Right = y.Left

	if y.Left != nil {
		y.Left.Parent = x
	}

	y.Parent = x.Parent

	if x.Parent == nil {
		t.Root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}

	y.Left = x
	x.Parent = y
}

func (t *Tree[T, O]) rotateRight(y *Node[T, O]) {
	x := y.Left
	y.Left = x.Right

	if x.Right != nil {
		x.Right.Parent = y
	}

	x.Parent = y.Parent

	if y.Parent == nil {
		t.Root = x
	} else if y == y.Parent.Left {
		y.Parent.Left = x
	} else {
		y.Parent.Right = x
	}

	x.Right = y
	y.Parent = x
}

func (t *Tree[T, O]) insertFixup(z *Node[T, O]) {
	for z.Parent != nil && z.Parent.Color == Red {
		if z.Parent == z.Parent.Parent.Left {
			y := z.Parent.Parent.Right
			if y != nil && y.Color == Red {
				z.Parent.Color = Black
				y.Color = Black
				z.Parent.Parent.Color = Red
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Right {
					z = z.Parent
					t.rotateLeft(z)
				}
				z.Parent.Color = Black
				z.Parent.Parent.Color = Red
				t.rotateRight(z.Parent.Parent)
			}
		} else {
			y := z.Parent.Parent.Left
			if y != nil && y.Color == Red {
				z.Parent.Color = Black
				y.Color = Black
				z.Parent.Parent.Color = Red
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Left {
					z = z.Parent
					t.rotateRight(z)
				}
				z.Parent.Color = Black
				z.Parent.Parent.Color = Red
				t.rotateLeft(z.Parent.Parent)
			}
		}
	}
	t.Root.Color = Black
}

func (t *Tree[T, O]) insertNode(node *Node[T, O]) {
	var parent *Node[T, O]
	current := t.Root

	// Find the correct position for the new node
	for current != nil {
		parent = current
		if node.Value.Get() < current.Value.Get() {
			current = current.Left
		} else if node.Value.Get() > current.Value.Get() {
			current = current.Right
		} else {
			// Value already exists, do not insert
			return
		}
	}

	// Set the parent of the new node
	node.Parent = parent

	// Insert the new node into the tree
	if parent == nil {
		t.Root = node // Tree was empty
	} else if node.Value.Get() < parent.Value.Get() {
		parent.Left = node
	} else {
		parent.Right = node
	}
	t.size++
}

func (t *Tree[T, O]) Iterate(callback Callback[T, O]) {
	t.iterateInOrder(t.Root, callback)
}

func (t *Tree[T, O]) iterateInOrder(node *Node[T, O], callback Callback[T, O]) {
	if node != nil {
		t.iterateInOrder(node.Left, callback)  // Visit left subtree
		callback(&node.Value)                  // Process current node
		t.iterateInOrder(node.Right, callback) // Visit right subtree
	}
}

func (t *Tree[T, O]) Insert(value O) {
	t.mux.Lock()
	t.mux.Unlock()
	newNode := NewNode[T, O](value) // Create a new node with the provided value
	t.insertNode(newNode)
	t.insertFixup(newNode)
}

func (t *Tree[T, O]) Size() int {
	return t.size
}

func (t *Tree[T, O]) Delete(value T) {
	t.mux.Lock()
	defer t.mux.Unlock()

	nodeToDelete := t.findNode(t.Root, value)
	if nodeToDelete == nil {
		return // Value not found in the tree
	}
	t.size--
	t.deleteNode(nodeToDelete)
	t.fixAfterDeletion(nodeToDelete)
}

func (t *Tree[T, O]) fixAfterDeletion(x *Node[T, O]) {
	for x != t.Root && x.Color == Black {
		if x == x.Parent.Left {
			w := x.Parent.Right
			if w.Color == Red {
				w.Color = Black
				x.Parent.Color = Red
				t.rotateLeft(x.Parent)
				w = x.Parent.Right
			}

			if w.Left.Color == Black && w.Right.Color == Black {
				w.Color = Red
				x = x.Parent
			} else {
				if w.Right.Color == Black {
					w.Left.Color = Black
					w.Color = Red
					t.rotateRight(w)
					w = x.Parent.Right
				}
				w.Color = x.Parent.Color
				x.Parent.Color = Black
				w.Right.Color = Black
				t.rotateLeft(x.Parent)
				x = t.Root
			}
		} else { // Symmetrical case: x is a right child
			w := x.Parent.Left
			if w.Color == Red {
				w.Color = Black
				x.Parent.Color = Red
				t.rotateRight(x.Parent)
				w = x.Parent.Left
			}

			if w.Right.Color == Black && w.Left.Color == Black {
				w.Color = Red
				x = x.Parent
			} else {
				if w.Left.Color == Black {
					w.Right.Color = Black
					w.Color = Red
					t.rotateLeft(w)
					w = x.Parent.Left
				}
				w.Color = x.Parent.Color
				x.Parent.Color = Black
				w.Left.Color = Black
				t.rotateRight(x.Parent)
				x = t.Root
			}
		}
	}
	x.Color = Black
}

func (t *Tree[T, O]) deleteNode(node *Node[T, O]) {
	if node.Left == nil && node.Right == nil {
		// Node is a leaf
		t.replaceNode(node, nil)
	} else if node.Left == nil || node.Right == nil {
		// Node has only one child
		var child *Node[T, O]
		if node.Left != nil {
			child = node.Left
		} else {
			child = node.Right
		}
		t.replaceNode(node, child)
	} else {
		// Node has two children
		successor := t.minimum(node.Right)
		node.Value = successor.Value
		t.deleteNode(successor)
	}
}

func (t *Tree[T, O]) replaceNode(oldNode, newNode *Node[T, O]) {
	if oldNode.Parent == nil {
		t.Root = newNode
	} else if oldNode == oldNode.Parent.Left {
		oldNode.Parent.Left = newNode
	} else {
		oldNode.Parent.Right = newNode
	}
	if newNode != nil {
		newNode.Parent = oldNode.Parent
	}
}

func (t *Tree[T, O]) minimum(node *Node[T, O]) *Node[T, O] {
	for node.Left != nil {
		node = node.Left
	}
	return node
}
func (t *Tree[T, O]) Search(value T) *Node[T, O] {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.findNode(t.Root, value)
}

func (t *Tree[T, O]) InOrderTraversal() []O {
	var result = make([]O, 0, t.size)
	t.Iterate(func(o *O) {
		result = append(result, *o)
	})
	return result
}

func (t *Tree[T, O]) validateProperties() error {
	if t.Root == nil {
		return nil // An empty tree is a valid Red-Black Tree
	}

	// Check if root is black
	if t.Root.Color != Black {
		return errors.New("root is not black")
	}

	// Check the Red-Black Tree properties starting from the root
	_, err := t.validateNodeProperties(t.Root, 0)
	return err
}

// validateNodeProperties checks the properties for each node recursively
func (t *Tree[T, O]) validateNodeProperties(node *Node[T, O], blackCount int) (int, error) {
	if node == nil {
		// Return the count of black nodes for this path
		return blackCount, nil
	}

	if node.Color == Black {
		blackCount++
	} else if node.Color == Red {
		// Check for consecutive red nodes
		if (node.Left != nil && node.Left.Color == Red) || (node.Right != nil && node.Right.Color == Red) {
			return 0, errors.New("red node with red child found")
		}
	}

	// Recursively check the left and right subtrees
	leftBlackCount, err := t.validateNodeProperties(node.Left, blackCount)
	if err != nil {
		return 0, err
	}

	rightBlackCount, err := t.validateNodeProperties(node.Right, blackCount)
	if err != nil {
		return 0, err
	}

	delta := leftBlackCount - rightBlackCount
	if delta < 0 {
		delta *= -1
	}
	// Ensure the number of black nodes is consistent in both subtrees
	if delta > 1 {
		return 0, fmt.Errorf("inconsistent black node count, left: %d, right: %d", leftBlackCount, rightBlackCount)
	}

	return leftBlackCount, nil // You could also return rightBlackCount as they are the same
}

func (t *Tree[T, O]) findNode(node *Node[T, O], value T) *Node[T, O] {
	if node == nil {
		return nil
	}
	if value < node.Value.Get() {
		return t.findNode(node.Left, value)
	} else if value > node.Value.Get() {
		return t.findNode(node.Right, value)
	} else {
		return node
	}
}

// NewTree creates a new instance of a Tree
func NewTree[T constraints.Ordered, O Ordered[T]]() *Tree[T, O] {
	return &Tree[T, O]{}
}
