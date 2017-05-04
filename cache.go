package arp

import (
	"net"
	"sync"
	"time"
)

type cache struct {
	sync.RWMutex
	table  ArpTable
	table2 ArpTable2

	IncludeLocal bool
	Updated      time.Time
	UpdatedCount int
}

func (c *cache) Refresh() {
	c.Lock()
	defer c.Unlock()

	c.table, c.table2 = Table12()
	if c.IncludeLocal {
		c.RefreshLocal()
	}
	c.Updated = time.Now()
	c.UpdatedCount += 1
}

func addressToIPString(addr net.Addr) string {
	switch x := addr.(type) {
	case *net.IPNet:
		return x.IP.String()
	case *net.IPAddr:
		return x.IP.String()
	}
	return ""
}

func (c *cache) RefreshLocal() {
	allInterfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, inf := range allInterfaces {
		if len(inf.HardwareAddr) == 0 {
			continue
		}
		addresses, err := inf.Addrs()
		if err != nil || len(addresses) == 0 {
			continue
		}
		for _, addr := range addresses {
			macStr := inf.HardwareAddr.String()
			addrStr := addressToIPString(addr)
			c.table[addrStr] = macStr
			c.table2[addrStr] = ArpInfo{
				IPAddr: addrStr,
				HWAddr: macStr,
				Flags:  "0x2",
				Device: inf.Name,
			}
		}
	}
}

func (c *cache) Search(ip string) string {
	c.RLock()
	defer c.RUnlock()

	mac, ok := c.table[ip]

	if !ok {
		c.RUnlock()
		c.Refresh()
		c.RLock()
		mac = c.table[ip]
	}

	return mac
}

func (c *cache) Search2(ip string) ArpInfo {
	c.RLock()
	defer c.RUnlock()

	info, ok := c.table2[ip]

	if !ok {
		c.RUnlock()
		c.Refresh()
		c.RLock()
		info = c.table2[ip]
	}

	return info
}
