package functools

import "sort"

// iterable holds a generic slice with two type parameters
type iterable[InputType any] struct {
	items []InputType
}

// Factory to create iterable with a variadic type for ResponseType
func Slicefy[InputType any](items []InputType) *iterable[InputType] {
	return &iterable[InputType]{items: items}
}

// Filter returns a new iterable with only the items that pass the filter (without changing the type)
func (c *iterable[InputType]) Filter(fn func(InputType) bool) *iterable[InputType] {
	var result []InputType
	for _, v := range c.items {
		if fn(v) {
			result = append(result, v)
		}
	}
	return &iterable[InputType]{items: result}
}

// ForEach executes the function fn on each item (no return)
func (c *iterable[InputType]) ForEach(fn func(InputType)) {
	for _, v := range c.items {
		fn(v)
	}
}

// Map applies the transformation function fn and returns a new iterable
func (c *iterable[InputType]) Map(fn func(InputType) any) *iterable[any] {
	var result []any
	for _, v := range c.items {
		result = append(result, fn(v))
	}
	return &iterable[any]{items: result}
}

// ToSlice returns the internal slice (for anyone to access directly)
func (c *iterable[InputType]) ToSlice() []InputType {
	return c.items
}

// Reduce reduces the iterable to a single value based on the provided function.
func (c *iterable[InputType]) Reduce(fn func(acc InputType, item InputType) InputType, initial InputType) InputType {
	acc := initial
	for _, v := range c.items {
		acc = fn(acc, v)
	}
	return acc
}

// Find returns the first element that satisfies the condition or nil.
func (c *iterable[InputType]) Find(fn func(InputType) bool) *InputType {
	for _, v := range c.items {
		if fn(v) {
			return &v
		}
	}
	return nil
}

// Some checks if at least one element satisfies the condition.
func (c *iterable[InputType]) Some(fn func(InputType) bool) bool {
	for _, v := range c.items {
		if fn(v) {
			return true
		}
	}
	return false
}

// Every checks if every element satisfies the condition.
func (c *iterable[InputType]) Every(fn func(InputType) bool) bool {
	for _, v := range c.items {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Sort sorts the elements in ascending order using a comparison function.
func (c *iterable[InputType]) Sort(fn func(a, b InputType) bool) *iterable[InputType] {
	sorted := append([]InputType{}, c.items...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return fn(sorted[i], sorted[j])
	})
	return &iterable[InputType]{items: sorted}
}

// Concat concatenates the current iterable with another iterable and returns a new iterable.
func (c *iterable[InputType]) Concat(other []InputType) *iterable[InputType] {
	// Combine the two slices
	combinedItems := append(c.items, other...)
	return &iterable[InputType]{items: combinedItems}
}

// Slice extracts a subset of the iterable (like slicing an array).
func (c *iterable[InputType]) Slice(start, end int) *iterable[InputType] {
	if start < 0 || end > len(c.items) || start > end {
		return &iterable[InputType]{items: []InputType{}}
	}
	return &iterable[InputType]{items: c.items[start:end]}
}

// ToStream converts an iterable to a streamable
func (c *iterable[InputType]) ToStream() *streamable[InputType] {
	ch := make(chan InputType)
	go func() {
		defer close(ch)
		for _, v := range c.items {
			ch <- v
		}
	}()
	return &streamable[InputType]{stream: ch}
}

// ToBufferedStream converts an iterable (using a slice) to a buffered streamable
func (c *iterable[InputType]) ToBufferedStream(bufferSize int) *bufferedStream[InputType] {
	ch := make(chan InputType, bufferSize)
	go func() {
		defer close(ch)
		for _, v := range c.items {
			ch <- v
		}
	}()
	return &bufferedStream[InputType]{stream: ch}
}

func RecastSlice[SliceType any](input *iterable[any]) *iterable[SliceType] {
	var result []SliceType
	for _, v := range input.items {
		// Assuming the input can be directly converted to SliceType.
		// You'll need to handle the conversion, potentially with type assertions.
		if casted, ok := v.(SliceType); ok {
			result = append(result, casted)
		}
	}
	return &iterable[SliceType]{items: result}
}
