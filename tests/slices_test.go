package tests

import (
	"reflect"
	"testing"

	functools "github.com/felipegenef/functools"
)

func TestSlicefy(t *testing.T) {
	items := []string{"apple", "banana", "cherry"}
	iter := functools.Slicefy(items)

	result := iter.ToSlice()
	if !reflect.DeepEqual(result, items) {
		t.Errorf("Expected %v, got %v", items, result)
	}
}

func TestIterableFilter(t *testing.T) {
	items := []int{1, 2, 3, 4}
	iter := functools.Slicefy(items)

	filtered := iter.Filter(func(x int) bool { return x%2 == 0 })
	result := filtered.ToSlice()
	expected := []int{2, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestIterableMap(t *testing.T) {
	items := []int{1, 2, 3, 4}
	iter := functools.Slicefy(items)

	mapped := iter.Map(func(x int) any { return x * x })
	result := mapped.ToSlice()
	expected := []interface{}{1, 4, 9, 16}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestIterableReduce(t *testing.T) {
	items := []int{1, 2, 3, 4}
	iter := functools.Slicefy(items)

	sum := iter.Reduce(func(acc, item int) int { return acc + item }, 0)
	expected := 10

	if sum != expected {
		t.Errorf("Expected %d, got %d", expected, sum)
	}
}

func TestIterableFind(t *testing.T) {
	items := []int{1, 2, 3, 4}
	iter := functools.Slicefy(items)

	result := iter.Find(func(x int) bool { return x == 3 })
	if *result != 3 {
		t.Errorf("Expected 3, got %v", result)
	}

	// Test when element is not found
	result = iter.Find(func(x int) bool { return x == 5 })
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestIterableSome(t *testing.T) {
	items := []int{1, 2, 3, 4}
	iter := functools.Slicefy(items)

	result := iter.Some(func(x int) bool { return x == 3 })
	if !result {
		t.Errorf("Expected true, got false")
	}

	result = iter.Some(func(x int) bool { return x == 5 })
	if result {
		t.Errorf("Expected false, got true")
	}
}

func TestIterableEvery(t *testing.T) {
	items := []int{1, 2, 3, 4}
	iter := functools.Slicefy(items)

	result := iter.Every(func(x int) bool { return x < 5 })
	if !result {
		t.Errorf("Expected true, got false")
	}

	result = iter.Every(func(x int) bool { return x < 4 })
	if result {
		t.Errorf("Expected false, got true")
	}
}

func TestIterableSort(t *testing.T) {
	items := []int{4, 3, 2, 1}
	iter := functools.Slicefy(items)

	sorted := iter.Sort(func(a, b int) bool { return a < b })
	result := sorted.ToSlice()
	expected := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestIterableConcat(t *testing.T) {
	items1 := []int{1, 2}
	items2 := []int{3, 4}
	iter1 := functools.Slicefy(items1)

	concatenated := iter1.Concat(items2)
	result := concatenated.ToSlice()
	expected := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestIterableForEach(t *testing.T) {
	items := []int{1, 2, 3, 4}
	stream := functools.Slicefy(items)

	var result []int
	stream.ForEach(func(x int) {
		result = append(result, x)
	})

	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestIterableSlice(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	iter := functools.Slicefy(items)

	sliced := iter.Slice(1, 3)
	result := sliced.ToSlice()
	expected := []int{2, 3}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test invalid slice
	sliced = iter.Slice(5, 3)
	if len(sliced.ToSlice()) != 0 {
		t.Errorf("Expected empty slice, got %v", sliced.ToSlice())
	}
}

func TestIterableToStream(t *testing.T) {
	items := []int{1, 2, 3, 4}
	iter := functools.Slicefy(items)

	stream := iter.ToStream()
	result := stream.ToSlice()
	expected := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestIterableToBufferedStream(t *testing.T) {
	items := []int{1, 2, 3, 4}
	iter := functools.Slicefy(items)

	buffered := iter.ToBufferedStream(2)
	result := buffered.ToSlice()
	expected := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
