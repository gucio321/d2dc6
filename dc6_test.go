package d2dc6

import (
	"testing"

	"github.com/gucio321/d2dc6/dc6frames"
	"github.com/stretchr/testify/assert"
)

func TestDC6New(t *testing.T) {
	if dc6 := New(); dc6 == nil {
		t.Error("d2dc6.New() method returned nil")
	}
}

func getExampleDC6() *DC6 {
	exampleDC6 := &DC6{
		Flags:              1,
		Encoding:           0,
		Termination:        []byte{238, 238, 238, 238},
		Directions:         0,
		FramesPerDirection: 0,
		// FramePointers:      []uint32{56, 100, 140, 180},
		FramePointers: []uint32{},
		Frames:        dc6frames.New(),
	}

	exampleDC6.Frames.SetNumberOfDirections(int(exampleDC6.Directions))
	exampleDC6.Frames.SetFramesPerDirection(int(exampleDC6.FramesPerDirection))
	/*
			grid: {
				{
					Flipped:    0,
					Width:      32,
					Height:     26,
					OffsetX:    45,
					OffsetY:    24,
					Unknown:    0,
					NextBlock:  50,
					FrameData:  []byte{2, 23, 34, 128, 53, 64, 39, 43, 123, 12},
					Terminator: []byte{2, 8, 5},
				},
				{
					Flipped:    0,
					Width:      62,
					Height:     36,
					OffsetX:    15,
					OffsetY:    28,
					Unknown:    0,
					NextBlock:  35,
					FrameData:  []byte{9, 33, 89, 148, 64, 64, 49, 81, 221, 19},
					Terminator: []byte{3, 7, 5},
				},
			},
			{
				{
					Flipped:    0,
					Width:      62,
					Height:     36,
					OffsetX:    15,
					OffsetY:    28,
					Unknown:    0,
					NextBlock:  35,
					FrameData:  []byte{9, 33, 89, 148, 64, 64, 49, 81, 121, 19},
					Terminator: []byte{3, 7, 5},
				},
				{
					Flipped:    0,
					Width:      32,
					Height:     26,
					OffsetX:    45,
					OffsetY:    24,
					Unknown:    0,
					NextBlock:  50,
					FrameData:  []byte{2, 23, 34, 128, 53, 64, 39, 43, 123, 12},
					Terminator: []byte{2, 8, 5},
				},
			},
		},
	*/

	return exampleDC6
}

func TestDC6Unmarshal(t *testing.T) {
	exampleDC6 := getExampleDC6()

	data := exampleDC6.Encode()

	extractedDC6, err := Load(data)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, exampleDC6, extractedDC6, "encoded and decoded dc6 isn't equal")
}

func TestDC6Clone(t *testing.T) {
	exampleDC6 := getExampleDC6()
	clonedDC6 := exampleDC6.Clone()

	assert.Equal(t, exampleDC6, clonedDC6, "cloned dc6 isn't equal to base")
}
