package main

import "math"

type Float interface {
	~float32 | ~float64
}

func point[T Float](x, y, z T) Point3[T] {
	return Point3[T]{x, y, z}
}

func sqrt[T Float](x T) T {
	return T(math.Sqrt(float64(x)))
}

func linearToGamma[T Float](x T) T {
	if x <= 0 {
		return T(0)
	}
	return T(math.Sqrt(float64(x)))
}
