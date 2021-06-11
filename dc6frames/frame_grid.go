package dc6frames

type FrameGrid []Direction

func (f FrameGrid) Direction(d int) Direction {
	if d > len(f) {
		return nil
	}

	return f[d]
}

type Direction []*Frame

func (d Direction) Frame(f int) *Frame {
	if f > len(d) {
		return nil
	}

	return d[f]
}
