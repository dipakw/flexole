package util

func PackUint16(n uint16) []byte {
	buf := make([]byte, 2)
	buf[0] = byte(n >> 8)
	buf[1] = byte(n)
	return buf
}

func UnpackUint16(buf []byte) uint16 {
	return uint16(buf[0])<<8 | uint16(buf[1])
}
