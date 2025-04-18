package babydump

import (
	"time"

	"github.com/ngockhanhnguyen/babydump/layers"
)

type PackageMetadata struct {
	Timestamp      time.Time
	InterfaceIndex int
}

type Package interface {
	Data() []byte
	String() string
	Layers() []layers.Layer
	Metadata() PackageMetadata
}

type packet struct {
	data   []byte
	layers []layers.Layer

	// last is the last layer added to the packet
	last     layers.Layer
	metadata PackageMetadata
}

func (p packet) Data() []byte {
	return p.data
}

func (p packet) String() string {
	return ""
}

func (p packet) Layers() []layers.Layer {
	return p.layers
}

func (p packet) Metadata() PackageMetadata {
	return p.metadata
}

func (p packet) decode() error {
	// decode each layer from outside -> inside

	// decode ethernet

	// while get next decoder
	// decode

	return nil
}

func NewPacket(data []byte, metadata PackageMetadata) packet {
	p := packet{
		data:     data,
		metadata: metadata,
	}

	p.decode()

	return p
}
