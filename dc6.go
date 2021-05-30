package d2dc6

import (
	"fmt"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"

	"github.com/gucio321/d2dc6/d2dc6frame"

	"github.com/gravestench/bitstream"
)

const (
	endOfScanLine = 0x80
	maxRunLength  = 0x7f

	terminationSize = 4

	bytesPerInt32 = 4
)

type scanlineState int

const (
	endOfLine scanlineState = iota
	runOfTransparentPixels
	runOfOpaquePixels
)

// DC6 represents a DC6 file.
type DC6 struct {
	Version            int32
	Flags              uint32
	Encoding           uint32
	Termination        []byte // 4 bytes
	Directions         uint32
	FramesPerDirection uint32
	FramePointers      []uint32                 // size is Directions*FramesPerDirection
	Frames             [][]*d2dc6frame.DC6Frame // size is Directions*FramesPerDirection
}

// New creates a new, empty DC6
func New() *DC6 {
	result := &DC6{
		Version:            0,
		Flags:              0,
		Encoding:           0,
		Termination:        make([]byte, 4),
		Directions:         0,
		FramesPerDirection: 0,
		FramePointers:      make([]uint32, 0),
		Frames:             make([][]*d2dc6frame.DC6Frame, 0),
	}

	return result
}

// Load loads a dc6 animation
func Load(data []byte) (*DC6, error) {
	d := New()

	if err := d.Load(data); err != nil {
		return nil, err
	}

	return d, nil
}

// Load converts bite slice into DC6 structure
func (d *DC6) Load(data []byte) error {
	var err error

	r := bitstream.FromBytes(data...)

	err = d.loadHeader(r)
	if err != nil {
		return err
	}

	frameCount := int(d.Directions * d.FramesPerDirection)
	fmt.Println(frameCount)

	d.FramePointers = make([]uint32, frameCount)
	for i := 0; i < frameCount; i++ {
		d.FramePointers[i], err = r.Next(bytesPerInt32).Bytes().AsUInt32()
		if err != nil {
			return err
		}
	}

	d.Frames = make([][]*d2dc6frame.DC6Frame, d.Directions)
	for i := range d.Frames {
		d.Frames[i] = make([]*d2dc6frame.DC6Frame, d.FramesPerDirection)
	}

	if err := d.loadFrames(r); err != nil {
		return fmt.Errorf("error loading frames: %w", err)
	}

	return nil
}

func (d *DC6) loadHeader(r *bitstream.BitStream) error {
	var err error

	r.Next(bytesPerInt32) // set readed data size to 4 bytes

	if d.Version, err = r.Bytes().AsInt32(); err != nil {
		return err
	}

	if d.Flags, err = r.Bytes().AsUInt32(); err != nil {
		return err
	}

	if d.Encoding, err = r.Bytes().AsUInt32(); err != nil {
		return err
	}

	if d.Termination, err = r.Next(terminationSize).Bytes().AsBytes(); err != nil {
		return err
	}

	r.Next(bytesPerInt32) // set readed data size to 4 bytes

	if d.Directions, err = r.Bytes().AsUInt32(); err != nil {
		return err
	}

	if d.FramesPerDirection, err = r.Bytes().AsUInt32(); err != nil {
		return err
	}

	return nil
}

func (d *DC6) loadFrames(r *bitstream.BitStream) error {
	var err error

	for dir := range d.Frames {
		// for i := 0; i < len(d.FramePointers); i++ {
		for f := range d.Frames[dir] {
			d.Frames[dir][f], err = d2dc6frame.Load(r)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Marshal encodes dc6 animation back into byte slice
func (d *DC6) Marshal() []byte {
	sw := d2datautils.CreateStreamWriter()

	// Encode header
	sw.PushInt32(d.Version)
	sw.PushUint32(d.Flags)
	sw.PushUint32(d.Encoding)

	sw.PushBytes(d.Termination...)

	sw.PushUint32(d.Directions)
	sw.PushUint32(d.FramesPerDirection)

	// load frames
	for _, i := range d.FramePointers {
		sw.PushUint32(i)
	}

	for dir := range d.Frames {
		for f := range d.Frames[dir] {
			data := d.Frames[dir][f].Encode()
			sw.PushBytes(data...)
		}
	}

	return sw.GetBytes()
}

// DecodeFrame decodes the given frame to an indexed color texture
func (d *DC6) DecodeFrame(directionIndex, frameIndex int) []byte {
	frame := d.Frames[directionIndex][frameIndex]

	indexData := make([]byte, frame.Width*frame.Height)
	x := 0
	y := int(frame.Height) - 1
	offset := 0

loop: // this is a label for the loop, so the switch can break the loop (and not the switch)
	for {
		b := int(frame.FrameData[offset])
		offset++

		switch scanlineType(b) {
		case endOfLine:
			if y == 0 {
				break loop
			}

			y--

			x = 0
		case runOfTransparentPixels:
			transparentPixels := b & maxRunLength
			x += transparentPixels
		case runOfOpaquePixels:
			for i := 0; i < b; i++ {
				indexData[x+y*int(frame.Width)+i] = frame.FrameData[offset]
				offset++
			}

			x += b
		}
	}

	return indexData
}

func scanlineType(b int) scanlineState {
	if b == endOfScanLine {
		return endOfLine
	}

	if (b & endOfScanLine) > 0 {
		return runOfTransparentPixels
	}

	return runOfOpaquePixels
}

// Clone creates a copy of the DC6
func (d *DC6) Clone() *DC6 {
	clone := *d
	copy(clone.Termination, d.Termination)
	copy(clone.FramePointers, d.FramePointers)
	clone.Frames = make([][]*d2dc6frame.DC6Frame, len(d.Frames))

	for dir := range clone.Frames {
		clone.Frames[dir] = make([]*d2dc6frame.DC6Frame, len(d.Frames[dir]))
	}

	for dir := range d.Frames {
		for f := range d.Frames[dir] {
			cloneFrame := *d.Frames[dir][f]
			clone.Frames[dir][f] = &cloneFrame
		}
	}

	return &clone
}
