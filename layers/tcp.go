package layers

import (
	"encoding/binary"
	"fmt"
)

type LayerTypeTCP struct{}

func (l LayerTypeTCP) GetNextLayer() Layer {
	return nil
}

type TCP struct {
	SrcPort                                    uint16
	DestPort                                   uint16
	SeqNum                                     uint32
	AckNum                                     uint32
	DataOffsetAndReserved                      uint8
	FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS bool
	Window                                     uint16
	Checksum                                   uint16
	Payload                                    []byte
}

func (t TCP) LayerPayload() []byte {
	return t.Payload
}

func (t TCP) LayerType() LayerType {
	return LayerTypeTCP{}
}

func (t *TCP) Decode(data []byte) error {
	t.SrcPort = binary.BigEndian.Uint16(data[0:2])
	t.DestPort = binary.BigEndian.Uint16(data[2:4])
	t.SeqNum = binary.BigEndian.Uint32(data[4:8])
	t.AckNum = binary.BigEndian.Uint32(data[8:12])
	t.DataOffsetAndReserved = data[12]
	t.FIN = (data[13] & 0x01) != 0
	t.SYN = (data[13] & 0x02) != 0
	t.RST = (data[13] & 0x04) != 0
	t.PSH = (data[13] & 0x08) != 0
	t.ACK = (data[13] & 0x10) != 0
	t.URG = (data[13] & 0x20) != 0
	t.ECE = (data[13] & 0x40) != 0
	t.CWR = (data[13] & 0x80) != 0
	t.Window = binary.BigEndian.Uint16(data[14:16])
	t.Checksum = binary.BigEndian.Uint16(data[16:18])

	dataOffset := t.DataOffsetAndReserved >> 4
	t.Payload = data[dataOffset*4:]

	return nil
}

func (t TCP) String() string {
	return fmt.Sprintf("; %d > %d", t.SrcPort, t.DestPort)
}
