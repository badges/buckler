package shield

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/raster"
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
	c        *freetype.Context
)

const (
	h  = 18
	op = 4
	ip = 4
)

func Init(dataPath string) {
	fi, err := os.Open(filepath.Join(dataPath, "edge.png"))
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	edge, err = png.Decode(fi)
	if err != nil {
		log.Fatal(err)
	}

	fi, err = os.Open(filepath.Join(dataPath, "gradient.png"))
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	gradient, err = png.Decode(fi)
	if err != nil {
		log.Fatal(err)
	}

	fontBytes, err := ioutil.ReadFile(filepath.Join(dataPath, "opensanssemibold.ttf"))
	if err != nil {
		log.Fatal(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	c = freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(10)
}

func hexColor(c string) (color.RGBA, bool) {
	if len(c) != 6 {
		return color.RGBA{}, false
	}

	r, rerr := strconv.ParseInt(c[0:2], 16, 16)
	g, gerr := strconv.ParseInt(c[2:4], 16, 16)
	b, berr := strconv.ParseInt(c[4:6], 16, 16)

	if rerr != nil || gerr != nil || berr != nil {
		return color.RGBA{}, false
	}

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}, true
}

func GetColor(cs string) (c color.RGBA, err error) {
	c, ok := Colors[cs]
	if !ok {
		c, ok = hexColor(cs)
		if !ok {
			err = errors.New("Unknown colour")
			return
		}
	}

	return
}

func getTextOffset(pt raster.Point) int {
	return int(math.Floor(float64(float32(pt.X)/256 + 0.5)))
}

func renderString(s string, c *freetype.Context) (*image.RGBA, int) {
	estWidth := 8 * len(s)
	dst := image.NewRGBA(image.Rect(0, 0, estWidth, h))

	c.SetDst(dst)
	c.SetClip(dst.Bounds())

	c.SetSrc(&image.Uniform{C: shadow})
	pt := freetype.Pt(0, 13)
	c.DrawString(s, pt)

	c.SetSrc(image.White)

	pt = freetype.Pt(0, 12)
	pt, _ = c.DrawString(s, pt)

	return dst, getTextOffset(pt)
}

func buildMask(mask *image.RGBA, imageWidth int, tmpl image.Image, imgOp draw.Op) {
	draw.Draw(mask, tmpl.Bounds(), tmpl, image.ZP, imgOp)

	sr := image.Rect(2, 0, 3, h)
	for i := 3; i <= imageWidth-3; i++ {
		dp := image.Point{i, 0}
		r := sr.Sub(sr.Min).Add(dp)
		draw.Draw(mask, r, tmpl, sr.Min, imgOp)
	}

	sr = image.Rect(0, 0, 1, h)
	dp := image.Point{imageWidth - 1, 0}
	r := sr.Sub(sr.Min).Add(dp)
	draw.Draw(mask, r, tmpl, sr.Min, imgOp)

	sr = image.Rect(1, 0, 2, h)
	dp = image.Point{imageWidth - 2, 0}
	r = sr.Sub(sr.Min).Add(dp)
	draw.Draw(mask, r, tmpl, sr.Min, imgOp)
}

func PNG(w io.Writer, d Data) {
	// render text to determine how wide the image has to be
	// we leave 6 pixels at the start and end, and 3 for each in the middle
	v, vw := renderString(d.Vendor, c)
	s, sw := renderString(d.Status, c)
	imageWidth := op + vw + ip*2 + sw + op

	img := image.NewRGBA(image.Rect(0, 0, imageWidth, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: d.Color}, image.ZP, draw.Src)

	rect := image.Rect(0, 0, op+vw+ip, h)
	draw.Draw(img, rect, &image.Uniform{C: Grey}, image.ZP, draw.Src)

	dst := image.NewRGBA(image.Rect(0, 0, imageWidth, h))

	mask := image.NewRGBA(image.Rect(0, 0, imageWidth, h))
	buildMask(mask, imageWidth, edge, draw.Src)
	draw.DrawMask(dst, dst.Bounds(), img, image.ZP, mask, image.ZP, draw.Over)

	buildMask(dst, imageWidth, gradient, draw.Over)

	draw.Draw(dst, dst.Bounds(), v, image.Point{-op, 0}, draw.Over)

	draw.Draw(dst, dst.Bounds(), s, image.Point{-(op + vw + ip*2), 0}, draw.Over)

	png.Encode(w, dst)
}
