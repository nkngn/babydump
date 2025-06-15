package babydump

import (
	"fmt"
	"time"

	"github.com/nkngn/babydump/internal/layers"
)

type PackageMetadata struct {
	Timestamp      time.Time
	InterfaceIndex int
}

type Packet struct {
	rawData  []byte
	layers   []layers.Layer
	last     layers.Layer
	metadata PackageMetadata
}

func (p Packet) Data() []byte {
	return p.rawData
}

func (p Packet) String() string {
	str := ""
	for _, layer := range p.layers {
		str += fmt.Sprintf("%v", layer)
	}
	return str
}

func (p Packet) Layers() []layers.Layer {
	return p.layers
}

func (p Packet) Metadata() PackageMetadata {
	return p.metadata
}

func (p *Packet) decode() error {
	ethernetFrm := &layers.EthernetFrame{}

	ethernetFrm.Decode(p.rawData)
	p.layers = append(p.layers, ethernetFrm)
	p.last = ethernetFrm

	for {
		nextLayer := p.last.LayerType().GetNextLayer()
		if nextLayer != nil {
			nextLayer.Decode(p.last.LayerPayload())
			p.layers = append(p.layers, nextLayer)
			p.last = nextLayer
		} else {
			break
		}
	}

	return nil
}

func NewPacket(data []byte, metadata PackageMetadata) Packet {
	p := Packet{
		rawData:  data,
		metadata: metadata,
	}

	p.decode()

	return p
}
