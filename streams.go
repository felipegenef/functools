package functools

// streamable is a collection that processes data on-demand via channels
type streamable[InputType any] struct {
	stream <-chan InputType
}

// Creates a streamable from a slice
func Streamify[InputType any](items []InputType) *streamable[InputType] {
	ch := make(chan InputType)
	go func() {
		defer close(ch)
		for _, v := range items {
			ch <- v
		}
	}()
	return &streamable[InputType]{stream: ch}
}

// CreateStream creates a streamable by receiving a generator function
// that generates values and sends them through the provided channel.
func CreateStream[InputType any](generator func(chan InputType)) *streamable[InputType] {
	ch := make(chan InputType)
	go func() {
		defer close(ch)
		generator(ch) // Call the generator with the channel
	}()
	return &streamable[InputType]{stream: ch}
}

// Pipe creates a new streamable by applying fn to each item
func (s *streamable[InputType]) Pipe(fn func(InputType) any) *streamable[any] {
	out := make(chan any)
	go func() {
		defer close(out)
		for v := range s.stream {
			out <- fn(v)
		}
	}()
	return &streamable[any]{stream: out}
}

// Filter creates a new streamable by filtering items with fn
func (s *streamable[InputType]) Filter(fn func(InputType) bool) *streamable[InputType] {
	out := make(chan InputType)
	go func() {
		defer close(out)
		for v := range s.stream {
			if fn(v) {
				out <- v
			}
		}
	}()
	return &streamable[InputType]{stream: out}
}

// ForEach consumes the stream by applying fn to each item
func (s *streamable[InputType]) ForEach(fn func(InputType)) {
	for v := range s.stream {
		fn(v)
	}
}

// ToSlice collects all items into a slice (may block until everything is consumed)
func (s *streamable[InputType]) ToSlice() []InputType {
	var result []InputType
	for v := range s.stream {
		result = append(result, v)
	}
	return result
}

func (s *streamable[InputType]) ToBufferedStream(bufferSize int) *bufferedStream[InputType] {
	ch := make(chan InputType, bufferSize)
	go func() {
		defer close(ch)
		for v := range s.stream {
			ch <- v
		}
	}()
	return &bufferedStream[InputType]{stream: ch}
}

func RecastStream[StreamType any](s *streamable[any]) *streamable[StreamType] {
	out := make(chan StreamType)
	go func() {
		defer close(out)
		for v := range s.stream {
			// Attempt to cast each item in the stream to OutputType
			if casted, ok := v.(StreamType); ok {
				out <- casted
			}
		}
	}()
	return &streamable[StreamType]{stream: out}
}
