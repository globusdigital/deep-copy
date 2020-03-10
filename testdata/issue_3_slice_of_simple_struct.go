package testdata

type I3WithSlice struct {
	a []I3SimpleStruct
	b string
}

type I3SimpleStruct struct {
	foo string
	bar int
}

type I3WithMap struct {
	a map[I3SimpleStruct]string
	b int
}

type I3WithMapVal struct {
	a map[string]I3SimpleStruct
	b int
}
