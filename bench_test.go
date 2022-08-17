package protorand

import (
	"testing"

	testpb "github.com/wattch/protorand/testdata"
)

func BenchmarkGen(b *testing.B) {
	pr := New()
	input := &testpb.TestMessage{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pr.Gen(input)
	}
}
