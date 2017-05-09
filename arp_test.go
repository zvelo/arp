package arp

import (
	"fmt"
	"net"
	"testing"
	"time"
)

var (
	_ = fmt.Println
)

func TestTable(t *testing.T) {

	table := Table()
	if table == nil {
		t.Errorf("Empty table")
	}
}

func TestCacheInfo(t *testing.T) {
	prevUpdated := CacheLastUpdate().UnixNano()
	prevCount := CacheUpdateCount()

	CacheUpdate()

	if prevUpdated == CacheLastUpdate().UnixNano() {
		t.Error()
	}

	if prevCount == CacheUpdateCount() {
		t.Error()
	}
}

func TestAutoRefresh(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping autorefresh test")
	}
	prevUpdated := CacheLastUpdate().UnixNano()
	prevCount := CacheUpdateCount()

	AutoRefresh(100 * time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	StopAutoRefresh()

	if prevUpdated == CacheLastUpdate().UnixNano() {
		t.Error()
	}

	if prevCount == CacheUpdateCount() {
		t.Error()
	}

	// test to make sure stop worked
	prevUpdated = CacheLastUpdate().UnixNano()
	prevCount = CacheUpdateCount()
	time.Sleep(200 * time.Millisecond)
	if prevUpdated != CacheLastUpdate().UnixNano() {
		t.Error()
	}

	if prevCount != CacheUpdateCount() {
		t.Error()
	}
}

func TestSearch(t *testing.T) {
	table := Table()

	for ip, test := range table {

		result := Search(ip)
		if test != result {
			t.Errorf("expected %s got %s", test, result)
		}
	}
}

func getTestInterface() (result net.Interface, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, inf := range interfaces {
		if inf.Flags&net.FlagUp != 0 && len(inf.HardwareAddr) != 0 {
			if addrs, err2 := inf.Addrs(); err2 == nil && len(addrs) != 0 {
				result = inf
				return
			}
		}
	}
	err = fmt.Errorf("unable to find test interface")
	return
}

func TestLocalOff(t *testing.T) {
	testInf, err := getTestInterface()
	if err != nil {
		t.Error(err)
	}
	addrs, err := testInf.Addrs()
	if err != nil {
		t.Error(err)
	}
	for _, addr := range addrs {
		addrStr := addressToIPString(addr)
		result := Search2(addrStr)
		if result.HWAddr != "" {
			t.Error("expected nothing got", result.HWAddr)
		}
	}
}

func TestLocal(t *testing.T) {
	CacheIncludeLocal()
	testInf, err := getTestInterface()
	if err != nil {
		t.Error(err)
	}
	macStr := testInf.HardwareAddr.String()
	addrs, err := testInf.Addrs()
	if err != nil {
		t.Error(err)
	}
	for _, addr := range addrs {
		addrStr := addressToIPString(addr)
		result := Search2(addrStr)
		if result.HWAddr != macStr {
			t.Error("expected", macStr, "got", result.HWAddr)
		}
	}

	result := Search2("127.0.0.1")
	if result.HWAddr == "" {
		t.Error("expected something for 127.0.0.1 got nothing")
	}

	result = Search2("::1")
	if result.HWAddr == "" {
		t.Error("expected something for ::1 got nothing")
	}
}

func BenchmarkSearch(b *testing.B) {
	table := Table()
	if len(table) == 0 {
		return
	}

	for ip, _ := range Table() {
		for i := 0; i < b.N; i++ {
			Search(ip)
		}

		// using the first key is enough
		break
	}
}
