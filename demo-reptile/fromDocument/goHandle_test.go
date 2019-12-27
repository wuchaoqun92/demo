package fromDocument

import (
	"testing"
)

func TestImageCompression(t *testing.T) {
	a := ImageCompression("/Users/wuchaoqun/Desktop/codeMaterial/3.jpg")
	if len(a) > 34655 {
		t.Error("compress failed")
	}
}

func Benchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ImageCompression("/Users/wuchaoqun/Desktop/codeMaterial/3.jpg")
	}
}
