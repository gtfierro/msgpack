package msgpack

func Encode(input interface{}) []byte {
	switch input.(type) {
	case int:
		if input.(int) < 128 {
			return []byte{byte(0x7f & input.(int))}
		} else {
			// go to int64
		}
	case uint:
	case int64:
	case uint64:
	case string:
	case map[string]interface{}:
	case []interface{}:
	case bool:
		if input.(bool) {
			return []byte{0xc3}
		} else {
			return []byte{0xc2}
		}
	case nil:
		return []byte{0xc0}
	default:
	}
	return []byte{}
}
