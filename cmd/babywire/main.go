package main

import (
	"fmt"
	"syscall"
)

func main() {
	fd, error := syscall.Socket(
		syscall.AF_PACKET,
		syscall.SOCK_RAW,
		int(HostToNetShort(syscall.ETH_P_ALL)),
	)
	defer syscall.Close(fd)
	if error != nil {
		panic(error)
	}

	err := syscall.BindToDevice(fd, "enp0s3")
	if err != nil {
		panic(err)
	}

	data := make([]byte, 1024)
	for {
		syscall.Recvfrom(fd, data, 0)
		fmt.Println(data)
		ethHeader := EthernetHeader{
			DestinationAddress: [6]byte(data[0:6]),
			SourceAddress:      [6]byte(data[6:12]),
			EtherType:          uint16(data[12] | data[13]),
		}
		fmt.Println("Destination Address: ", ethHeader.DestinationAddress)
		fmt.Println("Source Address: ", ethHeader.SourceAddress)
		fmt.Println("EtherType: ", ethHeader.EtherType)
	}
}

// HostToNetShort converts a uint16 host byte order to a uint16 network byte order.
func HostToNetShort(host uint16) uint16 {
	return (host&0xFF00)>>8 | (host&0x00FF)<<8
}

type EthernetHeader struct {
	DestinationAddress [6]byte
	SourceAddress      [6]byte
	EtherType          uint16
}
