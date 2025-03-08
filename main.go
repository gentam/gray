package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func main() {
	const width = 256
	const height = 256
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)
	for j := range height {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", height-j)
		for i := range width {
			r := float64(i) / (width - 1)
			g := float64(j) / (height - 1)
			b := 0.0
			v := NewVec3(r, g, b)
			img.Set(i, j, v.Color())
		}
	}
	if err := png.Encode(os.Stdout, img); err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stderr, "\rDone.                 ")
}
