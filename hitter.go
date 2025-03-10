package main

type Hitter[T Float] interface {
	Hit(r *Ray[T], rayT Interval[T], hitRecord *HitRecord[T]) bool
}

type HitRecord[T Float] struct {
	P         Point3[T]
	Normal    Vec3[T]
	Material  Material[T]
	T         T
	FrontFace bool
}

// outwardNormal is assumed to be have unit length
func (rec *HitRecord[T]) SetFaceNormal(r *Ray[T], outwardNormal Vec3[T]) {
	rec.FrontFace = r.Direction.Dot(outwardNormal) < 0
	rec.Normal = outwardNormal
	if !rec.FrontFace {
		rec.Normal.Negate()
	}
}
