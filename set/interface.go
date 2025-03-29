package set

// Interface is the interface used by all the Set types.
type Interface[E comparable] interface {
	// Contains returns true if the set contains the given value.
	Contains(E) bool

	// ContainsAll returns true if the set contains all the given values.
	ContainsAll(...E) bool

	// Add adds the given values to the set.
	Add(...E)

	// AddImmutable adds the given values to the set and returns a new set.
	AddImmutable(...E) Interface[E]

	// Remove removes the given values from the set.
	Remove(...E)

	// RemoveImmutable removes the given values from the set and returns
	// a new set.
	RemoveImmutable(...E) Interface[E]

	// Clear removes all values from the set.
	Clear()

	// Members returns the members of the set as a slice.
	Members() []E

	// String returns a string representation of the set.
	String() string

	// Union returns the union of the set with another set.
	Union(Interface[E]) Interface[E]

	// Intersection returns the intersection of the set with another
	// set.
	Intersection(Interface[E]) Interface[E]

	// Difference returns the difference of the set with another set.
	Difference(Interface[E]) Interface[E]

	// IsSubsetOf returns true if the set is a subset of another set.
	IsSubsetOf(Interface[E]) bool

	// IsSupersetOf returns true if the set is a superset of another
	IsSupersetOf(Interface[E]) bool

	// Equal returns true if the set is equal to another set.
	Equal(Interface[E]) bool

	// Len returns the number of elements in the set.
	Len() int

	// MarshalJSON implements the json.Marshaler interface to convert the map into
	// a JSON object.
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON implements the json.Unmarshaler interface to convert a json
	// object to a Set.
	UnmarshalJSON([]byte) error
}
