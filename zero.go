package option

// zero is a helper function that returns the zero value of a type.
// It simplifies returning abstract zero values in template-generated code.
func zero[T any]() T {
	var zero T

	return zero
}
