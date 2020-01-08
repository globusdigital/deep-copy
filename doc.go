// deep-copy is a tool for generating DeepCopy() functions for a given type.
//
// Given a package directory, and a type name that appears in that package, a
// DeepCopy method will be generated, to create a deep copy of the type value.
// Members of the type will also be copied deeply, recursively.
//
// To specify a pointer receiver for the method, an optional --pointer-receiver
// boolean flag can be specified. The flag will also govern whether the return
// type is a pointer as well.
//
// It might also be desirable to skip deeply copying certain fields. To achieve
// that, field selectors can be specified in the optional comma-separated --skip
// flag.
package main
