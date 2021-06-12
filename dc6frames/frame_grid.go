package dc6frames

// New creates a new frame grid
func New() *FrameGrid {
	return &FrameGrid{
		numberOfDirections: 0,
		framesPerDirection: 0,
		grid:               make([]Direction, 0),
	}
}

// FrameGrid represents a grid of frames [directions][framesPerDirection]
type FrameGrid struct {
	grid []Direction
	numberOfDirections,
	framesPerDirection int
}

// NumberOfDirections returns a number of directions in grid
func (f *FrameGrid) NumberOfDirections() int {
	return f.numberOfDirections
}

// SetNumberOfDirections sets a number of directions
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

// FramesPerDirection returns a number of frames per each direction
func (f *FrameGrid) FramesPerDirection() int {
	return f.framesPerDirection
}

// SetFramesPerDirection sets a number of frames per direction
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

// Direction returns a specified direction
func (f *FrameGrid) Direction(d int) Direction {
	if d > len(f.grid) {
		return nil
	}

	return f.grid[d]
}

// Clone clones frame grid
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

// Direction represents a frame set
type Direction []*Frame

// Frame returns a specified frame, if f > FramesPerDirection, returns nil
func (d Direction) Frame(f int) *Frame {
	if f > len(d) {
		return nil
	}

	return d[f]
}
