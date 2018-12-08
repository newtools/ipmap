package ipmap

import (
	"net"
	"testing"
)

func TestIPMapSet(t *testing.T) {
	t.Parallel()
	ipMap := new(IPv4Map)
	loop := net.IPv4(127, 0, 0, 1)
	ipMap.Set(loop)
	if !ipMap.IsSet(loop) {
		t.Fatalf("ip: %s, was set on map, but came back as not set", loop)
	}
}

func TestIPMapSimpleUnSet(t *testing.T) {
	t.Parallel()
	ipMap := new(IPv4Map)
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
	ipMap := new(IPv4Map)
	loop := net.IPv4(127, 0, 0, 1)
	loop1 := net.IPv4(127, 0, 0, 254)
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
	// 1.0.0.0.1
	var ipBits uint32 = 16777217
	for n := 0; n < b.N; n++ {
		one := uint8((ipBits & 0xff000000) >> 24)
		two := uint8((ipBits & 0xff0000) >> 16)
		three := uint8((ipBits & 0xff00) >> 8)
		four := uint8(ipBits & 0xff)
		ipArr = append(ipArr, net.IPv4(one, two, three, four))
		ipBits++
	}
	ipMap := new(IPv4Map)
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
	// 1.0.0.0.1
	var ipBits uint32 = 16777217
	ipMap := new(IPv4Map)
	for n := 0; n < b.N; n++ {
		one := uint8((ipBits & 0xff000000) >> 24)
		two := uint8((ipBits & 0xff0000) >> 16)
		three := uint8((ipBits & 0xff00) >> 8)
		four := uint8(ipBits & 0xff)
		ipArr = append(ipArr, net.IPv4(one, two, three, four))
		ipBits++
		ipMap.Set(ipArr[n])
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
	// 1.0.0.0.1
	var ipBits uint32 = 16777217
	ipMap := new(IPv4Map)
	for n := 0; n < b.N; n++ {
		one := uint8((ipBits & 0xff000000) >> 24)
		two := uint8((ipBits & 0xff0000) >> 16)
		three := uint8((ipBits & 0xff00) >> 8)
		four := uint8(ipBits & 0xff)
		ipArr = append(ipArr, net.IPv4(one, two, three, four))
		ipBits++
		ipMap.Set(ipArr[n])
	}
	b.ResetTimer()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		if !ipMap.IsSet(ipArr[n]) {
			b.Fatalf("is set reported %s as unset, when it was set", ipArr[n])
		}
	}
}
