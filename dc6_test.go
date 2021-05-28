package d2dc6

import (
	"testing"

	"github.com/gucio321/d2dc6/d2dc6frame"
	"github.com/stretchr/testify/assert"
)

func TestDC6New(t *testing.T) {
	if dc6 := New(); dc6 == nil {
		t.Error("d2dc6.New() method returned nil")
	}
}

func getExampleDC6() *DC6 {
	exampleDC6 := &DC6{
		Version:            6,
		Flags:              1,
		Encoding:           0,
		Termination:        []byte{238, 238, 238, 238},
		Directions:         1,
		FramesPerDirection: 1,
		FramePointers:      []uint32{56},
		Frames: []*d2dc6frame.DC6Frame{
			{
				Flipped:    0,
				Width:      32,
				Height:     26,
				OffsetX:    45,
				OffsetY:    24,
				Unknown:    0,
				NextBlock:  50,
				Length:     10,
				FrameData:  []byte{2, 23, 34, 128, 53, 64, 39, 43, 123, 12},
				Terminator: []byte{2, 8, 5},
			},
		},
	}

	return exampleDC6
}

func TestDC6Unmarshal(t *testing.T) {
	exampleDC6 := getExampleDC6()

	data := exampleDC6.Marshal()

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
