<a href="https://github.com/globusdigital/deep-copy/actions?query=workflow%3ACI">
    <img src="https://github.com/globusdigital/deep-copy/workflows/CI/badge.svg" alt="GitHub Actions: CI">
</a>
<a href='https://coveralls.io/github/globusdigital/deep-copy'>
    <img src='https://coveralls.io/repos/github/globusdigital/deep-copy/badge.svg' alt='Coverage Status' />
</a>
<a href="https://goreportcard.com/report/github.com/globusdigital/deep-copy">
    <img src="https://goreportcard.com/badge/github.com/globusdigital/deep-copy" alt="GoReportCard">
</a>

# deep-copy

deep-copy is a tool for generating DeepCopy() functions for a given type.

Given a package directory, and a type name that appears in that package, a
`DeepCopy` method will be generated, to create a deep copy of the type value.
Members of the type will also be copied deeply, recursively. If a member `T` of
the type has a method `DeepCopy() [*]T`, that method will be reused. Multiple
types can be specified for the given package, by adding more `--type`
parameters.

To specify a pointer receiver for the method, an optional `--pointer-receiver`
boolean flag can be specified. The flag will also govern whether the return
type is a pointer as well.

It might also be desirable to skip deeply copying certain fields, slice
members, or map members. To achieve that, selectors can be specified in the
optional comma-separated `--skip` flag. Multiple `--skip` flags can be
specified, to match the number of `--type` flags. For example, given the
following type:

```go
type Foo struct {
     J *int
     B Bar
}

type Bar struct {
    I *int
}
```

Leaving the 'B' field as a shallow copy can be achieved by specifying `--skip
B`. To skip deeply copying the inner 'I' field, one can specify `--skip B.I`.
Slice and Map members can also be skipped, by adding `[i]` and `[k]`
respectively.

To specify a max depth of deep copying, use `--maxdepth` option. It stops
deep copying at a given depth, with a warning message spotting a place
the deep copying has been stopped. It might especially be useful when
one or more structs have circular references.

To change a method name of deep copying, use `--method` option.

To match receiver names with existing ones, use `--reuse-receiver` option.

## Usage

Pass either path to the folder containing the types or the module name:

```bash
deep-copy <flags> /path/to/package/containing/type
deep-copy <flags> github.com/globusdigital/deep-copy
deep-copy <flags> github.com/globusdigital/deep-copy/some/sub/packages
```
Here is the full set of supported flags:

```bash
deep-copy \ 
  [-o /output/path.go] \
  [--method DeepCopy] \
  [--pointer-receiver] \
  [--skip Selector1,Selector.Two --skip Selector2[i], Selector.Three[k]] \
  [--type Type1 --type Type2] \
  [--reuse-receiver] \
  /path/to/package/containing/type
```

## Example

Given the following types:

```go
package pkg

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
```

Running `deep-copy --type Foo ./path/to/pkg` will generate:

```go
// Code generated by deep-copy --type Foo ./path/to/pkg; DO NOT EDIT.

package pkg

// DeepCopy generates a deep copy of Foo
func (o Foo) DeepCopy() Foo {
	var cp Foo
	cp = o
	if o.Map != nil {
		cp.Map = make(map[string]*Bar, len(o.Map))
		for k, v := range o.Map {
			var cpv *Bar
			if v != nil {
				cpv = new(Bar)
				*cpv = *v
				if v.Slice != nil {
					cpv.Slice = make([]string, len(v.Slice))
					copy(cpv.Slice, v.Slice)
				}
			}
			cp.Map[k] = cpv
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
}
```
