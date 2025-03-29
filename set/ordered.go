package set

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
)

// OrderedSet is a set that maintains the order of its elements when returned.
// It leverages an underlying map for efficient membership checks,
// but imposes sorting when elements are retrieved.
//
// Specifically the Members(), String(), and MarshalJSON() methods sort the
// elements.
// To get the unsorted elements, use InsertionOrderMembers().
type OrderedSet[E cmp.Ordered] struct {
	Set[E]
}

// NewOrdered creates a new ordered set with the given values.
//
// Example:
//
//	s := New(3, 1, 2)
func NewOrdered[E cmp.Ordered](vals ...E) OrderedSet[E] {
	return OrderedSet[E]{
		New(vals...),
	}
}

// Contains returns true if the set contains the given value.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	fmt.Println(s.Contains(2))
//	true
func (s OrderedSet[E]) Contains(v E) bool {
	return s.Set.Contains(v)
}

// ContainsAll returns true if the set contains all the given values.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	fmt.Println(s.ContainsAll(2, 3))
//	true
func (s OrderedSet[E]) ContainsAll(vals ...E) bool {
	return s.Set.ContainsAll(vals...)
}

// Add adds the given values to the set.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	s.Add(4, 5)
//	fmt.Println(s.Members())
//	[1 2 3 4 5]
func (s OrderedSet[E]) Add(vals ...E) {
	for _, v := range vals {
		s.Set.Add(v)
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
func (s OrderedSet[E]) AddImmutable(vals ...E) Interface[E] {
	return s.Set.AddImmutable(vals...)
}

// Remove removes the given values from the set.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	s.Remove(2, 3)
//	fmt.Println(s.Members())
//	[1]
func (s OrderedSet[E]) Remove(vals ...E) {
	s.Set.Remove(vals...)
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
func (s OrderedSet[E]) RemoveImmutable(vals ...E) Interface[E] {
	return s.Set.RemoveImmutable(vals...)
}

// Clear removes all values from the set.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	s.Clear()
//	fmt.Println(s.Members())
//	[]
func (s OrderedSet[E]) Clear() {
	s.Set.Clear()
}

// Members returns the members of the set as a slice.
// The result is sorted in ascending order.
//
// Example:
//
//	s := NewSet(3, 1, 3)
//	fmt.Println(s.Members())
//	[1 2 3]
func (s OrderedSet[E]) Members() []E {
	result := s.Set.Members()
	slices.Sort(result)
	return result
}

// InsertionOrderMembers returns the members of the set as a slice in the order
// they were inserted.
//
// Example:
//
//	s := NewSet(2, 3, 1)
//	fmt.Println(s.Members())
//	[1 2 3]
//	fmt.Println(s.InsertionOrderMembers())
//	[2 3 1]
func (s OrderedSet[E]) InsertionOrderMembers() []E {
	return s.Set.Members()
}

// String returns a string representation of the set.
// Order is guaranteed (ascending).
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	fmt.Println(s)
//	[1 2 3]
func (s OrderedSet[E]) String() string {
	result := s.Set.Members()
	slices.Sort(result)
	return fmt.Sprintf("%v", result)
}

// Union returns the union of the set with another set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(2, 3, 4)
//	fmt.Println(s1.Union(s2))
//	[1 2 3 4 5]
func (s OrderedSet[E]) Union(s2 Interface[E]) Interface[E] {
	return s.Set.Union(s2)
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
func (s OrderedSet[E]) Intersection(s2 Interface[E]) Interface[E] {
	return s.Set.Intersection(s2)
}

// Difference returns the difference of the set with another set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(2, 3, 4)
//	fmt.Println(s1.Difference(s2))
//	[1]
func (s OrderedSet[E]) Difference(s2 Interface[E]) Interface[E] {
	return s.Set.Difference(s2)
}

// IsSubsetOf returns true if the set is a subset of another set.
//
// Example:
//
//	s1 := NewSet(1, 2)
//	s2 := NewSet(1, 2, 3)
//	fmt.Println(s1.IsSubsetOf(s2))
//	true
func (s OrderedSet[E]) IsSubsetOf(s2 Interface[E]) bool {
	return s.Set.IsSubsetOf(s2)
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
func (s OrderedSet[E]) IsSupersetOf(s2 Interface[E]) bool {
	return s.Set.IsSupersetOf(s2)
}

// Equal returns true if the set is equal to another set.
//
// Example:
//
//	s1 := NewSet(1, 2, 3)
//	s2 := NewSet(1, 2, 3)
//	fmt.Println(s1.Equal(s2))
//	true
func (s OrderedSet[E]) Equal(s2 Interface[E]) bool {
	return s.Set.Equal(s2)
}

// Len returns the number of elements in the set.
//
// Example:
//
//	s := NewSet(1, 2, 3)
//	fmt.Println(s.Len())
//	3
func (s OrderedSet[E]) Len() int {
	return s.Set.Len()
}

// MarshalJSON implements the json.Marshaler interface to convert the map into
// a JSON object. The values are ordered (ascending).
func (s OrderedSet[E]) MarshalJSON() (out []byte, err error) {
	if s.Len() == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(s.Members())
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json
// object to a Set.
//
// The JSON must start with a list.
func (s OrderedSet[E]) UnmarshalJSON(in []byte) (err error) {
	return s.Set.UnmarshalJSON(in)
}
