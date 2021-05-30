package d2dc6frame

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"

	"github.com/gravestench/bitstream"
)

const (
	bytesPerInt32  = 4
	terminatorSize = 3
)

// DC6Frame represents a single frame in a DC6.
type DC6Frame struct {
	Flipped    uint32
	Width      uint32
	Height     uint32
	OffsetX    int32
	OffsetY    int32
	Unknown    uint32
	NextBlock  uint32
	FrameData  []byte
	Terminator []byte // 3 bytes
}

// Load loads frame data
func Load(r *bitstream.BitStream) (*DC6Frame, error) {
	var err error

	frame := &DC6Frame{}

	r.Next(bytesPerInt32) // set bytes len to uint32

	if frame.Flipped, err = r.Bytes().AsUInt32(); err != nil {
		return nil, err
	}

	if frame.Width, err = r.Bytes().AsUInt32(); err != nil {
		return nil, err
	}

	if frame.Height, err = r.Bytes().AsUInt32(); err != nil {
		return nil, err
	}

	if frame.OffsetX, err = r.Bytes().AsInt32(); err != nil {
		return nil, err
	}

	if frame.OffsetY, err = r.Bytes().AsInt32(); err != nil {
		return nil, err
	}

	if frame.Unknown, err = r.Bytes().AsUInt32(); err != nil {
		return nil, err
	}

	if frame.NextBlock, err = r.Bytes().AsUInt32(); err != nil {
		return nil, err
	}

	l, err := r.Bytes().AsUInt32()
	if err != nil {
		return nil, err
	}

	if frame.FrameData, err = r.Next(int(l)).Bytes().AsBytes(); err != nil {
		return nil, err
	}

	if frame.Terminator, err = r.Next(terminatorSize).Bytes().AsBytes(); err != nil {
		return nil, err
	}

	return frame, nil
}

// Encode encodes frame data into a byte slice
func (f *DC6Frame) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()
	sw.PushUint32(f.Flipped)
	sw.PushUint32(f.Width)
	sw.PushUint32(f.Height)
	sw.PushInt32(f.OffsetX)
	sw.PushInt32(f.OffsetY)
	sw.PushUint32(f.Unknown)
	sw.PushUint32(f.NextBlock)
	sw.PushUint32(uint32(len(f.FrameData)))
	sw.PushBytes(f.FrameData...)
	sw.PushBytes(f.Terminator...)

	return sw.GetBytes()
}
