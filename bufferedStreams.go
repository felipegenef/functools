package functools

// bufferedStream is a collection that processes data on-demand via buffered channels
type bufferedStream[InputType any] struct {
	stream <-chan InputType
}

func StreamifyWithBuffer[InputType any](items []InputType, bufferSize int) *bufferedStream[InputType] {
	ch := make(chan InputType, bufferSize)
	go func() {
		defer close(ch)
		for _, v := range items {
			ch <- v
		}
	}()
	return &bufferedStream[InputType]{stream: ch}
}

// CreateBufferedStream creates a buffered streamable with the given bufferSize
func CreateBufferedStream[InputType any](generator func(chan InputType), bufferSize int) *bufferedStream[InputType] {
	ch := make(chan InputType, bufferSize)
	go func() {
		defer close(ch)
		generator(ch) // Call the generator with the channel
	}()
	return &bufferedStream[InputType]{stream: ch}
}

// Pipe creates a new streamable by applying fn to each item
func (s *bufferedStream[InputType]) Pipe(fn func(InputType) InputType) *bufferedStream[any] {
	out := make(chan any)
	go func() {
		defer close(out)
		for v := range s.stream {
			out <- fn(v)
		}
	}()
	return &bufferedStream[any]{stream: out}
}

// Filter creates a new streamable by filtering items with fn
func (s *bufferedStream[InputType]) Filter(fn func(InputType) bool) *bufferedStream[InputType] {
	out := make(chan InputType)
	go func() {
		defer close(out)
		for v := range s.stream {
			if fn(v) {
				out <- v
			}
		}
	}()
	return &bufferedStream[InputType]{stream: out}
}

// ForEach consumes the stream by applying fn to each item
func (s *bufferedStream[InputType]) ForEach(fn func(InputType)) {
	for v := range s.stream {
		fn(v)
	}
}

// ToSlice collects all items from the buffered stream into a slice (may block until everything is consumed)
func (s *bufferedStream[InputType]) ToSlice() []InputType {
	var result []InputType
	for v := range s.stream {
		result = append(result, v)
	}
	return result
}

// ToStream converts a buffered streamable into a regular streamable (unbuffered channel)
func (s *bufferedStream[InputType]) ToStream() *streamable[InputType] {
	ch := make(chan InputType)
	go func() {
		defer close(ch)
		for v := range s.stream {
			ch <- v
		}
	}()
	return &streamable[InputType]{stream: ch}
}
