package issue22

type Foo struct{}

func (xxx *Foo) String() string {
	return "foo"
}
