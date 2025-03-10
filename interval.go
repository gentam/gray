package main

type Interval[T Float] struct {
	Min, Max T
}

func NewInterval[T Float](min, max T) Interval[T] {
	return Interval[T]{min, max}
}

func (i Interval[T]) Size() T {
	return i.Max - i.Min
}

func (i Interval[T]) Contains(x T) bool {
	return i.Min <= x && x <= i.Max
}

func (i Interval[T]) Surrounds(x T) bool {
	return i.Min < x && x < i.Max
}

func (i Interval[T]) Clamp(x T) T {
	if x < i.Min {
		return i.Min
	}
	if x > i.Max {
		return i.Max
	}
	return x
}
