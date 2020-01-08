package testdata

type Foo struct {
	Map map[string]*Bar
	ch  chan float32
	baz Baz
}

type Bar struct {
	IntV  int
	Slice []string
}

type Baz struct {
	String        string
	StringPointer *string
}
