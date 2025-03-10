package main

func main() {
	world := NewHitterList(
		NewSphere(point(0.0, 0, -1), 0.5),
		NewSphere(point(0.0, -100.5, -1), 100),
	)

	camera := NewCamera(800, 16.0/9.0, 100)
	camera.Render(world)
}

func point[T Float](x, y, z T) Point3[T] { return Point3[T]{x, y, z} }
