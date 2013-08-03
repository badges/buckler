package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestRenderString(t *testing.T) {
	i, _ := os.Open("test/vendor.data")
	e, _ := ioutil.ReadAll(i)

	r, _ := renderString("Vendor", c)
	if !bytes.Equal(r.Pix, e) {
		t.Error("Failure")
	}
}

func BenchmarkRenderString(b *testing.B) {
	// c, the freetype context, is set up in png.go's init
	for i := 0; i < b.N; i++ {
		renderString("test string", c)
	}
}
