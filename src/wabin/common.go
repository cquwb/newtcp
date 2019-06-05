package wabin

//"fmt"

//大端转小端
func byteToUint32(b []byte) uint32 {
	return (uint32(b[0]) << 24) | (uint32(b[1]) << 16) | (uint32(b[2]))<<8 | uint32(b[3])
}

func byteToUint64(b []byte) uint64 {
	return (uint64(b[0]) << 56) | (uint64(b[1]) << 48) | (uint64(b[2]) << 40) | (uint64(b[3]) << 32) |
		(uint64(b[4]) << 24) | (uint64(b[5]) << 16) | (uint64(b[6]) << 8) | uint64(b[7])
}

//小端转大端
func Uint32ToByte(n uint32, b []byte) {
	b[3] = byte(n)
	b[2] = byte(n >> 8)
	b[1] = byte(n >> 16)
	b[0] = byte(n >> 24)
}

func Uint64ToByte(n uint64, b []byte) {
	b[7] = byte(n)
	b[6] = byte(n >> 8)
	b[5] = byte(n >> 16)
	b[4] = byte(n >> 24)
	b[3] = byte(n >> 32)
	b[2] = byte(n >> 40)
	b[1] = byte(n >> 48)
	b[0] = byte(n >> 56)
}
