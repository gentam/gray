package main

import "math/rand"

func main() {
	world := NewHitterList[float64]()

	groundMaterial := NewLambertian(rgb(0.5, 0.5, 0.5))
	world.Add(NewSphere(point(0., -1000, 0), 1000, groundMaterial))

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

			var sphereMaterial Material[float64]
			chooseMat := rand.Float64()
			switch {
			case chooseMat < 0.8: // diffuse
				albedo := randomVecIn(0.0, 1.0).Multiplied(randomVecIn(0.0, 1.0))
				sphereMaterial = NewLambertian(albedo)
			case chooseMat < 0.95: // metal
				albedo := randomVecIn(0.5, 1)
				fuzz := randomFloatIn(0, 0.5)
				sphereMaterial = NewMetal(albedo, fuzz)
			default: // glass
				sphereMaterial = NewDielectric(1.5)
			}
			world.Add(NewSphere(center, 0.2, sphereMaterial))
		}
	}

	material1 := NewDielectric(1.5)
	world.Add(NewSphere(point(0., 1, 0), 1.0, material1))

	material2 := NewLambertian(rgb(0.4, 0.2, 0.1))
	world.Add(NewSphere(point(-4., 1, 0), 1.0, material2))

	material3 := NewMetal(rgb(0.7, 0.6, 0.5), 0.0)
	world.Add(NewSphere(point(4., 1, 0), 1.0, material3))

	camera := NewCamera[float64]()
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
