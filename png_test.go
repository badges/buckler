package main

import (
	"testing"
)

func BenchmarkRenderString(b *testing.B) {
	// c, the freetype context, is set up in png.go's init
	for i := 0; i < b.N; i++ {
		renderString("test string", c);
	}
}
