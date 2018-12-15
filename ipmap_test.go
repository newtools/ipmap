package ipmap

import (
	"net"
	"testing"
)

type ipv4 [16]byte

func TestIPMapSet(t *testing.T) {
	t.Parallel()
	ipMap := NewIPMap(false)
	loop := net.IPv4(127, 0, 0, 1)
	ipMap.Set(loop)
	if !ipMap.IsSet(loop) {
		t.Fatalf("ip: %s, was set on map, but came back as not set", loop)
	}
}

func TestIPMapSimpleUnSet(t *testing.T) {
	t.Parallel()
	ipMap := NewIPMap(false)
	loop := net.IPv4(127, 0, 0, 254)
	ipMap.Set(loop)
	if !ipMap.IsSet(loop) {
		t.Fatalf("ip: %s, was set on map, but came back as not set", loop)
	}
	if !ipMap.Unset(loop) {
		t.Fatalf("ip: %s, was set on map and failed to unset", loop)
	}
	if ipMap.IsSet(loop) {
		t.Fatalf("ip: %s, was unset on map, but came back as set", loop)
	}
}

func TestIPMapComplexUnSet(t *testing.T) {
	t.Parallel()
	ipMap := NewIPMap(false)
	loop := IPv6(0, 0, 0, 0, 0, 0, 0, 1)
	loop1 := IPv6(0, 0, 0, 0, 0, 0, 0, 0xfffe)
	ipMap.Set(loop)
	ipMap.Set(loop1)
	if !ipMap.IsSet(loop) {
		t.Fatalf("ip: %s, was set on map, but came back as not set", loop)
	}
	if !ipMap.IsSet(loop1) {
		t.Fatalf("ip: %s, was set on map, but came back as not set", loop1)
	}
	if !ipMap.Unset(loop) {
		t.Fatalf("ip: %s, was set on map and failed to unset", loop)
	}
	if ipMap.IsSet(loop) {
		t.Fatalf("ip: %s, was unset on map, but came back as set", loop)
	}
	if !ipMap.IsSet(loop1) {
		t.Fatalf("ip: %s, was set on map, but came back as not set", loop1)
	}
	if !ipMap.Unset(loop1) {
		t.Fatalf("ip: %s, was set on map and failed to unset", loop1)
	}
	if ipMap.IsSet(loop1) {
		t.Fatalf("ip: %s, was unset on map, but came back as set", loop1)
	}
}

func BenchmarkIPSet(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var ipArr []net.IP
	// 0000.0000.0000.0000.0000.0000.0000.0000.0000.0000
	var ar [8]uint16
	var cur int
	for n := 0; n < b.N; n++ {
		cur = ip6Add(&ar, cur)
		ip := IPv6(ar[0], ar[1], ar[2], ar[3], ar[4], ar[5], ar[6], ar[7])
		ipArr = append(ipArr, ip)
	}
	ipMap := NewIPMap(false)
	b.ResetTimer()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		ipMap.Set(ipArr[n])
	}
}

func BenchmarkIPUnset(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var ipArr []net.IP
	ipMap := NewIPMap(false)
	// 0000.0000.0000.0000.0000.0000.0000.0000.0000.0000
	var ar [8]uint16
	var cur int
	for n := 0; n < b.N; n++ {
		cur = ip6Add(&ar, cur)
		ip := IPv6(ar[0], ar[1], ar[2], ar[3], ar[4], ar[5], ar[6], ar[7])
		ipArr = append(ipArr, ip)
		ipMap.Set(ip)
	}
	b.ResetTimer()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		ipMap.Unset(ipArr[n])
	}
}

func BenchmarkIPIsSet(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var ipArr []net.IP
	ipMap := NewIPMap(false)
	// 0000.0000.0000.0000.0000.0000.0000.0000.0000.0000
	var ar [8]uint16
	var cur int
	for n := 0; n < b.N; n++ {
		cur = ip6Add(&ar, cur)
		ip := IPv6(ar[0], ar[1], ar[2], ar[3], ar[4], ar[5], ar[6], ar[7])
		ipArr = append(ipArr, ip)
		ipMap.Set(ip)
	}
	b.ResetTimer()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		if !ipMap.IsSet(ipArr[n]) {
			b.Fatalf("is set reported %s as unset, when it was set", ipArr[n])
		}
	}
}

