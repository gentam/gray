package main

func main() {
	world := NewHitterList[float64]()
	materialGround := NewLambertian(rgb(0.8, 0.8, 0.0))
	materialCenter := NewLambertian(rgb(0.1, 0.2, 0.5))
	materialLeft := NewDielectric(1.00 / 1.33)
	materialRight := NewMetal(rgb(0.8, 0.6, 0.2), 1.0)

	world.Add(NewSphere(point(0.0, -100.5, -1), 100, materialGround))
	world.Add(NewSphere(point(0.0, 0.0, -1.2), 0.5, materialCenter))
	world.Add(NewSphere(point(-1.0, 0.0, -1.0), 0.5, materialLeft))
	world.Add(NewSphere(point(1.0, 0.0, -1.0), 0.5, materialRight))

	camera := NewCamera(800, 16.0/9.0)
	camera.SamplesPerPixel = 100
	camera.MaxDepth = 50
	camera.Render(world)
}
