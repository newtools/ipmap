package ipmap

import (
	"math/big"
	"net"
)

// IPv4Map is a lookup table for ip addresses
type IPv4Map struct {
	bitMap   big.Int
	length   uint16
	subAddrs []*IPv4Map
}

// Set sets and ipv4 address as being prevent in the map
func (ipt *IPv4Map) Set(ip net.IP) {
	ipt.set(ip[12:], 3)
}

func (ipt *IPv4Map) set(ip []byte, dec int8) {
	if dec >= 0 {
		ipt.bitMap.SetBit(&ipt.bitMap, int(ip[0]), 1)
		if dec > 0 {
			idx := ip[1]
			ipIL := uint16(idx) + 1
			if ipIL > ipt.length {
				sl := make([]*IPv4Map, ipIL-ipt.length)
				ipt.subAddrs = append(ipt.subAddrs, sl...)
				ipt.length = ipIL
			}
			if ipt.subAddrs[idx] == nil {
				ipt.subAddrs[idx] = new(IPv4Map)
			}
			ipt.subAddrs[idx].set(ip[1:], dec-1)
		}
	}
}

// IsSet returns true if the ip address is present
func (ipt *IPv4Map) IsSet(ip net.IP) bool {
	return ipt.isSet(ip[12:], 3)
}

func (ipt *IPv4Map) isSet(ip []byte, dec int8) bool {
	if dec >= 0 {
		if ipt.bitMap.Bit(int(ip[0])) == 1 {
			if dec > 0 {
				idx := ip[1]
				ipIL := uint16(idx) + 1
				if ipIL > ipt.length || ipt.subAddrs[idx] == nil {
					return false
				}
				return ipt.subAddrs[idx].isSet(ip[1:], dec-1)
			}
			return true
		}
	}
	return false
}

// Unset unsets an ip address from the map, returns true if it was successful.
func (ipt *IPv4Map) Unset(ip net.IP) bool {
	unset, _ := ipt.unset(ip[12:], 3)
	return unset
}

func (ipt *IPv4Map) unset(ip []byte, dec uint8) (bool, bool) {
	if dec >= 0 {
		bit := int(ip[0])
		if ipt.bitMap.Bit(bit) == 1 {
			if dec > 0 {
				idx := ip[1]
				ipIL := uint16(idx) + 1
				if ipIL > ipt.length || ipt.subAddrs[idx] == nil {
					return false, false
				}
				childUnset, childEmpty := ipt.subAddrs[idx].unset(ip[1:], dec-1)
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
