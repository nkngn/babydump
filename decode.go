package babydump

type Decoder interface {
	// Decode decodes the bytes of a packet, sending decoded values and other
	// information to PacketBuilder, and returning an error if unsuccessful.  See
	// the PacketBuilder documentation for more details.
	Decode([]byte) error
}
