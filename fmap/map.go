package fmap

import (
	"math"
)

// Numeric interface defines the numeric types that can be used as values in the map.
type Numeric interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

// INT_PHI is a constant used in the hash function to scramble the keys.
// It is derived from the golden ratio and helps in distributing keys uniformly.
const INT_PHI = 0x9E3779B9

// FREE_KEY represents the zero value of int64 and is used to denote an empty slot in the keys array.
const FREE_KEY int64 = 0

// phiMix applies a supplemental hash function to the given int64 key.
// It multiplies the key with INT_PHI and then applies an XOR with the shifted value to improve the distribution.
func phiMix(x int64) int64 {
	h := x * INT_PHI
	return h ^ (h >> 16)
}

// NumericMap is a high-performance hash map for int64 keys and numeric values.
// It uses open addressing with linear probing and a custom hash function for int64 keys.
// The map is generic over the value type T, which must satisfy the Numeric interface.
// This implementation is optimized for performance and low memory overhead, and is not safe for concurrent use.
type NumericMap[T Numeric] struct {
	keys       []int64 // Array of keys
	data       []T     // Array of values corresponding to keys
	fillFactor float64 // Fill factor for resizing the map
	threshold  int     // Resize threshold based on capacity and fill factor
	size       int     // Number of elements in the map

	mask int64 // Mask for calculating indices during probing

	hasFreeKey bool // Indicates if the map contains the FREE_KEY
	freeVal    T    // Value associated with the FREE_KEY
}

// nextPowerOf2 returns the next power of two greater than or equal to x.
// It is used to ensure that the capacity of the map is always a power of two.
func nextPowerOf2(x uint64) uint64 {
	if x == 0 {
		return 1
	}
	if x&(x-1) == 0 {
		return x
	}
	return 1 << (64 - bitsLeadingZeros(x))
}

// bitsLeadingZeros returns the number of leading zeros in x.
func bitsLeadingZeros(x uint64) int {
	var n int
	for x != 0 {
		x >>= 1
		n++
	}
	return 64 - n
}

// arraySize calculates the required capacity based on the expected number of elements and fill factor.
// It ensures that the capacity is a power of two and at least 2.
func arraySize(expected int, fillFactor float64) int {
	s := nextPowerOf2(uint64(math.Ceil(float64(expected) / fillFactor)))
	if s < 2 {
		s = 2
	}
	if s > math.MaxInt64 {
		panic("Requested capacity exceeds maximum int64 value")
	}
	return int(s)
}

// Get retrieves the value associated with the given key.
// It returns the value and a boolean indicating whether the key was found.
func (m *NumericMap[T]) Get(key int64) (T, bool) {
	var zero T
	if key == FREE_KEY {
		if m.hasFreeKey {
			return m.freeVal, true
		}
		return zero, false
	}

	ptr := phiMix(key) & m.mask
	k := m.keys[ptr]

	if k == FREE_KEY {
		return zero, false
	}
	if k == key {
		return m.data[ptr], true
	}

	for {
		ptr = (ptr + 1) & m.mask
		k = m.keys[ptr]
		if k == FREE_KEY {
			return zero, false
		}
		if k == key {
			return m.data[ptr], true
		}
	}
}

// Put adds or updates the key with the value val.
func (m *NumericMap[T]) Put(key int64, val T) {
	if key == FREE_KEY {
		if !m.hasFreeKey {
			m.size++
		}
		m.hasFreeKey = true
		m.freeVal = val
		return
	}

	ptr := phiMix(key) & m.mask
	k := m.keys[ptr]

	if k == FREE_KEY { // Empty slot found
		m.keys[ptr] = key
		m.data[ptr] = val
		m.size++
		if m.size >= m.threshold {
			m.rehash()
		}
		return
	} else if k == key { // Key already exists, update value
		m.data[ptr] = val
		return
	}

	// Collision resolution via linear probing
	for {
		ptr = (ptr + 1) & m.mask
		k = m.keys[ptr]
		if k == FREE_KEY {
			m.keys[ptr] = key
			m.data[ptr] = val
			m.size++
			if m.size >= m.threshold {
				m.rehash()
			}
			return
		} else if k == key {
			m.data[ptr] = val
			return
		}
	}
}

// rehash resizes the map when the load factor exceeds the threshold.
// It doubles the capacity and reinserts all existing keys and values.
func (m *NumericMap[T]) rehash() {
	newCapacity := len(m.keys) * 2

	// Update mask and threshold based on new capacity
	m.mask = int64(newCapacity - 1)
	m.threshold = int(math.Floor(float64(newCapacity) * m.fillFactor))

	// Save old data
	oldKeys := m.keys
	oldData := m.data

	// Create new slices with updated capacity
	m.keys = make([]int64, newCapacity)
	m.data = make([]T, newCapacity)

	// Reset size and re-insert keys
	m.size = 0
	if m.hasFreeKey {
		m.size = 1
	}

	for i := 0; i < len(oldKeys); i++ {
		k := oldKeys[i]
		if k != FREE_KEY {

			m.Put(k, oldData[i])
		}
	}
}

// Size returns the number of elements in the map.
func (m *NumericMap[T]) Size() int {
	if m == nil {
		return 0
	}
	return m.size
}

// NewNumericMap creates a new NumericMap with the specified expected size and fill factor.
// The fill factor must be between 0 and 1 (exclusive), and determines when the map will be resized.
// The map will grow automatically as needed.
func NewNumericMap[T Numeric](expectedSize int, fillFactor float64) *NumericMap[T] {
	if fillFactor <= 0 || fillFactor >= 1 {
		panic("FillFactor must be in (0, 1)")
	}
	if expectedSize <= 0 {
		panic("Size must be positive")
	}

	capacity := arraySize(expectedSize, fillFactor)
	m := &NumericMap[T]{
		keys:       make([]int64, capacity),
		data:       make([]T, capacity),
		fillFactor: fillFactor,
		threshold:  int(math.Floor(float64(capacity) * fillFactor)),
		mask:       int64(capacity - 1),
	}
	return m
}