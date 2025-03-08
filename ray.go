package main

type Ray[T Float] struct {
	Origin    Vec3[T]
	Direction Vec3[T]
}

func NewRay[T Float](origin, direction Vec3[T]) *Ray[T] {
	return &Ray[T]{
		Origin:    origin,
		Direction: direction,
	}
}

func (r *Ray[T]) At(t T) Vec3[T] {
	return r.Origin.Added(r.Direction.Scaled(t))
}
