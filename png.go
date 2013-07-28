package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"code.google.com/p/freetype-go/freetype"
)

type Data struct {
	Vendor string
	Status string
	Color  color.RGBA
}

var (
	Gray        = color.RGBA{74, 74, 74, 255}
	BrightGreen = color.RGBA{69, 203, 20, 255}
	Green       = color.RGBA{124, 166, 0, 255}
	YellowGreen = color.RGBA{156, 158, 9, 255}
	Yellow      = color.RGBA{184, 148, 19, 255}
	Orange      = color.RGBA{184, 113, 37, 255}
	Red         = color.RGBA{186, 77, 56, 255}
	LightGray   = color.RGBA{131, 131, 131, 255}
	Blue        = color.RGBA{0, 126, 198, 255}

	Colors = map[string]color.RGBA{
		"gray":        Gray,
		"brightgreen": BrightGreen,
		"green":       Green,
		"yellowgreen": YellowGreen,
		"yellow":      Yellow,
		"orange":      Orange,
		"red":         Red,
		"lightgray":   LightGray,
		"blue":        Blue,
	}
)

const (
	h = 18
)

func makePngShield(w http.ResponseWriter, d Data) {
	w.Header().Add("content-type", "image/png")

	fi, _ := os.Open("edge.png")
	edge, _ := png.Decode(fi)
	defer fi.Close()

	fi, _ = os.Open("gradient.png")
	gradient, _ := png.Decode(fi)
	defer fi.Close()

	mask := image.NewRGBA(image.Rect(0, 0, 100, h))
	draw.Draw(mask, edge.Bounds(), edge, image.ZP, draw.Src)

	sr := image.Rect(2, 0, 3, h)
	for i := 3; i <= 97; i++ {
		dp := image.Point{i, 0}
		r := sr.Sub(sr.Min).Add(dp)
		draw.Draw(mask, r, edge, sr.Min, draw.Src)
	}

	sr = image.Rect(0, 0, 1, h)
	dp := image.Point{99, 0}
	r := sr.Sub(sr.Min).Add(dp)
	draw.Draw(mask, r, edge, sr.Min, draw.Src)

	sr = image.Rect(1, 0, 2, h)
	dp = image.Point{98, 0}
	r = sr.Sub(sr.Min).Add(dp)
	draw.Draw(mask, r, edge, sr.Min, draw.Src)

	img := image.NewRGBA(image.Rect(0, 0, 100, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{d.Color}, image.ZP, draw.Src)

	rect := image.Rect(0, 0, 50, h)
	draw.Draw(img, rect, &image.Uniform{Gray}, image.ZP, draw.Src)

	dst := image.NewRGBA(image.Rect(0, 0, 100, h))
	draw.DrawMask(dst, dst.Bounds(), img, image.ZP, mask, image.ZP, draw.Over)

	draw.Draw(dst, gradient.Bounds(), gradient, image.ZP, draw.Over)

	gsr := image.Rect(2, 0, 3, h)
	for i := 3; i <= 97; i++ {
		dp := image.Point{i, 0}
		gr := gsr.Sub(gsr.Min).Add(dp)
		draw.Draw(dst, gr, gradient, gsr.Min, draw.Over)
	}

	sr = image.Rect(0, 0, 1, h)
	dp = image.Point{99, 0}
	r = sr.Sub(sr.Min).Add(dp)
	draw.Draw(dst, r, gradient, sr.Min, draw.Over)

	sr = image.Rect(1, 0, 2, h)
	dp = image.Point{98, 0}
	r = sr.Sub(sr.Min).Add(dp)
	draw.Draw(dst, r, gradient, sr.Min, draw.Over)

	fontBytes, err := ioutil.ReadFile("opensanssemibold.ttf")
	if err != nil {
		log.Println(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(10)
	c.SetDst(dst)
	c.SetClip(dst.Bounds())

	shadow := color.RGBA{0, 0, 0, 125}
	c.SetSrc(&image.Uniform{shadow})
	pt := freetype.Pt(6, 13)
	offset, _ := c.DrawString(d.Vendor, pt)

	pt = freetype.Pt(53, 13)
	c.DrawString(d.Status, pt)

	c.SetSrc(image.White)

	pt = freetype.Pt(6, 12)
	offset, _ = c.DrawString(d.Vendor, pt)

	pt = freetype.Pt(53, 12)
	c.DrawString(d.Status, pt)

	println(offset.X, offset.Y)
	png.Encode(w, dst)
}
