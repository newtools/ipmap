package ipmap

import (
	"math/big"
	"net"
)

// IPv6 returns a new IPv6 address
func IPv6(a, b, c, d, e, f, g, h uint16) net.IP {
	return net.IP{
		byte(a & 0xFF00),
		byte(a & 0xFF),
		byte(b & 0xFF00),
		byte(b & 0xFF),
		byte(c & 0xFF00),
		byte(c & 0xFF),
		byte(d & 0xFF00),
		byte(d & 0xFF),
		byte(e & 0xFF00),
		byte(e & 0xFF),
		byte(f & 0xFF00),
		byte(f & 0xFF),
		byte(g & 0xFF00),
		byte(g & 0xFF),
		byte(h & 0xFF00),
		byte(h & 0xFF),
	}
}

// ipMap is a lookup table for ip addresses
type ipMap struct {
	bitMap            big.Int
	length            uint16
	subAddrs          []*ipMap
	startingIndex     uint8
	startingDecrement uint8
}

// NewIPMap creates a new IPMap. Regular contruction can work, except if you only want to deal
// with IPv4 addresses. You can pass a boolean to this method to indicate to change the map
// alogorithm to optimize for IPv4 only (IPv6 only has not available optimization, so you needn't
// worry about that). This should create faster setting and getting times.
func NewIPMap(ipv4Only bool) *ipMap {
	var startingIndex uint8
	var startingDecrement uint8 = 15
	if ipv4Only {
		startingIndex = 12
		startingDecrement = 3
	}
	return &ipMap{
		startingIndex:     startingIndex,
		startingDecrement: startingDecrement,
	}
}

// Set sets and ipv4 address as being prevent in the map
func (ipt *ipMap) Set(ip net.IP) {
	ipt.set(ip, ipt.startingIndex, ipt.startingDecrement)
}

func (ipt *ipMap) set(ip []byte, baseIdx, dec uint8) {
	if dec >= 0 {
		ipt.bitMap.SetBit(&ipt.bitMap, int(ip[baseIdx]), 1)
		if dec > 0 {
			plusOne := baseIdx + 1
			idx := ip[plusOne]
			ipIL := uint16(idx) + 1
			if ipIL > ipt.length {
				sl := make([]*ipMap, ipIL-ipt.length)
				ipt.subAddrs = append(ipt.subAddrs, sl...)
				ipt.length = ipIL
			}
			if ipt.subAddrs[idx] == nil {
				ipt.subAddrs[idx] = new(ipMap)
			}
			ipt.subAddrs[idx].set(ip, plusOne, dec-1)
		}
	}
}

// IsSet returns true if the ip address is present
func (ipt *ipMap) IsSet(ip net.IP) bool {
	return ipt.isSet(ip, ipt.startingIndex, ipt.startingDecrement)
}

func (ipt *ipMap) isSet(ip []byte, baseIdx, dec uint8) bool {
	if dec >= 0 {
		if ipt.bitMap.Bit(int(ip[baseIdx])) == 1 {
			if dec > 0 {
				plusOne := baseIdx + 1
				idx := ip[plusOne]
				ipIL := uint16(idx) + 1
				if ipIL > ipt.length || ipt.subAddrs[idx] == nil {
					return false
				}
				return ipt.subAddrs[idx].isSet(ip, plusOne, dec-1)
			}
			return true
		}
	}
	return false
}

// Unset unsets an ip address from the map, returns true if it was successful.
func (ipt *ipMap) Unset(ip net.IP) bool {
	unset, _ := ipt.unset(ip, ipt.startingIndex, ipt.startingDecrement)
	return unset
}

func (ipt *ipMap) unset(ip []byte, baseIdx, dec uint8) (bool, bool) {
	if dec >= 0 {
		bit := int(ip[baseIdx])
		if ipt.bitMap.Bit(bit) == 1 {
			if dec > 0 {
				plusOne := baseIdx + 1
				idx := ip[plusOne]
				ipIL := uint16(idx) + 1
				if ipIL > ipt.length || ipt.subAddrs[idx] == nil {
					return false, false
				}
				childUnset, childEmpty := ipt.subAddrs[idx].unset(ip, plusOne, dec-1)
				if !childUnset {
					return false, false
				}
				var amEmpty bool
				if childEmpty {
					ipt.subAddrs[idx] = nil
					if ipIL == ipt.length {
						var cut uint8
						for i := int(ipt.length) - 1; i >= 0; i-- {
							if ipt.subAddrs[i] != nil {
								cut = uint8(i)
								break
							}
						}
						ipt.subAddrs = ipt.subAddrs[:cut]
						ipt.length = uint16(cut)
						amEmpty = ipt.length == 0
					}
				}
				if amEmpty {
					ipt.bitMap.SetBit(&ipt.bitMap, bit, 0)
				}
				return true, amEmpty
			}
			ipt.bitMap.SetBit(&ipt.bitMap, bit, 0)
			return true, true
		}
	}
	return false, false
}
