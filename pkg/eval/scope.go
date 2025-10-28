package eval

type Scope interface {
	// GetArgument gets the value argument
	// in specified place.
	//
	// Returns an error if the given argument doesn't match the type,
	// place is higher that the arguments' length or lower than zero.
	GetArgument(place int, tType string) (any, error)

	// Set changes the argument's value to v
	// in specified order.
	//
	// Returns an error if place is higher that the arguments' length
	// or lower than zero.
	SetArgument(place int, tType string, v any) (any, error)

	// Get gets the global symbol' value, if it's variable or constant variable.
	//
	// Returns an error if it doesn't exist, or
	// it's not allowed kind.q
	Get(name string) (any, error)
}

// BuiltinFunction is a builtin function.
type BuiltinFunction func(s Scope) (any, error)
