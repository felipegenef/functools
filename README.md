# Functools - A Go Toolkit for Stream and Slice Operations
![Codecov](https://codecov.io/gh/felipegenef/functools/branch/main/graph/badge.svg)
![Test Status](https://img.shields.io/github/workflow/status/felipegenef/functools/coverage?label=Test%20Status&style=for-the-badge)


**Functools** is a versatile toolkit for processing and manipulating streams and slices in Go. It provides powerful functions for creating, transforming, filtering, and consuming data streams and slices, enabling you to write clean, functional-style code with minimal boilerplate.

This package provides a set of utilities that allow for smooth handling of buffered streams, slices, and transformations using the concept of **streams** (channels in Go) and **iterables** (slices). The toolkit includes operations for common functional patterns like `map`, `filter`, `reduce`, and more, all implemented in a way that minimizes memory overhead while maintaining high performance.

## Key Features

- **Streamable & BufferedStream**: Create streams from slices or generators, with support for buffered or unbuffered channels.
- **Functional Operations**: Chain operations such as `Map`, `Filter`, `ForEach`, `Reduce`, and more to easily transform and process data.
- **Flexible Slice Handling**: Apply functional methods directly to slices, supporting common operations like `Filter`, `Sort`, `Concat`, and `Slice`.
- **Efficient Consumption**: Consume data lazily with streaming, or eagerly with slices. 
- **Conversion Between Streams**: Easily convert between regular streams and buffered streams, providing flexibility depending on your needs. **Slices**, **Streams**, and **BufferedStreams** can be cast to each other using the methods `ToSlice()`, `ToStream()`, and `ToBufferedStream()` respectively.

## Stream-Based Memory Efficiency

One of the major advantages of using streams over slices or arrays is their ability to process data in **chunks** without loading the entire dataset into memory at once. 

In traditional slice-based approaches, the entire collection is loaded into memory, which can quickly become a bottleneck when dealing with large datasets. However, streams operate lazily, processing data item-by-item as needed. This reduces the memory footprint significantly, especially when dealing with large data sources that don’t need to be entirely loaded into memory at once.

For example, with **BufferedStreams**, you can define a buffer size that suits your available memory, and the data is processed in manageable chunks. This way, only a subset of the data is held in memory at any given time, drastically lowering the overall memory consumption while maintaining efficient processing.

This approach allows you to handle **large datasets** or **infinite data streams** without worrying about running out of memory, making the toolkit ideal for scalable, high-performance applications.


## Chaining Slice Functions for Easy Data Manipulation

One of the powerful features of **Functools** is the ability to **chain** multiple functional operations on slices with minimal boilerplate. Using a simple and readable syntax, you can easily perform complex data manipulations like filtering, transforming, reducing, and more, in a fluid sequence.

### Example: Chaining Slice Functions

Imagine you have a slice of integers, and you want to perform a series of operations such as:

1. Filter out even numbers.
2. Square the remaining numbers.
3. Sort them in ascending order.

Here’s how you can do it with **Functools**:

```go
package main

import (
	"fmt"
	"github.com/felipegenef/functools"
)

func main() {
	items := []int{1, 2, 3, 4, 5, 6}

	// Chain filter, map, and sort operations
	result := functools.Slicefy(items).
		Filter(func(i int) bool {
			return i%2 != 0 // Keep odd numbers
		}).
        // After map, due to the probable type modification, the result type must be `any`.
        // Further types in the chain must be cast to `int` where necessary.
		Map(func(i int) any {
			return i * i // Square the number
		}).
		Sort(func(a, b any) bool {
			// Cast `a` and `b` to int for comparison
			return a.(int) < b.(int) // Sort in ascending order
		}).
		ToSlice() // Convert the result back to a slice

	fmt.Println(result) // Output: [1 9 25]
}
```

#### Explanation:

- Slicefy converts the slice items into an iterable, enabling functional operations on it.
- Filter removes even numbers, keeping only the odd numbers.
- Map squares each remaining number.
- Sort sorts the squared numbers in ascending order.
- ToSlice collects the results back into a slice.

## Creating Streams from Slices and Generator Functions

Functools provides powerful methods to create streams either from slices or from custom generator functions, enabling easy handling of data in a memory-efficient, lazy-evaluated manner.


### Example: Creating a Stream from a Generator Function

With CreateStream, you can create a stream from a generator function. This is useful when you want to generate data dynamically, for example, from a file, database, or complex calculations. This is the most optimized way for memory consumption since you can load and pass data as a stream.

Creating a streamable from a slice means you had all the initial data in memory, meanwhile with a generator function you can have only one piece per loop iteration/pipeline batch.


```go
package main

import (
	"fmt"
	"github.com/felipegenef/functools"
)

func main() {
	// Generator function that yields values
	generator := func(ch chan int) {
		// Don't have to close the channel as CreateStream already does that after the function returns
		for i := 1; i <= 5; i++ {
			ch <- i // Send values to the channel
		}
	}

	// Create a stream from the generator
	stream := functools.CreateStream(generator)

	// Process the stream (example: multiply by 2)
	result := stream.
        // After pipe, due to the probable type modification, the result type must be `any`.
        // Further types in the chain must be cast to `int` where necessary.
		Pipe(func(i int) any {
			return i * 2 // Multiply by 2
		}).
		ToSlice() // Convert the result back to a slice

	fmt.Println(result) // Output: [2 4 6 8 10]
}

```

## BufferedStream: Efficient Handling of Backpressure & Large Datasets

BufferedStreams are ideal when you're dealing with backpressure or want to optimize memory usage when processing larger datasets in chunks.

By using buffered streams, you can control the buffer size to manage the amount of data held in memory at any time. This technique is useful when you need to load data efficiently while keeping memory usage under control.

### Example: BufferedStream from a Generator

This example demonstrates how to create a buffered stream from a generator function, efficiently processing large or infinite data streams in chunks.

This is the most optimized way for memory consumption once you can load and pass data as a stream. Creating a streamable from a slice means you had all the initial data in memory, meanwhile with a generator function you can have only one piece per loop iteration/pipeline batch.

```go
package main

import (
	"fmt"
	"github.com/felipegenef/functools"
)

func main() {
	// Generator function that produces values in chunks
	generator := func(ch chan int) {
		// Don't have to close the channel as CreateBufferedStream already does that after the function returns
		for i := 1; i <= 1000000; i++ {
			ch <- i
		}
	}

	// Create a buffered stream with a buffer size of 1000
	stream := functools.CreateBufferedStream(generator, 1000)

	// Process the stream (example: multiply by 2)
	result := stream.
		// After pipe, due to the probable type modification, the result type must be `any`.
		// Further types in the chain must be cast to `int` where necessary.
		Pipe(func(i int) any {
			return i * 2 // Multiply by 2
		}).
		ToSlice() // Convert the result back to a slice

	fmt.Println(result[:10]) // Output: [2 4 6 8 10 12 14 16 18 20]
}

```

#### Explanation:

- CreateBufferedStream creates a stream from the generator, processing data in buffered chunks of 1000 items at a time.
- The Map operation multiplies each number by 2.
- ToSlice collects the processed results into a slice.

## Type Safety and Casting

When using any in Go, be mindful of type safety. After performing transformations (such as Map or Pipe), the type of data may change. To maintain type safety, ensure that the correct type is cast when performing operations that expect a specific type.

For example, in the chaining operations where Map or Pipe may modify the type of data (e.g., converting integers to any), it’s important to cast the values back to their expected types during operations like Sort or Reduce.

```go
// Sort requires casting back to the original type since the Map step returned `any`
Sort(func(a, b any) bool {
    return a.(int) < b.(int) // Ensure correct type before comparison
})

```

## Additional Notes

- **Backpressure**: When using BufferedStream, backpressure occurs when the producer generates data faster than the consumer can process it