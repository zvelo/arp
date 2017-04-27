package arp

import (
	"sync"
	"time"
)

type cache struct {
	sync.RWMutex
	table  ArpTable
	table2 ArpTable2

	Updated      time.Time
	UpdatedCount int
}

func (c *cache) Refresh() {
	c.Lock()
	defer c.Unlock()

	c.table, c.table2 = Table12()
	c.Updated = time.Now()
	c.UpdatedCount += 1
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
