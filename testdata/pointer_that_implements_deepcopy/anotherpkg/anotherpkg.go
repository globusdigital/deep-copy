package anotherpkg

type AnotherStruct struct {
	Field int
}

func (s *AnotherStruct) DeepCopy() *AnotherStruct {
	return &AnotherStruct{
		Field: s.Field,
	}
}
