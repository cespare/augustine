package augustine

type key uint64

func (k key) Bytes() []byte {
	b := make([]byte, 8)
	b[0] = byte(k >> 56)
	b[1] = byte(k >> 48)
	b[2] = byte(k >> 40)
	b[3] = byte(k >> 32)
	b[4] = byte(k >> 24)
	b[5] = byte(k >> 16)
	b[6] = byte(k >> 8)
	b[7] = byte(k)
	return b
}

func keyFromBytes(b []byte) key {
	if len(b) != 8 {
		panic("not a key")
	}
	return key(b[7]) | key(b[6])<<8 | key(b[5])<<16 | key(b[4])<<24 |
		key(b[3])<<32 | key(b[2])<<40 | key(b[1])<<48 | key(b[0])<<56
}
