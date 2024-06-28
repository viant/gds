package cover

import (
	"fmt"
	"github.com/viant/bintly"
	"io"
	"reflect"
	"sync"
	"unsafe"
)

var writers = bintly.NewWriters()
var readers = bintly.NewReaders()

type values[T any] struct {
	Type  reflect.Type
	vType interface{}
	data  []T
	sync.RWMutex
}

func (v *values[T]) useType(Type reflect.Type) {
	v.Type = Type
	v.vType = reflect.New(Type).Elem().Interface()
}

func (v *values[T]) put(value T) int32 {
	v.Lock()
	defer v.Unlock()
	if v.vType == nil {
		v.useType(reflect.TypeOf(value))
	}
	ret := len(v.data)
	v.data = append(v.data, value)
	return int32(ret)
}

func (v *values[T]) Decode(reader io.Reader) error {
	buffer := readers.Get()
	defer readers.Put(buffer)
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if err = buffer.FromBytes(data); err != nil {
		return err
	}
	v.data = make([]T, 0)
	values := any(v.data)
	switch actual := values.(type) {
	case []string:
		buffer.Strings(&actual)
	case []int:
		buffer.Ints(&actual)
	case []float32:
		buffer.Float32s(&actual)
	case []float64:
		buffer.Float64s(&actual)
	case []bool:
		buffer.Bools(&actual)
	default:
		if v.Type.Comparable() {
			raw := unsafe.Slice((*byte)(unsafe.Pointer(&values)), v.Type.Size())
			buffer.Uint8s(&raw)
		} else {
			return v.decodeCustom(buffer)
		}
	}
	return err
}

func (v *values[T]) Encode(writer io.Writer) error {
	buffer := writers.Get()
	defer writers.Put(buffer)
	values := any(v.data)
	switch actual := values.(type) {
	case []string:
		buffer.Strings(actual)
	case []int:
		buffer.Ints(actual)
	case []float32:
		buffer.Float32s(actual)
	case []float64:
		buffer.Float64s(actual)
	case []bool:
		buffer.Bools(actual)
	default:
		if v.Type.Comparable() {
			raw := unsafe.Slice((*byte)(unsafe.Pointer(&values)), v.Type.Size())
			buffer.Uint8s(raw)
		} else {
			return v.encodeCustom(buffer)
		}
	}
	_, err := writer.Write(buffer.Bytes())
	return err
}

func (v *values[T]) encodeCustom(writer *bintly.Writer) error {
	writer.Alloc(int32(len(v.data)))
	for i := range v.data {
		encoder, ok := any(v.data[i]).(bintly.Encoder)
		if !ok {
			return fmt.Errorf("unable to cast Encoder from %T", v.data[i])
		}
		if err := encoder.EncodeBinary(writer); err != nil {
			return err
		}
	}
	return nil
}

func (v *values[T]) value(index int32) T {
	v.RLock()
	defer v.RUnlock()
	return v.data[index]
}

func (v *values[T]) decodeCustom(buffer *bintly.Reader) error {
	size := buffer.Alloc()
	v.data = make([]T, size)
	for i := range v.data {
		decoder, ok := any(v.data[i]).(bintly.Decoder)
		if !ok {
			return fmt.Errorf("unable to cast Encoder from %T", v.data[i])
		}
		if err := decoder.DecodeBinary(buffer); err != nil {
			return err
		}
	}
	return nil
}
