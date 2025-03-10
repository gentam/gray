package main

type Sphere[T Float] struct {
	Center Point3[T]
	Radius T
}

func NewSphere[T Float](center Point3[T], radius T) *Sphere[T] {
	return &Sphere[T]{Center: center, Radius: max(T(0), radius)}
}

func (s *Sphere[T]) Hit(r *Ray[T], rayT Interval[T], rec *HitRecord[T]) bool {
	oc := s.Center.Subtracted(r.Origin)
	a := r.Direction.LenSq()
	h := r.Direction.Dot(oc)
	c := oc.LenSq() - s.Radius*s.Radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return false
	}

	sqrtd := sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range
	root := (h - sqrtd) / a
	if !rayT.Surrounds(root) {
		root = (h + sqrtd) / a
		if !rayT.Surrounds(root) {
			return false
		}
	}

	rec.T = root
	rec.P = r.At(rec.T)
	outwardNormal := rec.P.Subtracted(s.Center).Divided(s.Radius)
	rec.SetFaceNormal(r, outwardNormal)

	return true
}
