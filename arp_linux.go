// +build linux

package arp

import (
	"bufio"
	"os"
	"strings"
)

const (
	f_IPAddr int = iota
	f_HWType
	f_Flags
	f_HWAddr
	f_Mask
	f_Device
)

func Table() ArpTable {
	table, _ := Table12()
	return table
}

func Table2() ArpTable2 {
	_, table := Table12()
	return table
}

func Table12() (ArpTable, ArpTable2) {
	f, err := os.Open("/proc/net/arp")

	if err != nil {
		return nil, nil
	}

	defer f.Close()

	s := bufio.NewScanner(f)
	s.Scan() // skip the field descriptions

	table1 := make(ArpTable)
	table2 := make(ArpTable2)

	for s.Scan() {
		line := s.Text()
		fields := strings.Fields(line)
		entry := ArpInfo{
			IPAddr: fields[f_IPAddr],
			HWType: fields[f_HWAddr],
			Flags:  fields[f_Flags],
			HWAddr: fields[f_HWAddr],
			Mask:   fields[f_Mask],
			Device: fields[f_Device],
		}
		table1[entry.IPAddr] = entry.HWAddr
		table2[entry.IPAddr] = entry
	}

	return table1, table2
}
