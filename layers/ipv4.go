package layers

import (
	"encoding/binary"
	"net"
)

type IPv4Protocol uint8

func (p IPv4Protocol) GetNextLayer() Layer {
	return nil
}

type IPv4 struct {
	VersionAndHeaderLength uint8
	TotalLength            uint16
	TimeToLive             uint8
	Protocol               IPv4Protocol
	HeaderChecksum         uint16
	SrcAddr                net.IP
	DstAddr                net.IP
	Payload                []byte
}

func (i IPv4) LayerPayload() []byte {
	return i.Payload
}

func (i IPv4) LayerType() LayerType {
	return i.Protocol
}

func (i *IPv4) Decode(data []byte) error {
	i.VersionAndHeaderLength = data[0]
	i.TotalLength = binary.BigEndian.Uint16(data[2:4])
	i.TimeToLive = data[8]
	i.Protocol = IPv4Protocol(data[9])
	i.HeaderChecksum = binary.BigEndian.Uint16(data[10:12])
	i.SrcAddr = net.IP(data[12:16])
	i.DstAddr = net.IP(data[16:20])

	headerLength := (i.VersionAndHeaderLength & 0x0F) * 32 / 8
	i.Payload = data[headerLength:i.TotalLength]

	return nil
}
