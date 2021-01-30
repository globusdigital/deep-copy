package somepkg

import (
	"github.com/globusdigital/deep-copy/testdata/pointer_that_implements_deepcopy/anotherpkg"
)

type SomeStruct struct {
	AnotherStruct *anotherpkg.AnotherStruct
}
