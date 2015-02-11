package msgpack

import (
	"unsafe"
)

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

func encodeInt(buf []byte, offset, val int) int {
	if val < 128 && val >= 0 {
		buf[offset] = byte(0x7f & val)
		offset += 1
	} else if val < 0 && val > -32 {
		buf[offset] = byte(0xe0 | (0x1f & -val))
		offset += 1
	} else { // go to int64
		// find the smallest mask we can use
		switch val {
		case val & 0xff: // maybe 8 bit
			offset = encodeInt8(buf, offset, int64(val))
		case val & 0xffff: // maybe 16 bit
			offset = encodeInt16(buf, offset, int64(val))
		case val & 0xffffffff: // maybe 32 bit
			offset = encodeInt32(buf, offset, int64(val))
		default:
			offset = encodeInt64(buf, offset, int64(val))
		}
	}
	return offset
}

func encodeInt8(buf []byte, offset int, val int64) int {
	if val > 0x7f {
		return encodeInt16(buf, offset, val)
	}
	buf[offset] = byte(0xd0)
	buf[offset+1] = byte(val)
	return offset + 2
}

func encodeInt16(buf []byte, offset int, val int64) int {
	if val > 0x7fff {
		return encodeInt32(buf, offset, val)
	}
	buf[offset] = byte(0xd1)
	buf[offset+1] = byte(val >> 8)
	buf[offset+2] = byte(val & 0xff)
	return offset + 3
}

func encodeInt32(buf []byte, offset int, val int64) int {
	if val > 0x7fffffff {
		return encodeInt64(buf, offset, val)
	}
	buf[offset] = byte(0xd2)
	buf[offset+1] = byte(val >> 24)
	buf[offset+2] = byte(val >> 16)
	buf[offset+3] = byte(val >> 8)
	buf[offset+4] = byte(val & 0xff)
	return offset + 5
}

func encodeInt64(buf []byte, offset int, val int64) int {
	buf[offset] = byte(0xd3)
	buf[offset+1] = byte(val >> 56)
	buf[offset+2] = byte(val >> 48)
	buf[offset+3] = byte(val >> 40)
	buf[offset+4] = byte(val >> 32)
	buf[offset+5] = byte(val >> 24)
	buf[offset+6] = byte(val >> 16)
	buf[offset+7] = byte(val >> 8)
	buf[offset+8] = byte(val & 0xff)
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

func Encode(input interface{}) []byte {
	// try to just do 1 allocation.
	//TODO: try to predict a slightly larger size. what is the rate of growth on msgpack encoding?
	ret := make([]byte, uint64(unsafe.Sizeof(input)))
	offset := 0
	switch input.(type) {
	case int:
		offset = encodeInt(ret, offset, input.(int))
	case uint:
		offset = encodeUint(ret, offset, input.(uint))
	case int64:
		offset = encodeInt64(ret, offset, input.(int64))
	case uint64:
		offset = encodeUint(ret, offset, input.(uint))
	case string:
	case map[string]interface{}:
	case []interface{}:
	case bool:
		offset = encodeBool(ret, offset, input.(bool))
	case nil:
		offset = encodeNil(ret, offset)
	default:
	}
	return ret[:offset]
}
