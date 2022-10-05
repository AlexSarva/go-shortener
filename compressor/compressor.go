package compressor

// Compressor is the interface to different compress types
type Compressor interface {
	Compress() []byte
}
