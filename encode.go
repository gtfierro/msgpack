// This MsgPack encoding/decoding library was born out of a partial dissatisfaction with
// the current implementation of (and lack of documentation of) MsgPack encoders/decoders
// for Go. I also wanted some practice writing optimized Go code that attempted to avoid
// allocations and "got out of the way" of libraries using it.
//
// This package is not finished yet, but it is close, and I've made an effort to do things
// in a standard way so that what is actually being done underneath the covers is easy
// to figure out. The focus of this package is not to provide the most fully featured
// MsgPack implementation out there, but rather support the main data types that I
// use in my work (most of the basic types, also []interface{} and map[string]interface{}).
// You'll notice that I don't use the `reflect` package at all, and that's because I
// wanted to see what code looked like when you just did a switch on types instead of
// using runtime reflection to figure out the types. This of course comes at the cost
// of not being able to extend this MsgPack implementation to arbitrary data types.
// Implementing a new type on your own is not difficult, though.
//
// Coming up next are some more convenience methods for doing decoding and encoding
// through writer/reader/buffer interfaces
package msgpack

import (
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v2"
	"math"
	"sync"
)

const DEFAULT_ARR_SIZE = 15

var arrpool = sync.Pool{
	New: func() interface{} {
		return make([]interface{}, DEFAULT_ARR_SIZE)
	},
}

func getNewArray(length int) []interface{} {
	if length <= 15 {
		return arrpool.Get().([]interface{})[:length]
	} else {
		return make([]interface{}, length)
	}
}

/** Each of these functions should take 3 arguments: the buffer to add into, an
 * offset into that buffer (where our writes start), and the value to encode.
 * After encoding the value and placing the byte sequence in the buffer
 * starting at the given offset, the function should return the next available
 * place to place bytes (the next offset to use)
 */

func encodeNil(buf []byte, offset int) int {
	buf[offset] = byte(0xc0)
	return offset + 1
}

func encodeBool(buf []byte, offset int, val bool) int {
	if val {
		buf[offset] = byte(0xc3)
	} else {
		buf[offset] = byte(0xc2)
	}
	return offset + 1
}

func encodeInt(buf []byte, offset int, val int64) int {
	if val < 128 && val >= 0 {
		buf[offset] = byte(0x7f & val)
		offset += 1
	} else if val < 0 && val > -32 {
		buf[offset] = byte(0xe0 | (0x1f & -val))
		offset += 1
	} else { // go to int64
		// find the smallest mask we can use
		switch {
		case val < -2147483648 || val >= 2147483648:
			buf[offset] = byte(0xd3)
			buf[offset+1] = byte(uint64(val) >> 56)
			buf[offset+2] = byte(uint64(val) >> 48)
			buf[offset+3] = byte(uint64(val) >> 40)
			buf[offset+4] = byte(uint64(val) >> 32)
			buf[offset+5] = byte(uint64(val) >> 24)
			buf[offset+6] = byte(uint64(val) >> 16)
			buf[offset+7] = byte(uint64(val) >> 8)
			buf[offset+8] = byte(uint64(val) & 0xff)
			offset += 9
		case val < -32768 || val >= 32768:
			buf[offset] = byte(0xd2)
			buf[offset+1] = byte(uint64(val) >> 24)
			buf[offset+2] = byte(uint64(val) >> 16)
			buf[offset+3] = byte(uint64(val) >> 8)
			buf[offset+4] = byte(uint64(val))
			offset += 5
		case val < -128 || val > 128:
			buf[offset] = byte(0xd1)
			buf[offset+1] = byte(uint64(val) >> 8)
			buf[offset+2] = byte(uint64(val))
			offset += 3
		case val < -32:
			buf[offset] = byte(0xd0)
			buf[offset+1] = byte(uint64(val))
			offset += 2
		default:
			println("default")
		}
	}
	return offset
}

