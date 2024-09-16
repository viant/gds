package fmap

import (
	"strconv"
	"testing"
)

// TestCase defines a single test case with inputs and expected outputs.
type TestCase[T any] struct {
	name        string
	key         int64
	value       T
	updateVal   T
	expectGet   T
	expectSize  int
	expectFound bool
}

// TestFastMapDataDriven tests the FastMap using a data-driven approach.
func TestFastMapDataDriven(t *testing.T) {
	// Define test cases
	testCases := []TestCase[int64]{
		{
			name:        "InsertKey1",
			key:         1,
			value:       100,
			expectGet:   100,
			expectSize:  1,
			expectFound: true,
		},
		{
			name:        "InsertKey2",
			key:         2,
			value:       200,
			expectGet:   200,
			expectSize:  2,
			expectFound: true,
		},
		{
			name:        "UpdateKey1",
			key:         1,
			value:       100,
			updateVal:   150,
			expectGet:   150,
			expectSize:  2,
			expectFound: true,
		},
		{
			name:        "NonExistentKey",
			key:         3,
			expectGet:   0, // zero value for int64
			expectSize:  2,
			expectFound: false,
		},
		{
			name:        "InsertFreeKey",
			key:         0,
			value:       500,
			expectGet:   500,
			expectSize:  3,
			expectFound: true,
		},
	}

	// Initialize a new FastMap
	m := NewFastMap[int64](4, 0.75)

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Insert or update key
			if tc.value != 0 {
				m.Put(tc.key, tc.value)
			}
			if tc.updateVal != 0 {
				m.Put(tc.key, tc.updateVal)
			}

			// Get value
			val, found := m.Get(tc.key)
			if found != tc.expectFound {
				t.Errorf("Test %s: Expected found=%v, got %v", tc.name, tc.expectFound, found)
			}
			if val != tc.expectGet {
				t.Errorf("Test %s: Expected value=%v, got %v", tc.name, tc.expectGet, val)
			}

			// Check size
			if m.Size() != tc.expectSize {
				t.Errorf("Test %s: Expected size=%d, got %d", tc.name, tc.expectSize, m.Size())
			}
		})
	}

	next := m.Iterator()
	for i := 0; i < m.Size(); i++ {
		k, v := next()
		t.Logf("Key: %d, Value: %d", k, v)
	}
}

// TestFastMapCollisionDataDriven tests collision handling in FastMap using data-driven approach.
func TestFastMapCollisionDataDriven(t *testing.T) {
	// Define test cases with keys that may cause collisions
	testCases := []TestCase[int]{
		{
			name:        "InsertKey1",
			key:         1,
			value:       100,
			expectGet:   100,
			expectSize:  1,
			expectFound: true,
		},
		{
			name:        "InsertKey17",
			key:         17,
			value:       200,
			expectGet:   200,
			expectSize:  2,
			expectFound: true,
		},
		{
			name:        "InsertKey33",
			key:         33,
			value:       300,
			expectGet:   300,
			expectSize:  3,
			expectFound: true,
		},
	}

	m := NewFastMap[int](4, 0.75)

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Insert key
			m.Put(tc.key, tc.value)

			// Get value
			val, found := m.Get(tc.key)
			if found != tc.expectFound {
				t.Errorf("Test %s: Expected found=%v, got %v", tc.name, tc.expectFound, found)
			}
			if val != tc.expectGet {
				t.Errorf("Test %s: Expected value=%v, got %v", tc.name, tc.expectGet, val)
			}

			// Check size
			if m.Size() != tc.expectSize {
				t.Errorf("Test %s: Expected size=%d, got %d", tc.name, tc.expectSize, m.Size())
			}
		})
	}
}

// TestFastMapRehashDataDriven tests rehashing using a data-driven approach.
func TestFastMapRehashDataDriven(t *testing.T) {
	// Define test cases
	var testCases []TestCase[int]
	for i := 1; i <= 20; i++ {
		tc := TestCase[int]{
			name:        "InsertKey" + strconv.Itoa(i),
			key:         int64(i),
			value:       i * 10,
			expectGet:   i * 10,
			expectSize:  i,
			expectFound: true,
		}
		testCases = append(testCases, tc)
	}

	m := NewFastMap[int](4, 0.75)

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Insert key
			m.Put(tc.key, tc.value)

			// Get value
			val, found := m.Get(tc.key)
			if found != tc.expectFound {
				t.Errorf("Test %s: Expected found=%v, got %v", tc.name, tc.expectFound, found)
			}
			if val != tc.expectGet {
				t.Errorf("Test %s: Expected value=%v, got %v", tc.name, tc.expectGet, val)
			}

			// Check size
			if m.Size() != tc.expectSize {
				t.Errorf("Test %s: Expected size=%d, got %d", tc.name, tc.expectSize, m.Size())
			}
		})
	}
}

// TestFastMapDifferentTypesDataDriven tests different Fast types using data-driven approach.
func TestFastMapDifferentTypesDataDriven(t *testing.T) {
	// Test with float64 values
	floatTestCases := []TestCase[float64]{
		{
			name:        "InsertFloatKey1",
			key:         1,
			value:       1.5,
			expectGet:   1.5,
			expectSize:  1,
			expectFound: true,
		},
		{
			name:        "InsertFloatKey2",
			key:         2,
			value:       2.5,
			expectGet:   2.5,
			expectSize:  2,
			expectFound: true,
		},
	}

	floatMap := NewFastMap[float64](4, 0.75)

	for _, tc := range floatTestCases {
		t.Run(tc.name, func(t *testing.T) {
			floatMap.Put(tc.key, tc.value)
			val, found := floatMap.Get(tc.key)
			if found != tc.expectFound {
				t.Errorf("Test %s: Expected found=%v, got %v", tc.name, tc.expectFound, found)
			}
			if val != tc.expectGet {
				t.Errorf("Test %s: Expected value=%v, got %v", tc.name, tc.expectGet, val)
			}
			if floatMap.Size() != tc.expectSize {
				t.Errorf("Test %s: Expected size=%d, got %d", tc.name, tc.expectSize, floatMap.Size())
			}
		})
	}

	// Test with uint64 values
	uintTestCases := []TestCase[uint64]{
		{
			name:        "InsertUintKey1",
			key:         1,
			value:       100,
			expectGet:   100,
			expectSize:  1,
			expectFound: true,
		},
		{
			name:        "InsertUintKey2",
			key:         2,
			value:       200,
			expectGet:   200,
			expectSize:  2,
			expectFound: true,
		},
	}

	uintMap := NewFastMap[uint64](4, 0.75)

	for _, tc := range uintTestCases {
		t.Run(tc.name, func(t *testing.T) {
			uintMap.Put(tc.key, tc.value)
			val, found := uintMap.Get(tc.key)
			if found != tc.expectFound {
				t.Errorf("Test %s: Expected found=%v, got %v", tc.name, tc.expectFound, found)
			}
			if val != tc.expectGet {
				t.Errorf("Test %s: Expected value=%v, got %v", tc.name, tc.expectGet, val)
			}
			if uintMap.Size() != tc.expectSize {
				t.Errorf("Test %s: Expected size=%d, got %d", tc.name, tc.expectSize, uintMap.Size())
			}
		})
	}
}
