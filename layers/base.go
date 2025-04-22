package layers

type Layer interface {
	LayerPayload() []byte
	LayerType() LayerType
	Decode(data []byte) error
}

type LayerType interface {
	GetNextLayer() Layer
}
