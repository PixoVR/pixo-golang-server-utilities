package common

// Package common provides common data structures and utilities.
//
// Set is a generic set implementation that can be used to store unique items of any comparable type.
//
// Usage:
//
//	// Create a new set of strings
//	stringSet := common.NewSet[string]()
//
//	// Add items to the set
//	stringSet.Add("apple")
//	stringSet.Add("banana")
//	stringSet.Add("apple") // Adding a duplicate item has no effect
//
//	// Check if an item is in the set
//	exists := stringSet.Contains("banana") // returns true
//
//	// Remove an item from the set
//	stringSet.Remove("banana")
//
//	// Get the size of the set
//	size := stringSet.Size() // returns 1
//
//	// Convert the set to a slice
//	slice := stringSet.ToSlice() // returns []string{"apple"} (order is not guaranteed)
//
//	// Clear the set
//	stringSet.Clear()
//

type Set[T comparable] struct {
	items map[T]struct{}
}

// NewSet creates and returns a new Set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		items: make(map[T]struct{}),
	}
}

// Add adds an item to the set
func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

// Remove removes an item from the set
func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

// Contains checks if an item exists in the set
func (s *Set[T]) Contains(item T) bool {
	_, exists := s.items[item]
	return exists
}

// Size returns the number of items in the set
func (s *Set[T]) Size() int {
	return len(s.items)
}

// Clear removes all items from the set
func (s *Set[T]) Clear() {
	s.items = make(map[T]struct{})
}

// ToSlice converts the set to a slice
func (s *Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s.items))
	for item := range s.items {
		slice = append(slice, item)
	}
	return slice
}

