package set

import (
	"encoding/json"
	"fmt"
)

// Set is a set that does not maintain the order of its elements when returned.
// It leverages an underlying map for efficient membership checks.
type Set[E comparable] map[E]struct{}

// New creates a new set with the given values.
//
// Example:
//
//	s := New(1, 2, 3)
func New[E comparable](vals ...E) Set[E] {
	s := Set[E]{}
	s.Add(vals...)
	return s
}

// Contains returns true if the set contains the given value.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	fmt.Println(s.Contains(2))
//	true
func (s Set[E]) Contains(v E) bool {
	_, ok := s[v]
	return ok
}

// ContainsAll returns true if the set contains all the given values.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	fmt.Println(s.ContainsAll(2, 3))
//	true
func (s Set[E]) ContainsAll(vals ...E) bool {
	for _, v := range vals {
		if !s.Contains(v) {
			return false
		}
	}
	return true
}

// Add adds the given values to the set.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	s.Add(4, 5)
//	fmt.Println(s.Members())
//	[1 2 3 4 5]
func (s Set[E]) Add(vals ...E) {
	for _, v := range vals {
		s[v] = struct{}{}
	}
}

// AddImmutable adds the given values to the set and returns a new set.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	s2 := s.AddImmutable(4, 5)
//	fmt.Println(s.Members())
//	[1 2 3]
//	fmt.Println(s2.Members())
//	[1 2 3 4 5]
func (s Set[E]) AddImmutable(vals ...E) Interface[E] {
	n := New(s.Members()...)
	n.Add(vals...)
	return n
}

// Remove removes the given values from the set.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	s.Remove(2, 3)
//	fmt.Println(s.Members())
//	[1]
func (s Set[E]) Remove(vals ...E) {
	for _, v := range vals {
		delete(s, v)
	}
}

// RemoveImmutable removes the given values from the set and returns
// a new set.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	s2 := s.RemoveImmutable(2, 3)
//	fmt.Println(s1.Members())
//	[1 2 3]
//	fmt.Println(s2.Members())
//	[1]
func (s Set[E]) RemoveImmutable(vals ...E) Interface[E] {
	n := New(s.Members()...)
	n.Remove(vals...)
	return n
}

// Clear removes all values from the set.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	s.Clear()
//	fmt.Println(s.Members())
//	[]
func (s Set[E]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Members returns the members of the set as a slice.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	fmt.Println(s.Members())
//	[1 2 3]
func (s Set[E]) Members() []E {
	result := make([]E, 0, len(s))
	for v := range s {
		result = append(result, v)
	}
	return result
}

// String returns a string representation of the set. Order is not
// guaranteed.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	fmt.Println(s)
//	[1 2 3] *or* [3 2 1] *or* [2 1 3] *or* ...
func (s Set[E]) String() string {
	return fmt.Sprintf("%v", s.Members())
}

// Union returns the union of the set with another set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(2, 3, 4)
//	fmt.Println(s1.Union(s2))
//	[1 2 3 4 5]
func (s Set[E]) Union(s2 Interface[E]) Interface[E] {
	result := New(s.Members()...)
	result.Add(s2.Members()...)
	return result
}

// Intersection returns the intersection of the set with another
// set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(2, 3, 4)
//	fmt.Println(s1.Intersection(s2))
//	[3]
func (s Set[E]) Intersection(s2 Interface[E]) Interface[E] {
	result := New[E]()
	for _, v := range s.Members() {
		if s2.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// Difference returns the difference of the set with another set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(2, 3, 4)
//	fmt.Println(s1.Difference(s2))
//	[1]
func (s Set[E]) Difference(s2 Interface[E]) Interface[E] {
	result := New(s.Members()...)
	for _, v := range s2.Members() {
		delete(result, v)
	}
	return result
}

// IsSubsetOf returns true if the set is a subset of another set.
//
// Example:
//
//	s1 := NewSet(1, 2)
//	s2 := NewSet(1, 2, 3)
//	fmt.Println(s1.IsSubsetOf(s2))
//	true
func (s Set[E]) IsSubsetOf(s2 Interface[E]) bool {
	for _, v := range s.Members() {
		if !s2.Contains(v) {
			return false
		}
	}
	return true
}

// IsSupersetOf returns true if the set is a superset of another
// set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(1, 2)
//	fmt.Println(s1.IsSupersetOf(s2))
//	true
func (s Set[E]) IsSupersetOf(s2 Interface[E]) bool {
	for _, v := range s2.Members() {
		if !s.Contains(v) {
			return false
		}
	}
	return true
}

// Equal returns true if the set is equal to another set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(1, 2, 3)
//	fmt.Println(s1.Equal(s2))
//	true
func (s Set[E]) Equal(s2 Interface[E]) bool {
	return s.IsSubsetOf(s2) && s.IsSupersetOf(s2)
}

// Len returns the number of elements in the set.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	fmt.Println(s.Len())
//	3
func (s Set[E]) Len() int {
	return len(s)
}

// MarshalJSON implements the json.Marshaler interface to convert the map into
// a JSON object.
func (s Set[E]) MarshalJSON() (out []byte, err error) {
	if s.Len() == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(s.Members())
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json
// object to a Set.
//
// The JSON must start with a list.
func (s Set[E]) UnmarshalJSON(in []byte) (err error) {
	var v []E

	if err = json.Unmarshal(in, &v); err == nil {
		s.Add(v...)
	}
	return
}
