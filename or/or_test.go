package or

import (
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	// Helper function to create a signal channel that closes after a duration
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	t.Run("test empty channels", func(t *testing.T) {
		result := Or()
		if result != nil {
			t.Errorf("Expected nil channel for empty input, got %v", result)
		}
	})

	t.Run("test single channel", func(t *testing.T) {
		ch := sig(100 * time.Millisecond)
		result := Or(ch)
		if result != ch {
			t.Error("Expected original channel for single input")
		}
	})

	t.Run("test multiple channels - fastest one wins", func(t *testing.T) {
		start := time.Now()

		<-Or(
			sig(2*time.Hour),
			sig(5*time.Minute),
			sig(100*time.Millisecond), // This should be the fastest
			sig(1*time.Hour),
			sig(1*time.Minute),
		)

		elapsed := time.Since(start)

		// Should finish in around 100ms (with some tolerance)
		if elapsed > 200*time.Millisecond {
			t.Errorf("Expected to finish in ~100ms, took %v", elapsed)
		}
		if elapsed < 50*time.Millisecond {
			t.Errorf("Finished too fast, expected ~100ms, took %v", elapsed)
		}
	})

	t.Run("test multiple channels - immediate closure", func(t *testing.T) {
		closedCh := make(chan interface{})
		close(closedCh)

		start := time.Now()

		<-Or(
			sig(1*time.Second),
			closedCh,
			sig(500*time.Millisecond),
		)

		elapsed := time.Since(start)

		// Should finish immediately since closedCh is already closed
		if elapsed > 50*time.Millisecond {
			t.Errorf("Expected immediate closure, took %v", elapsed)
		}
	})

	t.Run("test multiple goroutines don't cause race conditions", func(t *testing.T) {
		var chans []<-chan interface{}
		for i := 0; i < 10; i++ {
			chans = append(chans, sig(time.Duration(100+i*50)*time.Millisecond))
		}

		start := time.Now()

		<-Or(chans...)

		elapsed := time.Since(start)

		// Should finish in around 100ms (fastest channel)
		if elapsed > 150*time.Millisecond {
			t.Errorf("Expected to finish in ~100ms, took %v", elapsed)
		}
	})
}
