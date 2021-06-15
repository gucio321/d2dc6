package frames

// New creates a new frame grid
func New() *Frames {
	return &Frames{
		numberOfDirections: 0,
		framesPerDirection: 0,
		grid:               make([]Direction, 0),
	}
}

// Frames represents a grid of frames [directions][framesPerDirection]
type Frames struct {
	grid []Direction
	numberOfDirections,
	framesPerDirection int
}

// NumberOfDirections returns a number of directions in grid
func (f *Frames) NumberOfDirections() int {
	return f.numberOfDirections
}

// SetNumberOfDirections sets a number of directions
func (f *Frames) SetNumberOfDirections(n int) {
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
func (f *Frames) FramesPerDirection() int {
	return f.framesPerDirection
}

// SetFramesPerDirection sets a number of frames per direction
func (f *Frames) SetFramesPerDirection(n int) {
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
func (f *Frames) Direction(d int) Direction {
	if d > len(f.grid) {
		return nil
	}

	return f.grid[d]
}

// Clone clones frame grid
func (f *Frames) Clone() *Frames {
	clone := &Frames{}
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
