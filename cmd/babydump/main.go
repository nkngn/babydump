package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"

	babydump "github.com/ngockhanhnguyen/babydump/packet"
)

func main() {
	// interfaces, err := net.Interfaces()
	// for _, i := range interfaces {
	// 	fmt.Printf("%s %d\n", i.Name, i.MTU)
	// }
	// Open a raw socket
	fd, error := syscall.Socket(
		syscall.AF_PACKET,
		syscall.SOCK_RAW,
		int(HostToNetShort(syscall.ETH_P_ALL)),
	)
	defer syscall.Close(fd)
	if error != nil {
		panic(error)
	}

	iface, err := net.InterfaceByName("enp0s3")
	if err != nil {
		log.Fatal("Interface enp0s3 not found")
	}
	println(iface.MTU)

	// giới hạn gói tin gửi đi, nhưng không luôn đảm bảo giới hạn gói tin nhận được (tùy kernel version)
	// err = syscall.BindToDevice(fd, iface.Name)
	// if err != nil {
	// 	panic(err)
	// }

	// thực sự lọc nhận từ interface cụ thể, vì đang bind tại layer 2, không chỉ là layer socket.
	sll := &syscall.SockaddrLinklayer{
		Protocol: HostToNetShort(syscall.ETH_P_ALL),
		Ifindex:  iface.Index,
	}

	err = syscall.Bind(fd, sll)
	if err != nil {
		log.Fatal("Failed to bind to interface:", err)
	}

	data := make([]byte, iface.MTU+18)
	for {
		syscall.Recvfrom(fd, data, 0)
		fmt.Printf("%v\n", data)
		packet := babydump.NewPacket(data, babydump.PackageMetadata{
			Timestamp:      time.Now(),
			InterfaceIndex: iface.Index,
		})
		// fmt.Printf("%d\n", len(packet.Layers()))
		fmt.Printf("%s\n", packet.String())
		break

		// var ethFrame EthernetFrame
		// ethFrame.ReceivedAt = time.Now()
		// err := ethFrame.UnmarshalBinary(data)
		// if err != nil {
		// 	log.Fatal("decode error:", err)
		// }
		// ethFrame.Print()
	}
}

// HostToNetShort converts a uint16 host byte order to a uint16 network byte order.
func HostToNetShort(host uint16) uint16 {
	return (host&0xFF00)>>8 | (host&0x00FF)<<8
}

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

type IPPacket struct {
	Header  IPPacketHeader
	Payload IPPacketPayload
	// Data   []byte
}

func (i *IPPacket) UnmarshalBinary(data []byte) error {
	i.Header.UnmarshalBinary(data[14:30])
	if i.Header.Version() == 4 {
		i.Payload.UnmarshalBinary(data[30+i.Header.IHL()*4:])
	}
	return nil
}

type IPPacketHeader struct {
	VersionAndHeaderLength uint8
	TotalLength            uint16
	TimeToLive             uint8
	Protocol               uint8
	HeaderChecksum         uint16
	SourceAddress          [4]byte
	DestinationAddress     [4]byte
}

func (h *IPPacketHeader) UnmarshalBinary(data []byte) error {
	h.VersionAndHeaderLength = data[0]

	if h.Version() == 4 {
		h.TotalLength = binary.BigEndian.Uint16(data[2:4])
		h.TimeToLive = data[8]
		h.Protocol = data[9]
		h.HeaderChecksum = binary.BigEndian.Uint16(data[10:12])
		copy(h.SourceAddress[:], data[12:16])
		copy(h.DestinationAddress[:], data[16:20])
	}

	return nil
}

func (h *IPPacketHeader) Version() uint8 {
	return h.VersionAndHeaderLength >> 4
}

// Internet Header Length
func (h *IPPacketHeader) IHL() uint8 {
	return h.VersionAndHeaderLength << 4
}

type IPPacketPayload struct {
	SrcPort uint16
	DesPort uint16
}

func (h *IPPacketPayload) UnmarshalBinary(data []byte) error {
	h.SrcPort = binary.BigEndian.Uint16(data[0:2])
	h.DesPort = binary.BigEndian.Uint16(data[2:4])
	return nil
}

type EthernetFrame struct {
	DstMac     MacAddress
	SrcMac     MacAddress
	EtherType  EtherType
	ReceivedAt time.Time

	IPPacket IPPacket
}

func (e *EthernetFrame) UnmarshalBinary(data []byte) error {
	e.DstMac = data[0:6]
	e.SrcMac = make([]byte, 6)
	copy(e.SrcMac, data[6:12])
	e.EtherType = EtherType(binary.BigEndian.Uint16(data[12:14]))

	if e.EtherType == 0x0800 {
		err := e.IPPacket.UnmarshalBinary(data)
		if err != nil {
			return err
		}
	}
	// if e.EtherType >= 0x0600 {
	// 	fmt.Println("Ethernet II frame")
	// } else {
	// 	fmt.Println("IEEE 802.3 frame")
	// }

	return nil
}

func (e EthernetFrame) Print() {
	// 06:46:56.871318 52:54:00:12:35:02 > 02:f6:e6:18:e6:1a, ethertype IPv4 (0x0800), length 60: 10.0.2.2.49880 > 10.0.2.15.22: Flags [.], ack 13576, win 65535, length 0
	s := fmt.Sprintf(
		"%s %s > %s, ethertype %s (%#04x)",
		e.ReceivedAt.Format("15:04:05.000000"),
		e.DstMac,
		e.SrcMac,
		e.EtherType.VerboseString(),
		e.EtherType,
	)

	if e.EtherType == 0x0800 {
		s = fmt.Sprintf("%s, length %d", s, e.IPPacket.Header.TotalLength)
	}

	// if e.IPPacket.Header.TotalLength > 1500 {
	// 	log.Fatalf("packet too long %d", e.IPPacket.Header.TotalLength)
	// }

	if e.IPPacket.Header.Version() == 4 {
		// TCP
		if e.IPPacket.Header.Protocol == 0x06 {
			s = fmt.Sprintf("%s: %d.%d.%d.%d.%d > %d.%d.%d.%d.%d",
				s,
				e.IPPacket.Header.SourceAddress[0],
				e.IPPacket.Header.SourceAddress[1],
				e.IPPacket.Header.SourceAddress[2],
				e.IPPacket.Header.SourceAddress[3],
				e.IPPacket.Payload.SrcPort,
				e.IPPacket.Header.DestinationAddress[0],
				e.IPPacket.Header.DestinationAddress[1],
				e.IPPacket.Header.DestinationAddress[2],
				e.IPPacket.Header.DestinationAddress[3],
				e.IPPacket.Payload.DesPort)
		}
	} else {
		s = fmt.Sprintf("%s: %x > %x", s, e.IPPacket.Header.SourceAddress, e.IPPacket.Header.DestinationAddress)
	}

	println(s)
}
