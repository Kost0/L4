package or

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestOr_Success(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch1 := make(chan interface{})
		ch2 := make(chan interface{})
		ch3 := make(chan interface{})

		result := Or(ch1, ch2, ch3)

		close(ch1)

		synctest.Wait()

		select {
		case <-result:
		default:
			t.Error("Expected closed channel")
		}
	})
}

func TestOr_EmptyChannels(t *testing.T) {
	result := Or()
	select {
	case <-result:
	default:
		t.Error("Expected closed channel")
	}
}

func TestOr_SingleChannel(t *testing.T) {
	ch := make(chan interface{})
	result := Or(ch)

	if result != ch {
		t.Error("Expected to return the same channel")
	}
}

func TestOr_CloseSecondChannel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch1 := make(chan interface{})
		ch2 := make(chan interface{})
		ch3 := make(chan interface{})

		result := Or(ch1, ch2, ch3)

		close(ch2)

		synctest.Wait()

		select {
		case <-result:
		default:
			t.Error("Expected closed channel")
		}
	})
}

func TestOr_MultipleChannels(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch1 := make(chan interface{})
		ch2 := make(chan interface{})
		ch3 := make(chan interface{})

		result := Or(ch1, ch2, ch3)

		go close(ch1)
		go close(ch2)
		go close(ch3)

		synctest.Wait()

		select {
		case <-result:
		default:
			t.Error("Expected closed channel")
		}
	})
}

func TestOr_AlreadyClosedChannel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ch1 := make(chan interface{})
		close(ch1)

		ch2 := make(chan interface{})
		ch3 := make(chan interface{})

		result := Or(ch1, ch2, ch3)

		synctest.Wait()

		select {
		case <-result:
		default:
			t.Error("Expected closed channel")
		}
	})
}

func TestOr_NoCloseChannels(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})

	result := Or(ch1, ch2, ch3)

	select {
	case <-result:
		t.Error("Expected open channel")
	default:
	}
}

func TestOr_WithTimeout(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	elapsed := time.Since(start)

	if elapsed > 2*time.Second {
		t.Errorf("Expected to finish in 1 second, took %v", elapsed)
	}

	t.Logf("Finished in %v second", elapsed)
}

func TestOr_WithTimeout_Synctest(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.After(after)
		}()
		return c
	}

	synctest.Test(t, func(t *testing.T) {
		result := Or(
			sig(2*time.Hour),
			sig(5*time.Minute),
			sig(1*time.Second),
			sig(1*time.Hour),
			sig(1*time.Minute),
		)

		time.Sleep(1100 * time.Millisecond)

		synctest.Wait()

		select {
		case <-result:
		default:
			t.Error("Expected closed channel")
		}
	})
}
