package redblack_test

import (
	"fmt"
	"github.com/viant/gds/tree/redblack"
)

func ExampleNewTree() {
	tree := redblack.NewTree[int, redblack.Integer]()
	tree.Insert(1)
	tree.Insert(101)
	tree.Insert(54)
	sorted := tree.InOrderTraversal()
	fmt.Printf("%v\n", sorted)
}
