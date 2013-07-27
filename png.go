package main

import (
	"net/http"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func makePngShield (w http.ResponseWriter, d Data) {
	w.Header().Add("content-type", "image/png");

	fi, _ := os.Open("edge.png");
	edge, _ := png.Decode(fi);

	img := image.NewRGBA(image.Rect(0, 0, 100, 19));
	blue := color.RGBA{0, 0, 255, 255};
    draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src);
	red := color.RGBA{255, 0, 0, 255};
	rect := image.Rect(0, 0, 50, 19);
    draw.Draw(img, rect, &image.Uniform{red}, image.ZP, draw.Src);


	dst := image.NewRGBA(image.Rect(0, 0, 100, 19));
	draw.DrawMask(dst, dst.Bounds(), img, image.ZP, edge, image.ZP, draw.Over);
	png.Encode(w, dst);
}
