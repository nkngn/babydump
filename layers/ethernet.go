package layers

import "fmt"

type MacAddress []byte

func (m MacAddress) String() string {
	if len(m) != 6 {
		return "invalid"
	}
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", m[0], m[1], m[2], m[3], m[4], m[5])
}

type EtherType uint16

func (e EtherType) VerboseString() string {
	switch e {
	case 0x0800:
		return "IPv4"
	case 0x0806:
		return "ARP"
	case 0x86DD:
		return "IPv6"
	default:
		return "Unknown"
	}
}

type Ethernet struct {
	DstMac    MacAddress
	SrcMac    MacAddress
	EtherType EtherType
}
