package consumer

import (
	"testing"
)

func TestSimplePush(t *testing.T) {
	var q QueueMgr
	var in = []int{40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50}
	var size = len(in)

	var err error

	// Push them in and test for correct length and no errors
	for i, v := range in {
		err = q.Push(v)

		if err != nil {
			t.Errorf("Push got error: %v", err)
		}

		if q.Size != i+1 {
			t.Errorf("QueueMgr size not %d, got %d", i+1, q.Size)
		}
	}

	// Pop them out and make sure they are the correct value and the correct size
	for i := range in {
		b, err := q.Pop()
		if err != nil {
			t.Errorf("Pop returned an error, got %v", err)
		}
		if b != in[i] {
			t.Errorf("Pop did not work, got %v instead of %d", b, in[i])
		}

		if q.Size != (size - i - 1) {
			t.Errorf("QueueMgr size not %d, got %d", (size - i - 1), q.Size)
		}
	}
}
