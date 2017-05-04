package arp

import (
	"time"
)

type ArpInfo struct {
	IPAddr string
	HWType string
	Flags  string
	HWAddr string
	Mask   string
	Device string
}

type ArpTable map[string]string
type ArpTable2 map[string]ArpInfo

var (
	stop     = make(chan struct{})
	arpCache = &cache{
		table:  make(ArpTable),
		table2: make(ArpTable2),
	}
)

func AutoRefresh(t time.Duration) {
	go func() {
		for {
			select {
			case <-time.After(t):
				arpCache.Refresh()
			case <-stop:
				return
			}
		}
	}()
}

func StopAutoRefresh() {
	stop <- struct{}{}
}

func CacheUpdate() {
	arpCache.Refresh()
}

func CacheLastUpdate() time.Time {
	return arpCache.Updated
}

func CacheUpdateCount() int {
	return arpCache.UpdatedCount
}

func CacheIncludeLocal() {
	arpCache.IncludeLocal = true
}

// Search looks up the MAC address for an IP address
// in the arp table
func Search(ip string) string {
	return arpCache.Search(ip)
}

func Search2(ip string) ArpInfo {
	return arpCache.Search2(ip)
}