func BenchmarkIPSetv4Only(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var ipArr []net.IP
	// 0.0.0.0
	var ar [4]uint8
	var cur int
	for n := 0; n < b.N; n++ {
		cur = ip4Add(&ar, cur)
		ip := net.IPv4(ar[0], ar[1], ar[2], ar[3])
		ipArr = append(ipArr, ip)
	}
	ipMap := NewIPMap(true)
	b.ResetTimer()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		ipMap.Set(ipArr[n])
	}
}

func BenchmarkIPUnsetv4Only(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var ipArr []net.IP
	ipMap := NewIPMap(true)
	// 0.0.0.0
	var ar [4]uint8
	var cur int
	for n := 0; n < b.N; n++ {
		cur = ip4Add(&ar, cur)
		ip := net.IPv4(ar[0], ar[1], ar[2], ar[3])
		ipArr = append(ipArr, ip)
		ipMap.Set(ip)
	}
	b.ResetTimer()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		ipMap.Unset(ipArr[n])
	}
}

func BenchmarkIPIsSetv4Only(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var ipArr []net.IP
	ipMap := NewIPMap(true)
	// 0.0.0.0
	var ar [4]uint8
	var cur int
	for n := 0; n < b.N; n++ {
		cur = ip4Add(&ar, cur)
		ip := net.IPv4(ar[0], ar[1], ar[2], ar[3])
		ipArr = append(ipArr, ip)
		ipMap.Set(ip)
	}
	b.ResetTimer()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		if !ipMap.IsSet(ipArr[n]) {
			b.Fatalf("is set reported %s as unset, when it was set", ipArr[n])
		}
	}
}

func BenchmarkIP4SetOnMap(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var ipArr []ipv4
	// 0.0.0.0
	var ar [4]uint8
	var cur int
	for n := 0; n < b.N; n++ {
		cur = ip4Add(&ar, cur)
		ipArr = append(ipArr, ipv4{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, ar[0], ar[1], ar[2], ar[3]})
	}
	ipMap := make(map[ipv4]struct{})
	b.ResetTimer()
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		ipMap[ipArr[n]] = struct{}{}
	}
}

func BenchmarkIP4IsSetOnMap(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	var ipArr []ipv4
	ipMap := make(map[ipv4]struct{})
	// 0.0.0.0
	var ar [4]uint8
	var cur int
	for n := 0; n < b.N; n++ {
		cur = ip4Add(&ar, cur)
		ipv4t := ipv4{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, ar[0], ar[1], ar[2], ar[3]}
		ipArr = append(ipArr, ipv4t)
		ipMap[ipv4t] = struct{}{}
	}

	b.ResetTimer()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		if _, ok := ipMap[ipArr[n]]; !ok {
			b.Fatalf("is set reported %s as unset, when it was set", ipArr[n])
		}
	}
}

func ip4Add(ip *[4]uint8, cur int) int {
	ipC := *ip
	var set bool
	l := 4
	for t := l - 1; t > cur; t-- {
		if ipC[t] < 255 {
			ipC[t]++
			set = true
			break
		}
	}
	if !set {
		if ipC[cur] == 255 {
			cur++
		}
		if cur == l {
			cur = 0
			ipC[0] = 0
		}
		ipC[cur]++
		for i := cur + 1; i < l; i++ {
			ipC[i] = 0
		}
	}
	*ip = ipC
	return cur
}

func ip6Add(ip *[8]uint16, cur int) int {
	ipC := *ip
	var set bool
	l := 8
	for t := l - 1; t > cur; t-- {
		if ipC[t] < 65535 {
			ipC[t]++
			set = true
			break
		}
	}
	if !set {
		if ipC[cur] == 65535 {
			cur++
		}
		if cur == l {
			cur = 0
			ipC[0] = 0
		}
		ipC[cur]++
		for i := cur + 1; i < l; i++ {
			ipC[i] = 0
		}
	}
	*ip = ipC
	return cur
}

func incrementUint16(v *uint16) bool {
	if *v < 65535 {
		*v++
		return true
	}
	return false
}
