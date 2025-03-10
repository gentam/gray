package main

type Material[T Float] interface {
	Scatter(r *Ray[T], rec *HitRecord[T]) (ok bool, scattered *Ray[T], attenuation RGB[T])
}

type Lambertian[T Float] struct {
	Albedo RGB[T]
}

func NewLambertian[T Float](albedo RGB[T]) *Lambertian[T] {
	return &Lambertian[T]{Albedo: albedo}
}

func (l *Lambertian[T]) Scatter(r *Ray[T], rec *HitRecord[T]) (bool, *Ray[T], RGB[T]) {
	scatterDirection := rec.Normal.Added(randomUnitVec[T]())

	// Catch degenerate scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = rec.Normal
	}

	scattered := NewRay(rec.P, scatterDirection)
	return true, scattered, l.Albedo
}

type Metal[T Float] struct {
	Albedo RGB[T]
	Fuzz   T
}

func NewMetal[T Float](albedo RGB[T], fuzz T) *Metal[T] {
	return &Metal[T]{Albedo: albedo, Fuzz: min(fuzz, 1)}
}

func (m *Metal[T]) Scatter(r *Ray[T], rec *HitRecord[T]) (bool, *Ray[T], RGB[T]) {
	reflected := r.Direction.Reflected(rec.Normal)
	reflected = reflected.Normalized().Added(randomUnitVec[T]().Scaled(m.Fuzz))
	scattered := NewRay(rec.P, reflected)
	return scattered.Direction.Dot(rec.Normal) > 0, scattered, m.Albedo
}
