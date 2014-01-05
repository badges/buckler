package shield

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func init() {
	LoadResources("../data")
}

func TestRenderString(t *testing.T) {
	i, err := os.Open("testdata/vendor.data")
	if err != nil {
		t.Fatal(err)
	}

	e, err := ioutil.ReadAll(i)
	if err != nil {
		t.Fatal(err)
	}

	r, _ := renderString("Vendor", c)
	if !bytes.Equal(r.Pix, e) {
		t.Error("make png shield 'use buckler blue' bytes not equal")
	}
}

// simple regression test
func TestPNG(t *testing.T) {
	i, err := os.Open("testdata/use-buckler-blue.png")
	if err != nil {
		t.Fatal(err)
	}

	e, err := ioutil.ReadAll(i)
	if err != nil {
		t.Fatal(err)
	}

	var b bytes.Buffer
	PNG(&b, Data{"use", "buckler", Blue})
	if !bytes.Equal(b.Bytes(), e) {
		t.Error("render string 'Vendor' bytes not equal")
	}
}

func BenchmarkRenderString(b *testing.B) {
	// c, the freetype context, is set up in png.go's init
	for i := 0; i < b.N; i++ {
		renderString("test string", c)
	}
}

func BenchmarkPNG(b *testing.B) {
	d := Data{"test", "output", Blue}
	for i := 0; i < b.N; i++ {
		PNG(ioutil.Discard, d)
	}
}
