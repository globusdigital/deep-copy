package testdata

type Depth1 struct {
	a1 *Depth2
	a2 *Depth2
}

type Depth2 struct {
	b1 *Depth3
}

type Depth3 struct {
	c *Depth4
}

type Depth4 struct {
	d int
}
