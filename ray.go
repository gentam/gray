package main

type Ray[T Float] struct {
	Origin    Point3[T]
	Direction Vec3[T]
}

func NewRay[T Float](origin Point3[T], direction Vec3[T]) *Ray[T] {
	return &Ray[T]{
		Origin:    origin,
		Direction: direction,
	}
}

func (r *Ray[T]) At(t T) Point3[T] {
	return r.Origin.Added(r.Direction.Scaled(t))
}
