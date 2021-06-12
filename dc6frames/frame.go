package dc6frames

import (
	"fmt"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"

	"github.com/gravestench/bitstream"
)

const (
	bytesPerInt32  = 4
	terminatorSize = 3
)

func newFrame() *Frame {
	return &Frame{
		FrameData:  make([]byte, 0),
		Terminator: make([]byte, terminatorSize),
	}
}

// Frame represents a single frame in a DC6.
type Frame struct {
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
func (frame *Frame) Load(r *bitstream.BitStream) error {
	var err error

	r.Next(bytesPerInt32) // set bytes len to uint32

	if frame.Flipped, err = r.Bytes().AsUInt32(); err != nil {
		return fmt.Errorf("reading flipped: %w", err)
	}

	if frame.Width, err = r.Bytes().AsUInt32(); err != nil {
		return fmt.Errorf("reading width: %w", err)
	}

	if frame.Height, err = r.Bytes().AsUInt32(); err != nil {
		return fmt.Errorf("reading height: %w", err)
	}

	if frame.OffsetX, err = r.Bytes().AsInt32(); err != nil {
		return fmt.Errorf("reading x-offset: %w", err)
	}

	if frame.OffsetY, err = r.Bytes().AsInt32(); err != nil {
		return fmt.Errorf("reading y-offset: %w", err)
	}

	if frame.Unknown, err = r.Bytes().AsUInt32(); err != nil {
		return fmt.Errorf("reading frame unknown: %w", err)
	}

	if frame.NextBlock, err = r.Bytes().AsUInt32(); err != nil {
		return fmt.Errorf("reading next block: %w", err)
	}

	l, err := r.Bytes().AsUInt32()
	if err != nil {
		return fmt.Errorf("reading length of frame data: %w", err)
	}

	if frame.FrameData, err = r.Next(int(l)).Bytes().AsBytes(); err != nil {
		return fmt.Errorf("reading frame data: %w", err)
	}

	if frame.Terminator, err = r.Next(terminatorSize).Bytes().AsBytes(); err != nil {
		return fmt.Errorf("reading terminator: %w", err)
	}

	return nil
}

// Encode encodes frame data into a byte slice
func (f *Frame) Encode() []byte {
	fmt.Println("encode")
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
