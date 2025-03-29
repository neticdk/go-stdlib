package set

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrderedSet(t *testing.T) {
	s := NewOrdered(3, 1, 2)
	assert.Equal(t, 3, len(s.Set))
	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
	assert.True(t, s.Contains(3))
}

func TestOrderedSetMembers(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Three elements", []int{3, 1, 2}, []int{1, 2, 3}},
		{"No elements", []int{}, []int{}},
		{"One element", []int{100}, []int{100}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewOrdered(tc.input...)
			assert.Equal(t, tc.expected, s.Members())
		})
	}
}

func TestOrderedSetInsertionOrderMembers(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Three elements", []int{3, 1, 2}, []int{3, 1, 2}},
		{"No elements", []int{}, []int{}},
		{"One element", []int{100}, []int{100}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewOrdered(tc.input...)
			assert.True(t, s.ContainsAll(s.InsertionOrderMembers()...))
		})
	}
}

func TestOrderedSetString(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected string
	}{
		{"Three elements", []int{3, 1, 2}, "[1 2 3]"},
		{"No elements", []int{}, "[]"},
		{"One element", []int{100}, "[100]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewOrdered(tc.input...)
			assert.Equal(t, tc.expected, s.String())
		})
	}
}

func TestOrderedSetMarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected string
	}{
		{"Three elements", []int{3, 1, 2}, "[1,2,3]"},
		{"No elements", []int{}, "[]"},
		{"One element", []int{100}, "[100]"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewOrdered(tc.input...)
			out, err := s.MarshalJSON()

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, string(out))
		})
	}
}

func TestOrderedSetUnmarshalJSON(t *testing.T) {
	in := []byte(`[1, 2, 3]`)
	s := NewOrdered[int]()

	err := s.UnmarshalJSON(in)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{1, 2, 3}, s.Members())
}

func TestOrderedSetAdd(t *testing.T) {
	s := NewOrdered[int]()
	s.Add(1, 2, 3)
	assert.ElementsMatch(t, []int{1, 2, 3}, s.Members())
}

func TestOrderedSetAddImmutable(t *testing.T) {
	s := NewOrdered(1, 2, 3)
	s2 := s.AddImmutable(4, 5)
	assert.ElementsMatch(t, []int{1, 2, 3}, s.Members())
	assert.ElementsMatch(t, []int{1, 2, 3, 4, 5}, s2.Members())
}

func TestOrderedSetRemove(t *testing.T) {
	s := NewOrdered(1, 2, 3)
	s.Remove(2, 3)
	assert.ElementsMatch(t, []int{1}, s.Members())
}

func TestOrderedSetRemoveImmutable(t *testing.T) {
	s := NewOrdered(1, 2, 3)
	s2 := s.RemoveImmutable(2, 3)
	assert.ElementsMatch(t, []int{1, 2, 3}, s.Members())
	assert.ElementsMatch(t, []int{1}, s2.Members())
}

func TestOrderedSetClear(t *testing.T) {
	s := NewOrdered(1, 2, 3)
	s.Clear()
	assert.Equal(t, 0, s.Len())
}

func TestOrderedSetUnion(t *testing.T) {
	s1 := NewOrdered(1, 2, 3)
	s2 := NewOrdered(3, 4, 5)
	s3 := s1.Union(s2)
	assert.ElementsMatch(t, []int{1, 2, 3, 4, 5}, s3.Members())
}

func TestOrderedSetIntersection(t *testing.T) {
	s1 := NewOrdered(1, 2, 3)
	s2 := NewOrdered(2, 3, 4)
	s3 := s1.Intersection(s2)
	assert.ElementsMatch(t, []int{2, 3}, s3.Members())
}

func TestOrderedSetDifference(t *testing.T) {
	s1 := NewOrdered(1, 2, 3)
	s2 := NewOrdered(2, 3, 4)
	s3 := s1.Difference(s2)
	assert.ElementsMatch(t, []int{1}, s3.Members())
}

func TestOrderedSetIsSubsetOf(t *testing.T) {
	s1 := NewOrdered(1, 2)
	s2 := NewOrdered(1, 2, 3)
	assert.True(t, s1.IsSubsetOf(s2))
	assert.False(t, s2.IsSubsetOf(s1))
}

func TestOrderedSetIsSupersetOf(t *testing.T) {
	s1 := NewOrdered(1, 2, 3)
	s2 := NewOrdered(1, 2)
	assert.True(t, s1.IsSupersetOf(s2))
	assert.False(t, s2.IsSupersetOf(s1))
}

func TestOrderedSetEqual(t *testing.T) {
	s1 := NewOrdered(1, 2, 3)
	s2 := NewOrdered(1, 2, 3)
	s3 := NewOrdered(1, 2)

	assert.True(t, s1.Equal(s2))
	assert.False(t, s1.Equal(s3))
}

func TestOrderedSetLen(t *testing.T) {
	s := NewOrdered(1, 2, 3)
	assert.Equal(t, 3, s.Len())
}

func TestOrderedSetUnmarshalJSONError(t *testing.T) {
	in := []byte(`{ "hello": "world"}`)
	s := NewOrdered[string]()
	err := json.Unmarshal(in, &s)

	assert.Error(t, err)
}
