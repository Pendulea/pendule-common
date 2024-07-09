package pcommon

import (
	"bytes"
	"encoding/binary"
	"log"
	"reflect"
	"unsafe"
)

type superBytes struct{}

var Bytes = superBytes{}

func (sb superBytes) BytesCount(slice interface{}) int64 {
	// Use reflect to get the underlying value and type of the slice
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("SliceSizeInBytes: provided argument is not a slice")
	}

	// Get the size of the slice header
	sliceHeaderSize := int64(unsafe.Sizeof(reflect.SliceHeader{}))

	// Get the size of one element in the slice
	if v.Len() == 0 {
		return sliceHeaderSize
	}

	elemSize := int64(v.Type().Elem().Size())

	// Calculate the total size: header size + element size * length of slice
	totalSize := sliceHeaderSize + (elemSize * int64(v.Len()))
	return totalSize
}

// int64ToBytes converts an int64 to a byte slice.
func (sb superBytes) Int64ToBytes(n int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, n)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

// bytesToInt64 converts a byte slice to an int64.
func (sb superBytes) BytesToInt64(b []byte) int64 {
	buf := bytes.NewBuffer(b)
	var n int64
	err := binary.Read(buf, binary.BigEndian, &n)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

// bytesToFloat64 converts a byte slice to a float64.
func (sb superBytes) BytesToFloat64(b []byte) float64 {
	buf := bytes.NewBuffer(b)
	var n float64
	err := binary.Read(buf, binary.BigEndian, &n)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

// float64ToBytes converts a float64 to a byte slice. (length: 8 bytes)
func (sb superBytes) Float64ToBytes(n float64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, n)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}
