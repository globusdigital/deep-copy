// deep-copy is a tool for generating DeepCopy() functions for a given type.
//
// Given a package directory, and a type name that appears in that package, a
// DeepCopy method will be generated, to create a deep copy of the type value.
// Members of the type will also be copied deeply, recursively. If a member T
// of the type has a method "DeepCopy() [*]T", that method will be reused.
// Multiple types can be specified for the given package, by adding more --type
// parameters.
//
// To specify a pointer receiver for the method, an optional --pointer-receiver
// boolean flag can be specified. The flag will also govern whether the return
// type is a pointer as well.
//
// It might also be desirable to skip deeply copying certain fields, slice
// members, or map members. To achieve that, selectors can be specified in the
// optional comma-separated --skip flag. Multiple --skip flags can be
// specified, to match the number of --type flags.
package main
