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

type HitterList[T Float] struct {
	hitters []Hitter[T]
}

func NewHitterList[T Float](hitters ...Hitter[T]) *HitterList[T] {
	return &HitterList[T]{hitters: hitters}
}

func (hl *HitterList[T]) Add(hitter Hitter[T]) {
	hl.hitters = append(hl.hitters, hitter)
}

func (hl *HitterList[T]) Hit(r *Ray[T], rayT Interval[T], rec *HitRecord[T]) bool {
	tmpRec := &HitRecord[T]{}
	hitAnything := false
	closestSoFar := rayT.Max

	for _, h := range hl.hitters {
		if h.Hit(r, NewInterval(rayT.Min, closestSoFar), tmpRec) {
			hitAnything = true
			closestSoFar = tmpRec.T
			*rec = *tmpRec
		}
	}

	return hitAnything
}
