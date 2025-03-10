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
}

func NewMetal[T Float](albedo RGB[T]) *Metal[T] {
	return &Metal[T]{Albedo: albedo}
}

func (m *Metal[T]) Scatter(r *Ray[T], rec *HitRecord[T]) (bool, *Ray[T], RGB[T]) {
	reflected := r.Direction.Reflected(rec.Normal)
	scattered := NewRay(rec.P, reflected)
	return true, scattered, m.Albedo
}
