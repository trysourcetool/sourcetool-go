package sourcetool

type ContainerType int

const MAIN ContainerType = iota

type DeltaPath struct {
	Container  int
	ParentPath []int
	Index      int
}

type Cursor struct {
	rootContainer ContainerType
	parentPath    []int
	index         int
}

func NewCursor(container ContainerType) *Cursor {
	return &Cursor{
		rootContainer: container,
		parentPath:    []int{},
		index:         0,
	}
}

func (c *Cursor) GetDeltaPath() []int {
	path := []int{int(c.rootContainer)}
	path = append(path, c.parentPath...)
	path = append(path, c.index)
	return path
}

func (c *Cursor) Next() {
	c.index++
}

func (c *Cursor) EnterBlock() *Cursor {
	return &Cursor{
		rootContainer: c.rootContainer,
		parentPath:    append(c.parentPath, c.index),
		index:         0,
	}
}
