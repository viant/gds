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
		v.data = any(actual).([]T)
	case []int:
		buffer.Ints(&actual)
		v.data = any(actual).([]T)
	case []float32:
		buffer.Float32s(&actual)
		v.data = any(actual).([]T)
	case []float64:
		buffer.Float64s(&actual)
		v.data = any(actual).([]T)
	case []bool:
		buffer.Bools(&actual)
		v.data = any(actual).([]T)
	default:
		v.ensureType()
		if v.Type.Comparable() && v.Type.Kind() != reflect.Pointer {
			var raw []byte
			buffer.Uint8s(&raw)
			size := len(raw) / int(v.Type.Size())
			v.data = make([]T, size)
			data := unsafe.Slice((*byte)(unsafe.Pointer(&v.data)), len(v.data)*int(v.Type.Size()))
			copy(data, raw)
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
		v.ensureType()
		if v.Type.Comparable() && v.Type.Kind() != reflect.Pointer {
			raw := unsafe.Slice((*byte)(unsafe.Pointer(&v.data)), len(v.data)*int(v.Type.Size()))
			buffer.Uint8s(raw)
		} else {
			if err := v.encodeCustom(buffer); err != nil {
				return err
			}
		}
	}
	_, err := writer.Write(buffer.Bytes())
	return err
}

func (v *values[T]) ensureType() {
	if v.Type == nil {
		v.useType(reflect.TypeOf(v.data).Elem())
	}
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

	if v.Type.Kind() == reflect.Ptr {
		for i := range v.data {
			v.data[i] = reflect.New(v.Type.Elem()).Interface().(T)
		}
	}

	for i := range v.data {
		decoder, ok := any(v.data[i]).(bintly.Decoder)
		if !ok {
			decoder, ok = any(&v.data[i]).(bintly.Decoder)
			if !ok {
				return fmt.Errorf("unable to cast Encoder from %T", v.data[i])
			}
		}
		if err := decoder.DecodeBinary(buffer); err != nil {
			return err
		}
	}
	return nil
}
