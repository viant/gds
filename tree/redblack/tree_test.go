package redblack

import (
	"reflect"
	"testing"
)

func TestRedBlackTree_Delete(t *testing.T) {

	cases := []struct {
		name      string
		elements  []Integer
		deleteVal int
		want      []Integer // Expected inorder traversal after deletion
	}{
		{
			name:      "Delete leaf",
			elements:  []Integer{10, 5, 15, 3, 7, 12, 18},
			deleteVal: 3,
			want:      []Integer{5, 7, 10, 12, 15, 18},
		},
		{
			name:      "Delete node with one child",
			elements:  []Integer{10, 5, 15, 3, 7, 12, 18},
			deleteVal: 15,
			want:      []Integer{3, 5, 7, 10, 12, 18},
		},
		{
			name:      "Delete node with two children",
			elements:  []Integer{10, 5, 15, 3, 7, 12, 18},
			deleteVal: 10,
			want:      []Integer{3, 5, 7, 12, 15, 18},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tree := NewTree[int, Integer]()
			for _, v := range tc.elements {
				tree.Insert(v)
			}

			tree.Delete(tc.deleteVal)

			got := tree.InOrderTraversal()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("After deleting %v, got inorder traversal %v, want %v", tc.deleteVal, got, tc.want)
			}
			if err := tree.validateProperties(); err != nil {
				t.Errorf("Red-Black Tree properties violated after insertion: %v", err)
			}
		})
	}
}

func TestRedBlackTree_Insert(t *testing.T) {
	cases := []struct {
		name     string
		elements []Integer // Elements to be inserted
		want     []Integer // Expected inorder traversal after insertions
	}{
		{
			name:     "Insert single element",
			elements: []Integer{10},
			want:     []Integer{10},
		},
		{
			name:     "Insert multiple elements",
			elements: []Integer{10, 5, 15, 3, 7, 12, 18},
			want:     []Integer{3, 5, 7, 10, 12, 15, 18},
		},
		{
			name:     "Insert duplicate element",
			elements: []Integer{10, 5, 15, 3, 7, 12, 18, 10},
			want:     []Integer{3, 5, 7, 10, 12, 15, 18}, // Assuming duplicates are not allowed
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tree := NewTree[int, Integer]()
			for _, v := range tc.elements {
				tree.Insert(v)
			}

			got := tree.InOrderTraversal()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("After insertion, got inorder traversal %v, want %v", got, tc.want)
			}

			if err := tree.validateProperties(); err != nil {
				t.Errorf("Red-Black Tree properties violated after insertion: %v", err)
			}
		})
	}
}
