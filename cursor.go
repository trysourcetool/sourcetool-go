package sourcetool

type containerType int

const main containerType = iota

type cursor struct {
	rootContainer containerType
	parentPath    []int
	index         int
}

func newCursor(container containerType) *cursor {
	return &cursor{
		rootContainer: container,
		parentPath:    []int{},
		index:         0,
	}
}

func (c *cursor) getDeltaPath() []int {
	path := []int{int(c.rootContainer)}
	path = append(path, c.parentPath...)
	path = append(path, c.index)
	return path
}

func (c *cursor) next() {
	c.index++
}
