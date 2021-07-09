package testdata

type ParentHasChildValue struct {
	c Child
}

type ParentHasChildPointer struct {
	c *Child
}

type Child struct {
	s string
}