func encodeFloat32(buf []byte, offset int, val float32) int {
	bits := uint64(math.Float32bits(val))
	buf[offset] = byte(0xca)
	buf[offset+1] = byte(bits >> 24)
	buf[offset+2] = byte(bits >> 16)
	buf[offset+3] = byte(bits >> 8)
	buf[offset+4] = byte(bits)
	return offset + 5
}

func encodeFloat64(buf []byte, offset int, val float64) int {
	bits := math.Float64bits(val)
	buf[offset] = byte(0xcb)
	buf[offset+1] = byte(bits >> 56)
	buf[offset+2] = byte(bits >> 48)
	buf[offset+3] = byte(bits >> 40)
	buf[offset+4] = byte(bits >> 32)
	buf[offset+5] = byte(bits >> 24)
	buf[offset+6] = byte(bits >> 16)
	buf[offset+7] = byte(bits >> 8)
	buf[offset+8] = byte(bits)
	return offset + 9
}

func encodeUint(buf []byte, offset int, val uint) int {
	switch val {
	case val & 0xff: // uint8
		buf[offset] = byte(0xcc)
		buf[offset+1] = byte(val)
		offset += 2
	case val & 0xffff: // uint16
		buf[offset] = byte(0xcd)
		buf[offset+1] = byte(val >> 8)
		buf[offset+2] = byte(val & 0xff)
		offset += 3
	case val & 0xffffffff: // uint32
		buf[offset] = byte(0xce)
		buf[offset+1] = byte(val >> 24)
		buf[offset+2] = byte(val >> 16)
		buf[offset+3] = byte(val >> 8)
		buf[offset+4] = byte(val & 0xff)
		offset += 5
	default: // uint64
		buf[offset] = byte(0xcf)
		buf[offset+1] = byte(val >> 56)
		buf[offset+2] = byte(val >> 48)
		buf[offset+3] = byte(val >> 40)
		buf[offset+4] = byte(val >> 32)
		buf[offset+5] = byte(val >> 24)
		buf[offset+6] = byte(val >> 16)
		buf[offset+7] = byte(val >> 8)
		buf[offset+8] = byte(val & 0xff)
		offset += 9
	}
	return offset
}

// Encodes @val as a bigendian unsigned integer in buffer @buf
// starting at offset @offset. Attempts to make it fit in @length
// bytes, and will truncate if it cannot
func encodeLength(buf []byte, offset int, val uint, length int) int {
	switch {
	case length == 1: // uint8
		buf[offset] = byte(val)
		offset += 1
	case length == 2: // uint16
		buf[offset] = byte(val >> 8)
		buf[offset+1] = byte(val & 0xff)
		offset += 2
	case length == 4: // uint32
		buf[offset] = byte(val >> 24)
		buf[offset+1] = byte(val >> 16)
		buf[offset+2] = byte(val >> 8)
		buf[offset+3] = byte(val & 0xff)
		offset += 4
	default: // uint64
		buf[offset] = byte(val >> 56)
		buf[offset+1] = byte(val >> 48)
		buf[offset+2] = byte(val >> 40)
		buf[offset+3] = byte(val >> 32)
		buf[offset+4] = byte(val >> 24)
		buf[offset+5] = byte(val >> 16)
		buf[offset+6] = byte(val >> 8)
		buf[offset+7] = byte(val & 0xff)
		offset += 8
	}
	return offset
}

func encodeString(buf []byte, offset int, val string) int {
	l := len(val)
	switch {
	case l <= 31: // fixstr
		buf[offset] = byte(0xa0 | l)
		offset += 1
	case l <= 255: // str8
		buf[offset] = byte(0xd9)
		buf[offset+1] = byte(l)
		offset += 2
	case l <= 65535: // str16
		buf[offset] = byte(0xda)
		offset += 1
		offset = encodeLength(buf, offset, uint(l), 2)
	default: // str32
		buf[offset] = byte(0xdb)
		offset += 1
		offset = encodeLength(buf, offset, uint(l), 4)
	}
	for i := 0; i < l; i++ { // TODO fewer copies, e.g. not 1 byte at a time
		buf[offset+i] = val[i]
	}
	offset += l
	return offset
}

