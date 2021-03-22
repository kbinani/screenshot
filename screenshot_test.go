package screenshot

import (
	"testing"
)

func TestCaptureRect(t *testing.T) {
	bounds := GetDisplayBounds(0)
	_, err := CaptureRect(bounds)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkCaptureRect(t *testing.B) {
	bounds := GetDisplayBounds(0)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_, err := CaptureRect(bounds)
		if err != nil {
			t.Error(err)
		}
	}
}
