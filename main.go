package main

import "math"

func main() {
	world := NewHitterList[float64]()
	r := math.Cos(math.Pi / 4)

	materialLeft := NewLambertian(rgb(0., 0, 1))
	materialRight := NewLambertian(rgb(1., 0, 0))

	world.Add(NewSphere(point(-r, 0, -1), r, materialLeft))
	world.Add(NewSphere(point(r, 0, -1), r, materialRight))

	camera := NewCamera(800, 16.0/9.0)
	camera.SamplesPerPixel = 100
	camera.MaxDepth = 50
	camera.VFOV = 90
	camera.Render(world)
}
