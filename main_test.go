package main

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name     string
		types    typesVal
		path     string
		pointer  bool
		skips    skipsVal
		maxdepth int
		method   string
		want     []byte
	}{
		{name: "foo", types: typesVal{"Foo"}, path: "./testdata", want: []byte(FooFile)},
		{name: "foo - pointer", types: typesVal{"Foo"}, pointer: true, path: "./testdata", want: []byte(FooPointerFile)},
		{name: "foo - pointer, skip slice", types: typesVal{"Foo"}, pointer: true, skips: skipsVal{{"Slice": struct{}{}}}, path: "./testdata", want: []byte(FooPointerSkipSliceFile)},
		{name: "foo, skip map member", types: typesVal{"Foo"}, skips: skipsVal{{"Map[k]": struct{}{}}}, path: "./testdata", want: []byte(FooSkipMapFile)},
		{name: "alpha - with DeepCopy method", types: typesVal{"Alpha"}, path: "./testdata", want: []byte(AlphaPointer)},
		{name: "slicepointer, skip slice member", types: typesVal{"SlicePointer"}, skips: skipsVal{{"[i]": struct{}{}}}, path: "./testdata", want: []byte(SlicePointer)},
		{name: "foo, alpha, skips", types: typesVal{"Foo", "Alpha"}, skips: skipsVal{{"Map[k]": struct{}{}, "ch": struct{}{}}, {"D": struct{}{}, "E": struct{}{}}}, path: "./testdata", want: []byte(FooAlphaSkips)},
		{name: "foo, method=Clone", types: typesVal{"Foo"}, path: "./testdata", method: "Clone", want: []byte(FooCloneFile)},
		{name: "issue 3, struct with slice of simple structs", types: typesVal{"I3WithSlice"}, pointer: true, path: "./testdata", want: []byte(Issue3SliceSimpleStruct)},
		{name: "issue 3, struct with map of simple struct keys", types: typesVal{"I3WithMap"}, pointer: true, path: "./testdata", want: []byte(Issue3MapSimpleStructKey)},
		{name: "issue 3, struct with map of simple struct values", types: typesVal{"I3WithMapVal"}, path: "./testdata", want: []byte(Issue3MapSimpleStructVal)},
		{name: "issue 7, shadowed map vars", types: typesVal{"SomeStruct2"}, path: "./testdata", want: []byte(Issue7ShadowedMapVars)},
		{name: "issue 7, shadowed map vars 2", types: typesVal{"SomeStruct", "SomeStruct2"}, path: "./testdata", want: []byte(Issue7ShadowedMapVars2)},
		{name: "pointer that implements DeepCopy", types: typesVal{"SomeStruct"}, path: "./testdata/pointer_that_implements_deepcopy/somepkg", want: []byte(PointerThatImplementsDeepcopy)},
		{name: "issue 10, slice with element that contains pointer and value", types: typesVal{"StructCH"}, path: "./testdata", want: []byte(Issue10StructCH)},
		{name: "issue 12, nested slices", types: typesVal{"I12NestedSlices"}, path: "./testdata", want: []byte(Issue12NestedSlices)},
		{name: "issue 12, map with slice value", types: typesVal{"I12StructWithMapOfSlices"}, path: "./testdata", want: []byte(Issue12MapWithSliceValues)},
		{name: "issue 15, parent has child value, value receiver", types: typesVal{"ParentHasChildValue", "Child"}, path: "./testdata", want: []byte(I15ParentHasChildValueValueRecv)},
		{name: "issue 15, parent has child pointer, value receiver", types: typesVal{"ParentHasChildPointer", "Child"}, path: "./testdata", want: []byte(I15ParentHasChildPointerValueRecv)},
		{name: "issue 15, parent has child value, pointer receiver", pointer: true, types: typesVal{"ParentHasChildValue", "Child"}, path: "./testdata", want: []byte(I15ParentHasChildValuePointerRecv)},
		{name: "issue 15, parent has child pointer, pointer receiver", pointer: true, types: typesVal{"ParentHasChildPointer", "Child"}, path: "./testdata", want: []byte(I15ParentHasChildPointerPointerRecv)},
		{name: "issue 17, with maxdepth", types: typesVal{"Depth1"}, pointer: true, maxdepth: 2, path: "./testdata", want: []byte(Issue17MaxDepth)},
		{name: "alias import", types: typesVal{"Data"}, path: "./testdata/import_alias", want: []byte(AliasImport)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := "DeepCopy"
			if tt.method != "" {
				method = tt.method
			}
			a := &app{
				isPtrRecv: tt.pointer,
				maxDepth:  tt.maxdepth,
				method:    method,
			}
			got, err := a.run(tt.path, tt.types, tt.skips)
			if err != nil {
				t.Fatal(err)
			}
			got = normalizeComment(got)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("generateFile() diff = %s", diff)
			}
		})
	}
}

