package testdata

type SomeStruct struct {
	mapSlice map[string][]string
}

type SomeStruct2 struct {
	mapStruct map[string]SomeStruct
}
