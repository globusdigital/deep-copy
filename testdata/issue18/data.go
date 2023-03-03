package issue18

type Foo struct {
	b1 Boo
	b2 *Boo
}

type Boo struct {
	s string
}
