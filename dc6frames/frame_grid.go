package dc6frames

// New creates a new frame grid
func New() *FrameGrid {
	return &FrameGrid{
		numberOfDirections: 0,
		framesPerDirection: 0,
		grid:               make([]Direction, 0),
	}
}

type FrameGrid struct {
	grid []Direction
	numberOfDirections,
	framesPerDirection int
}

func (f *FrameGrid) NumberOfDirections() int {
	return f.numberOfDirections
}

func (f *FrameGrid) SetNumberOfDirections(n int) {
	if n == f.numberOfDirections {
		return
	}

	if n > f.numberOfDirections {
		for i := 0; i < n-f.framesPerDirection; i++ {
			f.grid = append(f.grid, make(Direction, f.framesPerDirection))
		}

		f.numberOfDirections = n
		return
	}

	f.grid = f.grid[:n]

	f.numberOfDirections = n
}

func (f *FrameGrid) FramesPerDirection() int {
	return f.framesPerDirection
}

func (f *FrameGrid) SetFramesPerDirection(n int) {
	if n == f.framesPerDirection {
		return
	}

	for i := range f.grid {
		if l := len(f.grid[i]); n > l {
			for j := 0; j < n-l; j++ {
				f.grid[i] = append(f.grid[i], newFrame())
			}
		}
	}

	f.framesPerDirection = n
}

func (f *FrameGrid) Direction(d int) Direction {
	if d > len(f.grid) {
		return nil
	}

	return f.grid[d]
}

func (f *FrameGrid) Clone() *FrameGrid {
	clone := &FrameGrid{}
	clone.SetNumberOfDirections(f.numberOfDirections)
	clone.SetFramesPerDirection(f.framesPerDirection)

	for dir := range clone.grid {
		for frame := range clone.grid[dir] {
			*clone.grid[dir][frame] = *f.grid[dir][frame]
		}
	}

	return clone
}

type Direction []*Frame

func (d Direction) Frame(f int) *Frame {
	if f > len(d) {
		return nil
	}

	return d[f]
}
