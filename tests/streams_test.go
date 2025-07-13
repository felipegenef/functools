package tests

import (
	"reflect"
	"testing"

	functools "github.com/felipegenef/functools"
)

func TestStreamify(t *testing.T) {
	items := []int{1, 2, 3, 4}
	stream := functools.Streamify(items)

	result := stream.ToSlice()
	if !reflect.DeepEqual(result, items) {
		t.Errorf("Expected %v, got %v", items, result)
	}
}

func TestCreateStream(t *testing.T) {
	generator := func(ch chan int) {
		for i := 1; i <= 3; i++ {
			ch <- i
		}
	}

	stream := functools.CreateStream(generator)
	result := stream.ToSlice()
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestStreamPipe(t *testing.T) {
	// Test Pipe method for a stream
	items := []int{1, 2, 3, 4}
	stream := functools.Streamify(items)

	// Apply Pipe to double each item
	transformed := stream.Pipe(func(x int) any { return x * 2 })

	// Collect the transformed stream into a slice (which will be of type []any)
	result := functools.RecastStream[int](transformed).ToSlice()

	// Expected result
	expected := []int{2, 4, 6, 8}

	// Compare the result to the expected value
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestStreamFilter(t *testing.T) {
	items := []int{1, 2, 3, 4}
	stream := functools.Streamify(items)

	filtered := stream.Filter(func(x int) bool { return x%2 == 0 })
	result := filtered.ToSlice()
	expected := []int{2, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestStreamForEach(t *testing.T) {
	items := []int{1, 2, 3, 4}
	stream := functools.Streamify(items)

	var result []int
	stream.ForEach(func(x int) {
		result = append(result, x)
	})

	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestStreamToBufferedStream(t *testing.T) {
	items := []int{1, 2, 3, 4}
	stream := functools.Streamify(items)

	buffered := stream.ToBufferedStream(2)
	result := buffered.ToSlice()
	expected := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
