package ipmap

const octetSize = 16

// BitMap is a 256 bit bitmap
type BitMap [octetSize]uint64

// Set sets z to x, with x's i'th bit set to b (0 or 1).
// That is, if b is 1 SetBit sets z = x | (1 << i);
// if b is 0 SetBit sets z = x &^ (1 << i). If b is not 0 or 1,
// SetBit will panic.
func (bm *BitMap) Set(i uint, v uint8) *BitMap {
	j := i / octetSize
	m := uint64(1 << (i % octetSize))
	switch v {
	case 0:
		(*bm)[j] &^= m
		return bm
	case 1:
		(*bm)[j] |= m
		return bm
	}
	panic("bit is not 0 or 1")
}

// IsSet returns the true if the i'th bit of x is 1. That is, it
// returns ((x>>i)&1 == 0).
func (bm *BitMap) IsSet(i uint) bool {
	j := i / octetSize
	m := uint64(1 << (i % octetSize))
	return (*bm)[j]&m == m
}
