package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
)

type Camera[T Float] struct {
	aspectRatio T         // Ratio of image width over height
	imageWidth  int       // Rendered image width in pixel coun
	imageHeight int       // Rendered image height
	center      Point3[T] // Camera center
	pixel00Loc  Point3[T] // Location of pixel 0, 0
	pixelDeltaU Vec3[T]   // Offset to pixel to the right
	pixelDeltaV Vec3[T]   // Offset to pixel below
}

func NewCamera[T Float](width int, aspectRatio T) *Camera[T] {
	c := &Camera[T]{}
	c.imageWidth = width
	c.aspectRatio = aspectRatio
	c.imageHeight = max(1, int(T(c.imageWidth)/c.aspectRatio))

	c.center = Vec3[T]{0, 0, 0}

	// Determine viewport dimensions.
	focalLength := T(1.0)
	viewportHeight := T(2.0)
	viewPortWidth := viewportHeight * (T(c.imageWidth) / T(c.imageHeight))

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := Vec3[T]{viewPortWidth, 0, 0}
	viewportV := Vec3[T]{0, -viewportHeight, 0}

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	c.pixelDeltaU = viewportU.Divided(T(c.imageWidth))
	c.pixelDeltaV = viewportV.Divided(T(c.imageHeight))

	// Calculate the location of the upper left pixel.
	viewportUpperLeft := c.center.
		Subtracted(Vec3[T]{0, 0, focalLength}).
		Subtracted(viewportU.Divided(2)).
		Subtracted(viewportV.Divided(2))
	c.pixel00Loc = viewportUpperLeft.Added(c.pixelDeltaU.Added(c.pixelDeltaV).Scaled(0.5))

	return c
}

func (c *Camera[T]) rayColor(r *Ray[T], world Hitter[T]) RGB[T] {
	rec := &HitRecord[T]{}
	if world.Hit(r, NewInterval(0, T(math.Inf(1))), rec) {
		return rec.Normal.Added(rgb[T](1, 1, 1)).Scaled(0.5)
	}

	unitDirection := r.Direction.Normalized()
	a := 0.5 * (unitDirection.Y() + 1.0)
	white := rgb[T](1.0, 1.0, 1.0)
	blue := rgb[T](0.5, 0.7, 1.0)
	return white.Scaled(1.0 - a).Added(blue.Scaled(a))
}

func (c *Camera[T]) Render(world Hitter[T]) {
	rect := image.Rect(0, 0, c.imageWidth, c.imageHeight)
	img := image.NewRGBA(rect)

	for j := range c.imageHeight {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", c.imageHeight-j)
		for i := range c.imageWidth {
			pixelCenter := c.pixel00Loc.
				Added(c.pixelDeltaU.Scaled(T(i))).
				Added(c.pixelDeltaV.Scaled(T(j)))
			rayDirection := pixelCenter.Subtracted(c.center)

			ray := NewRay(c.center, rayDirection)
			pixelColor := c.rayColor(ray, world)
			img.Set(i, j, pixelColor.RGBA())
		}
	}

	if err := png.Encode(os.Stdout, img); err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stderr, "\rDone.                 ")
}