var re = regexp.MustCompile(`generated by .*deep-copy.*; DO NOT EDIT.`)

func normalizeComment(in []byte) []byte {
	return re.ReplaceAll(bytes.TrimSpace(in), []byte("generated by deep-copy; DO NOT EDIT."))
}

const (
	FooFile = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of Foo
func (o Foo) DeepCopy() Foo {
	var cp Foo = o
	if o.Map != nil {
		cp.Map = make(map[string]*Bar, len(o.Map))
		for k2, v2 := range o.Map {
			var cp_Map_v2 *Bar
			cp_Map_v2 = v2
			if v2 != nil {
				cp_Map_v2 = new(Bar)
				*cp_Map_v2 = *v2
				if v2.Slice != nil {
					cp_Map_v2.Slice = make([]string, len(v2.Slice))
					copy(cp_Map_v2.Slice, v2.Slice)
				}
			}
			cp.Map[k2] = cp_Map_v2
		}
	}
	if o.ch != nil {
		cp.ch = make(chan float32, cap(o.ch))
	}
	if o.baz.StringPointer != nil {
		cp.baz.StringPointer = new(string)
		*cp.baz.StringPointer = *o.baz.StringPointer
	}
	return cp
}`
	FooPointerFile = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of *Foo
func (o *Foo) DeepCopy() *Foo {
	var cp Foo = *o
	if o.Map != nil {
		cp.Map = make(map[string]*Bar, len(o.Map))
		for k2, v2 := range o.Map {
			var cp_Map_v2 *Bar
			cp_Map_v2 = v2
			if v2 != nil {
				cp_Map_v2 = new(Bar)
				*cp_Map_v2 = *v2
				if v2.Slice != nil {
					cp_Map_v2.Slice = make([]string, len(v2.Slice))
					copy(cp_Map_v2.Slice, v2.Slice)
				}
			}
			cp.Map[k2] = cp_Map_v2
		}
	}
	if o.ch != nil {
		cp.ch = make(chan float32, cap(o.ch))
	}
	if o.baz.StringPointer != nil {
		cp.baz.StringPointer = new(string)
		*cp.baz.StringPointer = *o.baz.StringPointer
	}
	return &cp
}`
	FooPointerSkipSliceFile = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of *Foo
func (o *Foo) DeepCopy() *Foo {
	var cp Foo = *o
	if o.Map != nil {
		cp.Map = make(map[string]*Bar, len(o.Map))
		for k2, v2 := range o.Map {
			var cp_Map_v2 *Bar
			cp_Map_v2 = v2
			if v2 != nil {
				cp_Map_v2 = new(Bar)
				*cp_Map_v2 = *v2
			}
			cp.Map[k2] = cp_Map_v2
		}
	}
	if o.ch != nil {
		cp.ch = make(chan float32, cap(o.ch))
	}
	if o.baz.StringPointer != nil {
		cp.baz.StringPointer = new(string)
		*cp.baz.StringPointer = *o.baz.StringPointer
	}
	return &cp
}`
	FooSkipMapFile = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of Foo
func (o Foo) DeepCopy() Foo {
	var cp Foo = o
	if o.Map != nil {
		cp.Map = make(map[string]*Bar, len(o.Map))
		for k2, v2 := range o.Map {
			cp.Map[k2] = v2
		}
	}
	if o.ch != nil {
		cp.ch = make(chan float32, cap(o.ch))
	}
	if o.baz.StringPointer != nil {
		cp.baz.StringPointer = new(string)
		*cp.baz.StringPointer = *o.baz.StringPointer
	}
	return cp
}`
	AlphaPointer = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of Alpha
func (o Alpha) DeepCopy() Alpha {
	var cp Alpha = o
	if o.B != nil {
		cp.B = o.B.DeepCopy()
	}
	cp.G = o.G.DeepCopy()
	if o.D != nil {
		retV := o.D.DeepCopy()
		cp.D = &retV
	}
	{
		retV := o.E.DeepCopy()
		cp.E = *retV
	}
	return cp
}`
	SlicePointer = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of SlicePointer
func (o SlicePointer) DeepCopy() SlicePointer {
	var cp SlicePointer = o
	if o != nil {
		cp = make([]*int, len(o))
		copy(cp, o)
	}
	return cp
}`
	FooAlphaSkips = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of Foo
func (o Foo) DeepCopy() Foo {
	var cp Foo = o
	if o.Map != nil {
		cp.Map = make(map[string]*Bar, len(o.Map))
		for k2, v2 := range o.Map {
			cp.Map[k2] = v2
		}
	}
	if o.baz.StringPointer != nil {
		cp.baz.StringPointer = new(string)
		*cp.baz.StringPointer = *o.baz.StringPointer
	}
	return cp
}

// DeepCopy generates a deep copy of Alpha
func (o Alpha) DeepCopy() Alpha {
	var cp Alpha = o
	if o.B != nil {
		cp.B = o.B.DeepCopy()
	}
	cp.G = o.G.DeepCopy()
	return cp
}`

	FooCloneFile = `// generated by deep-copy; DO NOT EDIT.

package testdata

// Clone generates a deep copy of Foo
func (o Foo) Clone() Foo {
	var cp Foo = o
	if o.Map != nil {
		cp.Map = make(map[string]*Bar, len(o.Map))
		for k2, v2 := range o.Map {
			var cp_Map_v2 *Bar
			if v2 != nil {
				cp_Map_v2 = new(Bar)
				*cp_Map_v2 = *v2
				if v2.Slice != nil {
					cp_Map_v2.Slice = make([]string, len(v2.Slice))
					copy(cp_Map_v2.Slice, v2.Slice)
				}
			}
			cp.Map[k2] = cp_Map_v2
		}
	}
	if o.ch != nil {
		cp.ch = make(chan float32, cap(o.ch))
	}
	if o.baz.StringPointer != nil {
		cp.baz.StringPointer = new(string)
		*cp.baz.StringPointer = *o.baz.StringPointer
	}
	return cp
}`

	Issue3SliceSimpleStruct = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of *I3WithSlice
func (o *I3WithSlice) DeepCopy() *I3WithSlice {
	var cp I3WithSlice = *o
	if o.a != nil {
		cp.a = make([]I3SimpleStruct, len(o.a))
		copy(cp.a, o.a)
	}
	return &cp
}`
	Issue3MapSimpleStructKey = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of *I3WithMap
func (o *I3WithMap) DeepCopy() *I3WithMap {
	var cp I3WithMap = *o
	if o.a != nil {
		cp.a = make(map[I3SimpleStruct]string, len(o.a))
		for k2, v2 := range o.a {
			cp.a[k2] = v2
		}
	}
	return &cp
}`
	Issue3MapSimpleStructVal = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of I3WithMapVal
func (o I3WithMapVal) DeepCopy() I3WithMapVal {
	var cp I3WithMapVal = o
	if o.a != nil {
		cp.a = make(map[string]I3SimpleStruct, len(o.a))
		for k2, v2 := range o.a {
			cp.a[k2] = v2
		}
	}
	return cp
}`

	Issue7ShadowedMapVars = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of SomeStruct2
func (o SomeStruct2) DeepCopy() SomeStruct2 {
	var cp SomeStruct2 = o
	if o.mapStruct != nil {
		cp.mapStruct = make(map[string]SomeStruct, len(o.mapStruct))
		for k2, v2 := range o.mapStruct {
			var cp_mapStruct_v2 SomeStruct
			cp_mapStruct_v2 = v2
			if v2.mapSlice != nil {
				cp_mapStruct_v2.mapSlice = make(map[string][]string, len(v2.mapSlice))
				for k4, v4 := range v2.mapSlice {
					var cp_mapStruct_v2_mapSlice_v4 []string
					cp_mapStruct_v2_mapSlice_v4 = v4
					if v4 != nil {
						cp_mapStruct_v2_mapSlice_v4 = make([]string, len(v4))
						copy(cp_mapStruct_v2_mapSlice_v4, v4)
					}
					cp_mapStruct_v2.mapSlice[k4] = cp_mapStruct_v2_mapSlice_v4
				}
			}
			cp.mapStruct[k2] = cp_mapStruct_v2
		}
	}
	return cp
}`

	Issue7ShadowedMapVars2 = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of SomeStruct
func (o SomeStruct) DeepCopy() SomeStruct {
	var cp SomeStruct = o
	if o.mapSlice != nil {
		cp.mapSlice = make(map[string][]string, len(o.mapSlice))
		for k2, v2 := range o.mapSlice {
			var cp_mapSlice_v2 []string
			cp_mapSlice_v2 = v2
			if v2 != nil {
				cp_mapSlice_v2 = make([]string, len(v2))
				copy(cp_mapSlice_v2, v2)
			}
			cp.mapSlice[k2] = cp_mapSlice_v2
		}
	}
	return cp
}

// DeepCopy generates a deep copy of SomeStruct2
func (o SomeStruct2) DeepCopy() SomeStruct2 {
	var cp SomeStruct2 = o
	if o.mapStruct != nil {
		cp.mapStruct = make(map[string]SomeStruct, len(o.mapStruct))
		for k2, v2 := range o.mapStruct {
			var cp_mapStruct_v2 SomeStruct
			cp_mapStruct_v2 = v2
			cp_mapStruct_v2 = v2.DeepCopy()
			cp.mapStruct[k2] = cp_mapStruct_v2
		}
	}
	return cp
}`

	Issue10StructCH = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of StructCH
func (o StructCH) DeepCopy() StructCH {
	var cp StructCH = o
	if o.Nested != nil {
		cp.Nested = make([]StructNested, len(o.Nested))
		copy(cp.Nested, o.Nested)
		for i2 := range o.Nested {
			if o.Nested[i2].B != nil {
				cp.Nested[i2].B = new(int)
				*cp.Nested[i2].B = *o.Nested[i2].B
			}
		}
	}
	return cp
}`

	PointerThatImplementsDeepcopy = `// generated by deep-copy; DO NOT EDIT.

package somepkg

// DeepCopy generates a deep copy of SomeStruct
func (o SomeStruct) DeepCopy() SomeStruct {
	var cp SomeStruct = o
	if o.AnotherStruct != nil {
		cp.AnotherStruct = o.AnotherStruct.DeepCopy()
	}
	return cp
}`

	Issue12NestedSlices = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of I12NestedSlices
func (o I12NestedSlices) DeepCopy() I12NestedSlices {
	var cp I12NestedSlices = o
	if o.Slices != nil {
		cp.Slices = make([][][]int, len(o.Slices))
		copy(cp.Slices, o.Slices)
		for i2 := range o.Slices {
			if o.Slices[i2] != nil {
				cp.Slices[i2] = make([][]int, len(o.Slices[i2]))
				copy(cp.Slices[i2], o.Slices[i2])
				for i3 := range o.Slices[i2] {
					if o.Slices[i2][i3] != nil {
						cp.Slices[i2][i3] = make([]int, len(o.Slices[i2][i3]))
						copy(cp.Slices[i2][i3], o.Slices[i2][i3])
					}
				}
			}
		}
	}
	return cp
}`

	Issue12MapWithSliceValues = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of I12StructWithMapOfSlices
func (o I12StructWithMapOfSlices) DeepCopy() I12StructWithMapOfSlices {
	var cp I12StructWithMapOfSlices = o
	if o.Sc1 != nil {
		cp.Sc1 = make(map[string][]I12StructWithSlices, len(o.Sc1))
		for k2, v2 := range o.Sc1 {
			var cp_Sc1_v2 []I12StructWithSlices
			cp_Sc1_v2 = v2
			if v2 != nil {
				cp_Sc1_v2 = make([]I12StructWithSlices, len(v2))
				copy(cp_Sc1_v2, v2)
				for i3 := range v2 {
					if v2[i3].Name != nil {
						cp_Sc1_v2[i3].Name = make([]string, len(v2[i3].Name))
						copy(cp_Sc1_v2[i3].Name, v2[i3].Name)
					}
				}
			}
			cp.Sc1[k2] = cp_Sc1_v2
		}
	}
	return cp
}`

	I15ParentHasChildValueValueRecv = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of ParentHasChildValue
func (o ParentHasChildValue) DeepCopy() ParentHasChildValue {
	var cp ParentHasChildValue = o
	cp.c = o.c.DeepCopy()
	return cp
}

// DeepCopy generates a deep copy of Child
func (o Child) DeepCopy() Child {
	var cp Child = o
	return cp
}`

	I15ParentHasChildPointerValueRecv = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of ParentHasChildPointer
func (o ParentHasChildPointer) DeepCopy() ParentHasChildPointer {
	var cp ParentHasChildPointer = o
	if o.c != nil {
		retV := o.c.DeepCopy()
		cp.c = &retV
	}
	return cp
}

// DeepCopy generates a deep copy of Child
func (o Child) DeepCopy() Child {
	var cp Child = o
	return cp
}`

	I15ParentHasChildValuePointerRecv = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of *ParentHasChildValue
func (o *ParentHasChildValue) DeepCopy() *ParentHasChildValue {
	var cp ParentHasChildValue = *o
	{
		retV := o.c.DeepCopy()
		cp.c = *retV
	}
	return &cp
}

// DeepCopy generates a deep copy of *Child
func (o *Child) DeepCopy() *Child {
	var cp Child = *o
	return &cp
}`

	I15ParentHasChildPointerPointerRecv = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of *ParentHasChildPointer
func (o *ParentHasChildPointer) DeepCopy() *ParentHasChildPointer {
	var cp ParentHasChildPointer = *o
	if o.c != nil {
		cp.c = o.c.DeepCopy()
	}
	return &cp
}

// DeepCopy generates a deep copy of *Child
func (o *Child) DeepCopy() *Child {
	var cp Child = *o
	return &cp
}`

	Issue17MaxDepth = `// generated by deep-copy; DO NOT EDIT.

package testdata

// DeepCopy generates a deep copy of *Depth1
func (o *Depth1) DeepCopy() *Depth1 {
	var cp Depth1 = *o
	if o.a1 != nil {
		cp.a1 = new(Depth2)
		*cp.a1 = *o.a1
	}
	if o.a2 != nil {
		cp.a2 = new(Depth2)
		*cp.a2 = *o.a2
	}
	return &cp
}`

	AliasImport = `// generated by deep-copy; DO NOT EDIT.

package import_alias

import (
	github_com_globusdigital_deep_copy_testdata_import_alias_another_item "github.com/globusdigital/deep-copy/testdata/import_alias/another/item"
	"github.com/globusdigital/deep-copy/testdata/import_alias/item"
)

// DeepCopy generates a deep copy of Data
func (o Data) DeepCopy() Data {
	var cp Data = o
	if o.Items != nil {
		cp.Items = make([]item.Item, len(o.Items))
		copy(cp.Items, o.Items)
	}
	if o.AnotherItems != nil {
		cp.AnotherItems = make([]github_com_globusdigital_deep_copy_testdata_import_alias_another_item.Item, len(o.AnotherItems))
		copy(cp.AnotherItems, o.AnotherItems)
	}
	return cp
}`
)
