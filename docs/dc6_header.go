// Package docs contains a documentation files (unused in real module)
package docs

// DC6Header represents the file header of a DC6 file.
// this structure is unused in this module and is only a documentation
type DC6Header struct {
	Termination        []byte `struct:"[4]byte"`
	Version            int32  `struct:"int32"`
	Flags              uint32 `struct:"uint32"`
	Encoding           uint32 `struct:"uint32"`
	Directions         int32  `struct:"int32"`
	FramesPerDirection int32  `struct:"int32"`
}
