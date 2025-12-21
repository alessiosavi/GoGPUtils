package collection

// Set is a generic unordered collection of unique elements.
// The zero value is not usable; use NewSet to create a Set.
type Set[T comparable] struct {
	items map[T]struct{}
}

// NewSet creates a new empty Set.
//
// Example:
//
//	set := NewSet[int]()
//	set.Add(1, 2, 3)
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		items: make(map[T]struct{}),
	}
}

// NewSetWithCapacity creates a Set with pre-allocated capacity.
func NewSetWithCapacity[T comparable](capacity int) *Set[T] {
	return &Set[T]{
		items: make(map[T]struct{}, capacity),
	}
}

// NewSetFrom creates a Set from a slice.
//
// Example:
//
//	set := NewSetFrom([]int{1, 2, 2, 3}) // {1, 2, 3}
func NewSetFrom[T comparable](items []T) *Set[T] {
	s := NewSetWithCapacity[T](len(items))
	for _, item := range items {
		s.items[item] = struct{}{}
	}

	return s
}

// Add adds one or more elements to the set.
// Elements already in the set are ignored.
//
// Example:
//
//	set.Add(1, 2, 3)
func (s *Set[T]) Add(items ...T) {
	for _, item := range items {
		s.items[item] = struct{}{}
	}
}

// Remove removes one or more elements from the set.
// Elements not in the set are ignored.
func (s *Set[T]) Remove(items ...T) {
	for _, item := range items {
		delete(s.items, item)
	}
}

// Contains returns true if the element is in the set.
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.items[item]

	return ok
}

// ContainsAll returns true if all elements are in the set.
func (s *Set[T]) ContainsAll(items ...T) bool {
	for _, item := range items {
		if _, ok := s.items[item]; !ok {
			return false
		}
	}

	return true
}

// ContainsAny returns true if any element is in the set.
func (s *Set[T]) ContainsAny(items ...T) bool {
	for _, item := range items {
		if _, ok := s.items[item]; ok {
			return true
		}
	}

	return false
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	return len(s.items)
}

// IsEmpty returns true if the set has no elements.
func (s *Set[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Clear removes all elements from the set.
func (s *Set[T]) Clear() {
	s.items = make(map[T]struct{})
}

// Values returns all elements in the set.
// Order is not guaranteed.
func (s *Set[T]) Values() []T {
	result := make([]T, 0, len(s.items))
	for item := range s.items {
		result = append(result, item)
	}

	return result
}

// Union returns a new set containing all elements from both sets.
//
// Example:
//
//	a := NewSetFrom([]int{1, 2})
//	b := NewSetFrom([]int{2, 3})
//	c := a.Union(b) // {1, 2, 3}
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSetWithCapacity[T](len(s.items) + len(other.items))
	for item := range s.items {
		result.items[item] = struct{}{}
	}

	for item := range other.items {
		result.items[item] = struct{}{}
	}

	return result
}

// Intersection returns a new set containing elements in both sets.
//
// Example:
//
//	a := NewSetFrom([]int{1, 2, 3})
//	b := NewSetFrom([]int{2, 3, 4})
//	c := a.Intersection(b) // {2, 3}
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	// Iterate over smaller set for efficiency
	smaller, larger := s, other
	if len(s.items) > len(other.items) {
		smaller, larger = other, s
	}

	result := NewSet[T]()

	for item := range smaller.items {
		if _, ok := larger.items[item]; ok {
			result.items[item] = struct{}{}
		}
	}

	return result
}

// Difference returns a new set containing elements in s but not in other.
//
// Example:
//
//	a := NewSetFrom([]int{1, 2, 3})
//	b := NewSetFrom([]int{2, 3, 4})
//	c := a.Difference(b) // {1}
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := NewSet[T]()

	for item := range s.items {
		if _, ok := other.items[item]; !ok {
			result.items[item] = struct{}{}
		}
	}

	return result
}

// SymmetricDifference returns elements in either set but not both.
//
// Example:
//
//	a := NewSetFrom([]int{1, 2, 3})
//	b := NewSetFrom([]int{2, 3, 4})
//	c := a.SymmetricDifference(b) // {1, 4}
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T] {
	result := NewSet[T]()

	for item := range s.items {
		if _, ok := other.items[item]; !ok {
			result.items[item] = struct{}{}
		}
	}

	for item := range other.items {
		if _, ok := s.items[item]; !ok {
			result.items[item] = struct{}{}
		}
	}

	return result
}

// IsSubset returns true if all elements of s are in other.
func (s *Set[T]) IsSubset(other *Set[T]) bool {
	for item := range s.items {
		if _, ok := other.items[item]; !ok {
			return false
		}
	}

	return true
}

// IsSuperset returns true if s contains all elements of other.
func (s *Set[T]) IsSuperset(other *Set[T]) bool {
	return other.IsSubset(s)
}

// Equal returns true if both sets contain the same elements.
func (s *Set[T]) Equal(other *Set[T]) bool {
	if len(s.items) != len(other.items) {
		return false
	}

	for item := range s.items {
		if _, ok := other.items[item]; !ok {
			return false
		}
	}

	return true
}

// Clone returns a copy of the set.
func (s *Set[T]) Clone() *Set[T] {
	result := NewSetWithCapacity[T](len(s.items))
	for item := range s.items {
		result.items[item] = struct{}{}
	}

	return result
}

// ForEach calls a function for each element in the set.
func (s *Set[T]) ForEach(fn func(T)) {
	for item := range s.items {
		fn(item)
	}
}

// Filter returns a new set with elements that satisfy the predicate.
func (s *Set[T]) Filter(predicate func(T) bool) *Set[T] {
	result := NewSet[T]()

	for item := range s.items {
		if predicate(item) {
			result.items[item] = struct{}{}
		}
	}

	return result
}
