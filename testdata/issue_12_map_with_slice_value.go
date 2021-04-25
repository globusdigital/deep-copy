package testdata

type I12StructWithMapOfSlices struct {
	Sc1 map[string][]I12StructWithSlices
}

type I12StructWithSlices struct {
	Name []string
}
