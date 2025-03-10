package main

func main() {
	world := NewHitterList[float64]()

	materialGround := NewLambertian(rgb(0.8, 0.8, 0.0))
	materialCenter := NewLambertian(rgb(0.1, 0.2, 0.5))
	materialLeft := NewDielectric(1.50)
	materialBubble := NewDielectric(1.00 / 1.50)
	materialRight := NewMetal(rgb(0.8, 0.6, 0.2), 1.0)

	world.Add(NewSphere(point(0.0, -100.5, -1.0), 100.0, materialGround))
	world.Add(NewSphere(point(0.0, 0.0, -1.2), 0.5, materialCenter))
	world.Add(NewSphere(point(-1.0, 0.0, -1.0), 0.5, materialLeft))
	world.Add(NewSphere(point(-1.0, 0.0, -1.0), 0.4, materialBubble))
	world.Add(NewSphere(point(1.0, 0.0, -1.0), 0.5, materialRight))

	camera := NewCamera[float64]()
	camera.AspectRatio = 16.0 / 9.0
	camera.ImageWidth = 800
	camera.SamplesPerPixel = 100
	camera.MaxDepth = 50

	camera.VFOV = 20
	camera.LookFrom = point(-2.0, 2, 1)
	camera.LookAt = point(0.0, 0.0, -1.0)
	camera.VUp = point(0.0, 1.0, 0.0)

	camera.DefocusAngle = 10.0
	camera.FocusDistance = 3.4

	camera.Render(world)
}
