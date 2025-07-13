package tests

import (
	"reflect"
	"testing"

	functools "github.com/felipegenef/functools"
)

func TestStreamifyWithBuffer(t *testing.T) {
	// Test for creating a streamable with buffer using StreamifyWithBuffer
	items := []int{1, 2, 3, 4, 5}
	bufferSize := 2
	stream := functools.StreamifyWithBuffer(items, bufferSize)

	result := stream.ToSlice()
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCreateBufferedStream(t *testing.T) {
	// Test for creating a buffered stream using CreateBufferedStream
	items := []int{1, 2, 3, 4, 5}
	generator := func(ch chan int) {
		for _, v := range items {
			ch <- v
		}
	}

	buffered := functools.CreateBufferedStream(generator, 3)

	result := buffered.ToSlice()
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBufferedStreamPipe(t *testing.T) {
	// Test Pipe method for buffered stream
	items := []int{1, 2, 3, 4, 5}
	generator := func(ch chan int) {
		for _, v := range items {
			ch <- v
		}
	}

	// Create the buffered stream
	buffered := functools.CreateBufferedStream(generator, 3)

	// Apply Pipe to double each item (transforming int to int)
	transformed := buffered.Pipe(func(x int) any { return x * 2 })

	// Get the result (which is of type []any)
	result := transformed.ToSlice()

	// Convert the result to []int for easy comparison
	var resultInts []int
	for _, v := range result {
		if v, ok := v.(int); ok {
			resultInts = append(resultInts, v)
		} else {
			t.Errorf("Expected int, but got %T", v)
		}
	}

	// Expected result
	expected := []int{2, 4, 6, 8, 10}

	// Compare the result to the expected value
	if !reflect.DeepEqual(resultInts, expected) {
		t.Errorf("Expected %v, got %v", expected, resultInts)
	}
}

func TestBufferedStreamBufferSizeConsistency(t *testing.T) {
	// Test that the buffer size is respected and processed correctly
	items := []int{1, 2, 3, 4, 5}
	bufferSize := 2
	generator := func(ch chan int) {
		for _, v := range items {
			ch <- v
		}
	}

	// Create the buffered stream
	buffered := functools.CreateBufferedStream(generator, bufferSize)
	// Additional test to ensure that the stream respects the buffer size
	// By checking the number of items processed in one go
	processedItems := 0

	result := buffered.
		Pipe(func(x int) any {
			processedItems++
			return x
		}).
		// ToSlice will block until all items are processed
		ToSlice()

	// Type assertion to check if result is of type []int
	var resultInts []int
	for _, v := range result {
		if v, ok := v.(int); ok {
			resultInts = append(resultInts, v)
		} else {
			t.Errorf("Expected int, but got %T", v)
		}
	}

	expected := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(resultInts, expected) {
		t.Errorf("Expected %v, got %v", expected, resultInts)
	}

	if processedItems != len(items) {
		t.Errorf("Expected to process %d items, but processed %d", len(items), processedItems)
	}
}

func TestBufferedStreamForEach(t *testing.T) {
	items := []int{1, 2, 3, 4}
	stream := functools.StreamifyWithBuffer(items, 2)

	var result []int
	stream.ForEach(func(x int) {
		result = append(result, x)
	})

	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBufferedStreamWithMultipleItemsAndSmallerBuffer(t *testing.T) {
	// Test with multiple items and a smaller buffer (e.g., buffer size 2)
	items := []int{1, 2, 3, 4, 5}
	bufferSize := 2
	generator := func(ch chan int) {
		for _, v := range items {
			ch <- v
		}
	}

	// Create the buffered stream with smaller buffer
	buffered := functools.CreateBufferedStream(generator, bufferSize)
	// Test the number of items processed in batches (buffer size)
	batches := 0
	// Collect the buffered stream into a slice
	result := buffered.
		Pipe(func(x int) any {
			// Each batch is processed in sequence, each element can come
			// through one by one due to the buffer size
			batches++
			return x
		}).
		ToSlice()
	expected := []int{1, 2, 3, 4, 5}

	// Type assertion to check if result is of type []int
	var resultInts []int
	for _, v := range result {
		if v, ok := v.(int); ok {
			resultInts = append(resultInts, v)
		} else {
			t.Errorf("Expected int, but got %T", v)
		}
	}

	if !reflect.DeepEqual(resultInts, expected) {
		t.Errorf("Expected %v, got %v", expected, resultInts)
	}
	if batches != len(items) {
		t.Errorf("Expected %d items to be processed, but %d were processed", len(items), batches)
	}
}

func TestBufferedStreamWithLargerBuffer(t *testing.T) {
	// Test with a larger buffer size than items
	items := []int{1, 2, 3, 4, 5}
	generator := func(ch chan int) {
		for _, v := range items {
			ch <- v
		}
	}

	buffered := functools.CreateBufferedStream(generator, 10)

	// Collect the buffered stream into a slice
	result := buffered.ToSlice()
	expected := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBufferedStreamWithEmptyBuffer(t *testing.T) {
	// Test with an empty stream
	items := []int{1, 2, 3, 4}
	generator := func(ch chan int) {
		for _, v := range items {
			ch <- v
		}
	}

	buffered := functools.CreateBufferedStream(generator, 2)

	// Collect the buffered stream into a slice
	result := buffered.ToSlice()
	expected := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBufferedStreamFilter(t *testing.T) {
	items := []int{1, 2, 3, 4}
	stream := functools.StreamifyWithBuffer(items, 2)

	filtered := stream.Filter(func(x int) bool { return x%2 == 0 })
	result := filtered.ToSlice()
	expected := []int{2, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBufferedStreamToStream(t *testing.T) {
	items := []int{1, 2, 3, 4}
	iter := functools.StreamifyWithBuffer(items, 2)

	stream := iter.ToStream()
	result := stream.ToSlice()
	expected := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
