package main

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
