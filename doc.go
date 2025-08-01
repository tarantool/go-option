// Package option provides a type-safe way to represent optional values in Go.
// An Optional[T] can either contain a value of type T (Some) or be empty (None).
//
// This is useful for:
// - Clearly representing nullable fields in structs.
// - Avoiding nil pointer dereferences.
// - Providing explicit intent about optional values.
package option
