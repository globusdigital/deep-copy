package main

import (
	"log"
	"reflect"

	. "github.com/globusdigital/deep-copy/testdata"
)

func main() {
	f := Foo{Map: map[string]*Bar{
		"key1": {Slice: []string{"s1", "s2"}},
	}}

	cp := f.DeepCopy()

	if !reflect.DeepEqual(f, cp) {
		log.Fatalf("source and sink differ")
	}

	cp.Map["key1"] = &Bar{IntV: 42}

	switch {
	case f.Map["key1"] == cp.Map["key1"]:
		log.Fatalf("key1 values match")
	}
}
