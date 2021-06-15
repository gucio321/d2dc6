package d2dc6

import (
	"errors"
	"fmt"

	"github.com/gravestench/bitstream"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"

	"github.com/gucio321/d2dc6/pkg/dc6/frames"
)

const (
	endOfScanLine = 0x80
	maxRunLength  = 0x7f

	terminationSize = 4

	bytesPerInt32 = 4

	expectedDC6Version = 6
)

type scanlineState int

const (
	endOfLine scanlineState = iota
	runOfTransparentPixels
	runOfOpaquePixels
)

// DC6 represents a DC6 file.
type DC6 struct {
	Frames      *frames.Frames
	Flags       uint32
	Encoding    uint32
	Termination [terminationSize]byte
}

// New creates a new, empty DC6
func New() *DC6 {
	result := &DC6{
		Flags:    0,
		Encoding: 0,
		Frames:   frames.New(),
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

	r := bitstream.ReaderFromBytes(data...)

	err = d.loadHeader(r)
	if err != nil {
		return err
	}

	frameCount := d.Frames.NumberOfDirections() * d.Frames.FramesPerDirection()

	// frame pointers - skip
	_, err = r.Next(frameCount * bytesPerInt32).Bytes().AsBytes()
	if err != nil {
		return fmt.Errorf("reading frame pointers: %w", err)
	}

	return d.loadFrames(r)
}

func (d *DC6) loadHeader(r *bitstream.Reader) error {
	var err error

	r.Next(bytesPerInt32) // set readed data size to 4 bytes

	version, err := r.Bytes().AsInt32()
	if err != nil {
		return fmt.Errorf("reading version: %w", err)
	}

	if version != expectedDC6Version {
		return errors.New("unexpected dc6 version")
	}

	if d.Flags, err = r.Bytes().AsUInt32(); err != nil {
		return fmt.Errorf("reading flags: %w", err)
	}

	if d.Encoding, err = r.Bytes().AsUInt32(); err != nil {
		return fmt.Errorf("reading encoding type: %w", err)
	}

	termination, err := r.Next(terminationSize).Bytes().AsBytes()
	if err != nil {
		return fmt.Errorf("reading termination: %w", err)
	}

	copy(d.Termination[:], termination)

	r.Next(bytesPerInt32) // set readed data size to 4 bytes

	directions, err := r.Bytes().AsUInt32()
	if err != nil {
		return fmt.Errorf("reading directions number: %w", err)
	}

	d.Frames.SetNumberOfDirections(int(directions))

	framesPerDirection, err := r.Bytes().AsUInt32()
	if err != nil {
		return fmt.Errorf("error reading a number of frames per direction: %w", err)
	}

	d.Frames.SetFramesPerDirection(int(framesPerDirection))

	return nil
}

func (d *DC6) loadFrames(r *bitstream.Reader) error {
	var err error

	for dir := 0; dir < d.Frames.NumberOfDirections(); dir++ {
		for f := 0; f < d.Frames.FramesPerDirection(); f++ {
			err = d.Frames.Direction(dir).Frame(f).Load(r)
			if err != nil {
				return fmt.Errorf("error loading frame %d at direction %d: %w", f, dir, err)
			}
		}
	}

	return nil
}

// Encode encodes dc6 animation back into byte slice
func (d *DC6) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	// Encode header
	sw.PushInt32(expectedDC6Version)
	sw.PushUint32(d.Flags)
	sw.PushUint32(d.Encoding)

	sw.PushBytes(d.Termination[:]...)

	sw.PushUint32(uint32(d.Frames.NumberOfDirections()))
	sw.PushUint32(uint32(d.Frames.FramesPerDirection()))

	numDirs := d.Frames.NumberOfDirections()
	fpd := d.Frames.FramesPerDirection()

	// encode frames
	framesData := make([][][]byte, numDirs)
	for dir := 0; dir < numDirs; dir++ {
		framesData[dir] = make([][]byte, fpd)
		for f := 0; f < fpd; f++ {
			framesData[dir][f] = d.Frames.Direction(dir).Frame(f).Encode()
		}
	}

	// current position in stream - terrible workaround,
	// but d2datautils.StreamWriter doesn't currently hav any
	// method to get byte position
	currentPosition := 24

	// frames data starts afte a frame pointers section
	currentPosition += numDirs * fpd * bytesPerInt32

	// encode frame pointers
	for dir := 0; dir < numDirs; dir++ {
		for f := 0; f < fpd; f++ {
			sw.PushUint32(uint32(currentPosition))
			currentPosition += len(framesData[dir][f])
		}
	}

	for _, dirData := range framesData {
		for _, frameData := range dirData {
			sw.PushBytes(frameData...)
		}
	}

	return sw.GetBytes()
}

// DecodeFrame decodes the given frame to an indexed color texture
func (d *DC6) DecodeFrame(directionIndex, frameIndex int) []byte {
	frame := d.Frames.Direction(directionIndex).Frame(frameIndex)

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
	clone.Frames = d.Frames.Clone()

	return &clone
}
