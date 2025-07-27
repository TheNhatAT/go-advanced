package micro

import "testing"

func BenchmarkSum_template(b *testing.B) {
	b.ReportAllocs()
	// TODO(bwplotka): Add any initialization that is needed.
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// TODO(bwplotka): Add tested functionality.
	}
}
