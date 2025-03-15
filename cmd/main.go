package main

import (
	"gray"
	"math/rand"
)

func main() {
	world := gray.NewHitterList[float64]()

	groundMaterial := gray.NewLambertian(rgb(0.5, 0.5, 0.5))
	world.Add(gray.NewSphere(point(0., -1000, 0), 1000, groundMaterial))

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := point(
				float64(a)+0.9*rand.Float64(),
				0.2,
				float64(b)+0.9*rand.Float64(),
			)
			if center.Subtracted(point(4, 0.2, 0)).Len() <= 0.9 {
				continue
			}

			var sphereMaterial gray.Material[float64]
			chooseMat := rand.Float64()
			switch {
			case chooseMat < 0.8: // diffuse
				albedo := gray.RandomVecIn(0.0, 1.0).Multiplied(gray.RandomVecIn(0.0, 1.0))
				sphereMaterial = gray.NewLambertian(albedo)
			case chooseMat < 0.95: // metal
				albedo := gray.RandomVecIn(0.5, 1)
				fuzz := gray.RandomFloatIn(0, 0.5)
				sphereMaterial = gray.NewMetal(albedo, fuzz)
			default: // glass
				sphereMaterial = gray.NewDielectric(1.5)
			}
			world.Add(gray.NewSphere(center, 0.2, sphereMaterial))
		}
	}

	material1 := gray.NewDielectric(1.5)
	world.Add(gray.NewSphere(point(0., 1, 0), 1.0, material1))

	material2 := gray.NewLambertian(rgb(0.4, 0.2, 0.1))
	world.Add(gray.NewSphere(point(-4., 1, 0), 1.0, material2))

	material3 := gray.NewMetal(rgb(0.7, 0.6, 0.5), 0.0)
	world.Add(gray.NewSphere(point(4., 1, 0), 1.0, material3))

	camera := gray.NewCamera[float64]()
	camera.AspectRatio = 16.0 / 9.0
	camera.ImageWidth = 1200
	camera.SamplesPerPixel = 500
	camera.MaxDepth = 50

	camera.VFOV = 20
	camera.LookFrom = point(13., 2., 3.)
	camera.LookAt = point(0., 0, 0)
	camera.VUp = point(0., 1, 0)

	camera.DefocusAngle = 0.6
	camera.FocusDistance = 10

	camera.Render(world)
}

func point[T gray.Float](x, y, z T) gray.Point3[T] {
	return gray.Point3[T]{x, y, z}
}

func rgb[T gray.Float](r, g, b T) gray.RGB[T] {
	return gray.RGB[T]{r, g, b}
}
