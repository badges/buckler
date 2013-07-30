package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/raster"
	"code.google.com/p/freetype-go/freetype/truetype"
)

type Data struct {
	Vendor string
	Status string
	Color  color.RGBA
}

var (
	Grey        = color.RGBA{74, 74, 74, 255}
	BrightGreen = color.RGBA{69, 203, 20, 255}
	Green       = color.RGBA{124, 166, 0, 255}
	YellowGreen = color.RGBA{156, 158, 9, 255}
	Yellow      = color.RGBA{184, 148, 19, 255}
	Orange      = color.RGBA{184, 113, 37, 255}
	Red         = color.RGBA{186, 77, 56, 255}
	LightGrey   = color.RGBA{131, 131, 131, 255}
	Blue        = color.RGBA{0, 126, 198, 255}

	Colors = map[string]color.RGBA{
		"grey":        Grey,
		"brightgreen": BrightGreen,
		"green":       Green,
		"yellowgreen": YellowGreen,
		"yellow":      Yellow,
		"orange":      Orange,
		"red":         Red,
		"lightgrey":   LightGrey,
		"blue":        Blue,

		// US spelling
		"gray":      Grey,
		"lightgray": LightGrey,
	}

	shadow = color.RGBA{0, 0, 0, 125}

	edge     image.Image
	gradient image.Image
	font     *truetype.Font
	c        *freetype.Context
)

const (
	h  = 18
	op = 4
	ip = 4
)

func init() {
	log.Println("Initializing png")

	fi, _ := os.Open("edge.png")
	edge, _ = png.Decode(fi)
	defer fi.Close()

	fi, _ = os.Open("gradient.png")
	gradient, _ = png.Decode(fi)
	defer fi.Close()

	fontBytes, err := ioutil.ReadFile("opensanssemibold.ttf")
	if err != nil {
		log.Println(err)
	}

	font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
	}

	c = freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(10)
}

func getTextOffset(pt raster.Point) int {
	return int(math.Floor(float64(float32(pt.X)/256 + 0.5)))
}

func renderString(s string, c *freetype.Context) (image.Image, int) {
	estWidth := 8 * len(s)
	dst := image.NewRGBA(image.Rect(0, 0, estWidth, h))

	c.SetDst(dst)
	c.SetClip(dst.Bounds())

	c.SetSrc(&image.Uniform{shadow})
	pt := freetype.Pt(0, 13)
	c.DrawString(s, pt)

	c.SetSrc(image.White)

	pt = freetype.Pt(0, 12)
	pt, _ = c.DrawString(s, pt)

	return dst, getTextOffset(pt)
}

func makePngShield(w http.ResponseWriter, d Data) {
	w.Header().Add("content-type", "image/png")

	// render text to determine how wide the image has to be
	// we leave 6 pixels at the start and end, and 3 for each in the middle
	v, vw := renderString(d.Vendor, c)
	s, sw := renderString(d.Status, c)
	imageWidth := op + vw + ip*2 + sw + op

	mask := image.NewRGBA(image.Rect(0, 0, imageWidth, h))
	draw.Draw(mask, edge.Bounds(), edge, image.ZP, draw.Src)

	sr := image.Rect(2, 0, 3, h)
	for i := 3; i <= imageWidth-3; i++ {
		dp := image.Point{i, 0}
		r := sr.Sub(sr.Min).Add(dp)
		draw.Draw(mask, r, edge, sr.Min, draw.Src)
	}

	sr = image.Rect(0, 0, 1, h)
	dp := image.Point{imageWidth - 1, 0}
	r := sr.Sub(sr.Min).Add(dp)
	draw.Draw(mask, r, edge, sr.Min, draw.Src)

	sr = image.Rect(1, 0, 2, h)
	dp = image.Point{imageWidth - 2, 0}
	r = sr.Sub(sr.Min).Add(dp)
	draw.Draw(mask, r, edge, sr.Min, draw.Src)

	img := image.NewRGBA(image.Rect(0, 0, imageWidth, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{d.Color}, image.ZP, draw.Src)

	rect := image.Rect(0, 0, op+vw+ip, h)
	draw.Draw(img, rect, &image.Uniform{Grey}, image.ZP, draw.Src)

	dst := image.NewRGBA(image.Rect(0, 0, imageWidth, h))
	draw.DrawMask(dst, dst.Bounds(), img, image.ZP, mask, image.ZP, draw.Over)

	draw.Draw(dst, gradient.Bounds(), gradient, image.ZP, draw.Over)

	gsr := image.Rect(2, 0, 3, h)
	for i := 3; i <= imageWidth-3; i++ {
		dp := image.Point{i, 0}
		gr := gsr.Sub(gsr.Min).Add(dp)
		draw.Draw(dst, gr, gradient, gsr.Min, draw.Over)
	}

	sr = image.Rect(0, 0, 1, h)
	dp = image.Point{imageWidth - 1, 0}
	r = sr.Sub(sr.Min).Add(dp)
	draw.Draw(dst, r, gradient, sr.Min, draw.Over)

	sr = image.Rect(1, 0, 2, h)
	dp = image.Point{imageWidth - 2, 0}
	r = sr.Sub(sr.Min).Add(dp)
	draw.Draw(dst, r, gradient, sr.Min, draw.Over)

	draw.Draw(dst, dst.Bounds(), v, image.Point{-op, 0}, draw.Over)

	draw.Draw(dst, dst.Bounds(), s, image.Point{-(op + vw + ip*2), 0}, draw.Over)

	png.Encode(w, dst)
}
