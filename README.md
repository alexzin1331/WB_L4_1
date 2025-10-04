# OR Channel Implementation

This project implements a Go function that combines multiple "done" channels into a single channel that closes as soon as any of the source channels close.

## The Problem

When working with concurrent operations, you often need to wait for multiple channels to signal completion. Instead of manually monitoring each channel, you want a single channel that signals completion when **any** of the input channels signal completion.

## The Solution

The `Or` function takes any number of input channels and returns a single channel that:
- Closes immediately when **any** input channel closes
- Uses an infinite loop with `select` to monitor all channels simultaneously
- Safely closes the output channel using `sync.Once` to prevent race conditions

## Usage Example

```go
sig := func(after time.Duration) <-chan interface{} {
    c := make(chan interface{})
    go func() {
        defer close(c)
        time.Sleep(after)
    }()
    return c
}

start := time.Now()
<-or.Or(
    sig(2*time.Hour),
    sig(5*time.Minute),
    sig(1*time.Second),    // This closes first!
    sig(1*time.Hour),
    sig(1*time.Minute),
)
fmt.Printf("done after %v", time.Since(start)) // Output: ~1 second
```

## Key Concepts

- **Channel orchestration**: Combining multiple channels into one
- **Race condition safety**: Using `sync.Once` for safe channel closing
- **Concurrent monitoring**: Using `select` in goroutines to read from multiple channels
- **Fastest-first closing**: The output channel closes when the fastest input channel closes

## Running the Project

```bash
# Run tests
go test ./or

# Run example
cd example && go run main.go
```
