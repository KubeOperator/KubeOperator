package escape

import "unsafe"

func Clean(str string) {
	clearStr := *(*[]byte)(unsafe.Pointer(&str))
	for i := 0; i < len(clearStr); i++ {
		clearStr[i] = 0
	}
}

func GetByte(tmp interface{}) []byte {
	k, ok := tmp.(string)
	if !ok {
		return []byte{}
	} else {
		return []byte(k)
	}
}
