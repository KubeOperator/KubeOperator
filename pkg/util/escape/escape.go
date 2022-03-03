package escape

func Clean(str []byte) {
	for i := 0; i < len(str); i++ {
		str[i] = 0
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