func encodeArray(buf []byte, offset int, val []interface{}) int {
	l := len(val)
	switch {
	case l <= 15:
		buf[offset] = byte(0x90 | l)
		offset += 1
	case l <= 65535: // (2^16 - 1)
		buf[offset] = byte(0xdc)
		offset += 1
		offset = encodeLength(buf, offset, uint(l), 2)
	default: // up to 4294967295 (2^32 - 1)
		buf[offset] = byte(0xdd)
		offset += 1
		offset = encodeLength(buf, offset, uint(l), 4)
	}
	for i := 0; i < l; i++ {
		offset = doEncode(val[i], &buf, offset)
	}
	return offset
}

func encodeMap(buf []byte, offset int, val map[string]interface{}) int {
	l := len(val)
	switch {
	case l <= 15:
		buf[offset] = byte(0x80 | l)
		offset += 1
	case l <= 65535: // 2^16 - 1
		buf[offset] = byte(0xde)
		offset += 1
		offset = encodeLength(buf, offset, uint(l), 2)
	default: // up to 4294967295 (2^32 - 1)
		buf[offset] = byte(0xdf)
		offset += 1
		offset = encodeLength(buf, offset, uint(l), 4)
	}
	for k, v := range val {
		offset = doEncode(k, &buf, offset)
		offset = doEncode(v, &buf, offset)
	}
	return offset
}

func doEncode(input interface{}, ret *[]byte, offset int) int {
	switch input.(type) {
	case int:
		offset = encodeInt(*ret, offset, int64(input.(int)))
	case int8:
		offset = encodeInt(*ret, offset, int64(input.(int8)))
	case int16:
		offset = encodeInt(*ret, offset, int64(input.(int16)))
	case int32:
		offset = encodeInt(*ret, offset, int64(input.(int32)))
	case uint:
		offset = encodeUint(*ret, offset, input.(uint))
	case uint8:
		offset = encodeUint(*ret, offset, uint(input.(uint8)))
	case uint16:
		offset = encodeUint(*ret, offset, uint(input.(uint16)))
	case uint32:
		offset = encodeUint(*ret, offset, uint(input.(uint32)))
	case int64:
		offset = encodeInt(*ret, offset, input.(int64))
	case uint64:
		offset = encodeUint(*ret, offset, uint(input.(uint64)))
	case float32:
		offset = encodeFloat32(*ret, offset, input.(float32))
	case float64:
		offset = encodeFloat64(*ret, offset, input.(float64))
	case string:
		offset = encodeString(*ret, offset, input.(string))
	case map[string]interface{}:
		offset = encodeMap(*ret, offset, input.(map[string]interface{}))
	case []interface{}:
		offset = encodeArray(*ret, offset, input.([]interface{}))
	case bool:
		offset = encodeBool(*ret, offset, input.(bool))
	case nil:
		offset = encodeNil(*ret, offset)
	case interface{}:
		offset = doEncodeReflect(input, ret, offset)
	default:
		offset = doEncodeReflect(input, ret, offset)
	}
	return offset
}

func doEncodeReflect(input interface{}, ret *[]byte, offset int) int {
	b, err := msgpack.Marshal(input)
	if err != nil {
		fmt.Errorf("Error %v\n", err)
	}
	copy((*ret)[offset:], b)
	return offset + len(b)
}

// Encodes the input as a msgpack byte array, which is provided
// by the user. This allows the user to control how many allocations
// are done. Returns the length of the encoded message, but does
// not adjust the length of the input array.
func Encode(input interface{}, ret *[]byte) int {
	return doEncode(input, ret, 0)
}
