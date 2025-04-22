package layers

import (
	"encoding/binary"
	"net"
)

type EtherType uint16

func (e EtherType) GetNextLayer() Layer {
	switch e {
	case 0x0800:
		return &IPv4{}
	// case 0x0806:
	// 	return &ARP{}
	// case 0x86DD:
	// 	return &IPv6{}
	default:
		return nil
	}
}

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

type EthernetFrame struct {
	DstMac    net.HardwareAddr
	SrcMac    net.HardwareAddr
	EtherType LayerType
	Payload   []byte
}

func (e EthernetFrame) LayerPayload() []byte {
	return e.Payload
}

func (e EthernetFrame) LayerType() LayerType {
	return e.EtherType
}

func (e *EthernetFrame) Decode(data []byte) error {
	e.DstMac = net.HardwareAddr(data[0:6])
	e.SrcMac = net.HardwareAddr(data[6:12])
	e.EtherType = EtherType(binary.BigEndian.Uint16(data[12:14]))
	e.Payload = data[14:]

	return nil
}

func (e EthernetFrame) String() string {
	return e.DstMac.String() + " -> " + e.SrcMac.String() + ", ethertype " + (e.EtherType.(EtherType)).VerboseString() + " :"
}
