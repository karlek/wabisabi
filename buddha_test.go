package main

import (
	"testing"
)

// var benchTries float64 = 100000

func BenchmarkFillHistogram1_1(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 1, 128)
	}
}
func BenchmarkFillHistogram1_2(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 1, 128)
	}
}
func BenchmarkFillHistogram1_4(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 1, 128)
	}
}
func BenchmarkFillHistogram1_8(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 1, 128)
	}
}
func BenchmarkFillHistogram1_16(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 1, 128)
	}
}
func BenchmarkFillHistogram1_32(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 1, 128)
	}
}
func BenchmarkFillHistogram1_64(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 1, 128)
	}
}
func BenchmarkFillHistogram1_128(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 1, 128)
	}
}

func BenchmarkFillHistogram2_1(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 2, 1)
	}
}
func BenchmarkFillHistogram2_2(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 2, 2)
	}
}
func BenchmarkFillHistogram2_4(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 2, 4)
	}
}
func BenchmarkFillHistogram2_8(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 2, 8)
	}
}
func BenchmarkFillHistogram2_16(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 2, 16)
	}
}
func BenchmarkFillHistogram2_32(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 2, 32)
	}
}
func BenchmarkFillHistogram2_64(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 2, 64)
	}
}
func BenchmarkFillHistogram2_128(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 2, 128)
	}
}

func BenchmarkFillHistogram4_128(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 4, 128)
	}
}
func BenchmarkFillHistogram8_128(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 8, 128)
	}
}
func BenchmarkFillHistogram32_128(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 32, 128)
	}
}
func BenchmarkFillHistogram64_128(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 64, 128)
	}
}
func BenchmarkFillHistogram128_128(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 128, 128)
	}
}
func BenchmarkFillHistogram16_1(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 1)
	}
}
func BenchmarkFillHistogram16_2(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 2)
	}
}
func BenchmarkFillHistogram16_4(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 4)
	}
}
func BenchmarkFillHistogram16_8(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 8)
	}
}
func BenchmarkFillHistogram16_16(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 16)
	}
}
func BenchmarkFillHistogram16_32(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 256)
	}
}
func BenchmarkFillHistogram16_64(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 256)
	}
}
func BenchmarkFillHistogram16_128(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 128)
	}
}
func BenchmarkFillHistogram16_256(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 256)
	}
}
func BenchmarkFillHistogram16_512(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 512)
	}
}
func BenchmarkFillHistogram16_1024(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 1024)
	}
}
func BenchmarkFillHistogram16_2048(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 2048)
	}
}
func BenchmarkFillHistogram16_4096(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 4096)
	}
}
func BenchmarkFillHistogram16_8192(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 8192)
	}
}
func BenchmarkFillHistogram16_16384(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 16384)
	}
}
func BenchmarkFillHistogram16_32768(b *testing.B) {
	red, green, blue := &Histo{}, &Histo{}, &Histo{}
	// tries = benchTries
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillHistograms(red, green, blue, 16, 32768)
	}
}
