package cover

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValues_Encode(t *testing.T) {

	type Foo struct {
		ID     int
		Amount float64
	}

	var testCases = []struct {
		name   string
		data   interface{}
		encode func(v interface{}) ([]byte, error)
		decode func(data []byte) (interface{}, error)
	}{
		{
			name: "int",
			data: []int{1, 2, 1011},
			encode: func(data interface{}) ([]byte, error) {
				values := values[int]{data: data.([]int)}
				buffer := new(bytes.Buffer)
				err := values.Encode(buffer)
				return buffer.Bytes(), err
			},
			decode: func(data []byte) (interface{}, error) {
				values := values[int]{data: make([]int, 0)}
				err := values.Decode(bytes.NewReader(data))
				return values.data, err
			},
		},
		{
			name: "floats",
			data: []float64{1.1, 2.4, 1011.2},
			encode: func(data interface{}) ([]byte, error) {
				values := values[float64]{data: data.([]float64)}
				buffer := new(bytes.Buffer)
				err := values.Encode(buffer)
				return buffer.Bytes(), err
			},
			decode: func(data []byte) (interface{}, error) {
				values := values[float64]{data: make([]float64, 0)}
				err := values.Decode(bytes.NewReader(data))
				return values.data, err
			},
		},
		{
			name: "strings",
			data: []string{"abc", "def", "ghi"},
			encode: func(data interface{}) ([]byte, error) {
				values := values[string]{data: data.([]string)}
				buffer := new(bytes.Buffer)
				err := values.Encode(buffer)
				return buffer.Bytes(), err
			},
			decode: func(data []byte) (interface{}, error) {
				values := values[string]{data: make([]string, 0)}
				err := values.Decode(bytes.NewReader(data))
				return values.data, err
			},
		},
		{
			name: "comparable",
			data: []Foo{{ID: 1, Amount: 1.1}, {ID: 2, Amount: 2.2}},
			encode: func(data interface{}) ([]byte, error) {
				values := values[Foo]{data: data.([]Foo)}
				buffer := new(bytes.Buffer)
				err := values.Encode(buffer)
				return buffer.Bytes(), err
			},
			decode: func(data []byte) (interface{}, error) {
				values := values[Foo]{data: make([]Foo, 0)}
				err := values.Decode(bytes.NewReader(data))
				return values.data, err
			},
		},
	}

	for _, testCase := range testCases {
		data, err := testCase.encode(testCase.data)
		if err != nil {
			t.Errorf("failed to encode %v: %v", testCase.name, err)
		}
		if len(data) == 0 {
			t.Errorf("expected data for %v, but got empty", testCase.name)
		}
		actual, err := testCase.decode(data)
		if err != nil {
			t.Errorf("failed to decode %v: %v", testCase.name, err)
		}
		assert.EqualValues(t, testCase.data, actual, testCase.name)
	}
}
