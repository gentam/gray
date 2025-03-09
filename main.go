package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
)

func main() {
	// Image
	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	imageHeight := int(float64(imageWidth) / aspectRatio)
	if imageHeight == 0 {
		imageHeight = 1
	}

	// Camera
	focalLength := 1.0
	viewportHeight := 2.0
	viewPortWidth := viewportHeight * (float64(imageWidth) / float64(imageHeight))
	cameraCenter := point(0.0, 0.0, 0.0)

	viewportU := point(viewPortWidth, 0, 0)
	viewportV := point(0, -viewportHeight, 0)

	pixelDeltaU := viewportU.Divided(float64(imageWidth))
	pixelDeltaV := viewportV.Divided(float64(imageHeight))

	viewportUpperLeft := cameraCenter.
		Subtracted(point(0, 0, focalLength)).
		Subtracted(viewportU.Divided(2)).
		Subtracted(viewportV.Divided(2))
	pixel00Loc := viewportUpperLeft.Added(pixelDeltaU.Added(pixelDeltaV).Scaled(0.5))

	// Render
	rect := image.Rect(0, 0, imageWidth, imageHeight)
	img := image.NewRGBA(rect)
	for j := range imageHeight {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", imageHeight-j)
		for i := range imageWidth {
			pixelCenter := pixel00Loc.
				Added(pixelDeltaU.Scaled(float64(i))).
				Added(pixelDeltaV.Scaled(float64(j)))
			rayDirection := pixelCenter.Subtracted(cameraCenter)

			ray := NewRay(cameraCenter, rayDirection)
			pixelColor := rayColor(ray)
			img.Set(i, j, pixelColor.RGBA())
		}
	}
	if err := png.Encode(os.Stdout, img); err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stderr, "\rDone.                 ")
}

func point[T Float](x, y, z T) Point3[T] { return Point3[T]{x, y, z} }
func rgb[T Float](r, g, b T) RGB[T]      { return RGB[T]{r, g, b} }

func rayColor[T Float](r *Ray[T]) RGB[T] {
	c := point[T](0, 0, -1)
	if t := hitSphere(r, c, 0.5); t > 0 {
		n := r.At(t).Subtracted(c).Normalized()
		return rgb(n.X()+1, n.Y()+1, n.Z()+1).Scaled(0.5)
	}

	unitDirection := r.Direction.Normalized()
	a := 0.5 * (unitDirection.Y() + 1.0)
	white := rgb[T](1.0, 1.0, 1.0)
	blue := rgb[T](0.5, 0.7, 1.0)
	return white.Scaled(1.0 - a).Added(blue.Scaled(a))
}

func hitSphere[T Float](r *Ray[T], center Point3[T], radius T) T {
	oc := center.Subtracted(r.Origin)
	a := r.Direction.LenSq()
	h := r.Direction.Dot(oc)
	c := oc.LenSq() - radius*radius
	discriminant := h*h - a*c

	if discriminant < 0 {
		return -1.0
	}
	return (h - T(math.Sqrt(float64(discriminant)))) / a
}
