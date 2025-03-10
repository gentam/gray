package main

import (
	"image/color"
	"math"
)

type Vec3[T Float] [3]T

type Float interface {
	~float32 | ~float64
}

func NewVec3[T Float](x, y, z T) Vec3[T] {
	return Vec3[T]{x, y, z}
}

func (v Vec3[T]) X() T { return v[0] }
func (v Vec3[T]) Y() T { return v[1] }
func (v Vec3[T]) Z() T { return v[2] }

func (v *Vec3[T]) Negate() *Vec3[T] {
	v[0] = -v[0]
	v[1] = -v[1]
	v[2] = -v[2]
	return v
}
func (v Vec3[T]) Negated() Vec3[T] {
	return Vec3[T]{-v[0], -v[1], -v[2]}
}

func (v *Vec3[T]) Add(u Vec3[T]) *Vec3[T] {
	v[0] += u[0]
	v[1] += u[1]
	v[2] += u[2]
	return v
}
func (v Vec3[T]) Added(u Vec3[T]) Vec3[T] {
	return Vec3[T]{v[0] + u[0], v[1] + u[1], v[2] + u[2]}
}

func (v *Vec3[T]) Sub(u Vec3[T]) *Vec3[T] {
	v[0] -= u[0]
	v[1] -= u[1]
	v[2] -= u[2]
	return v
}
func (v Vec3[T]) Subtracted(u Vec3[T]) Vec3[T] {
	return Vec3[T]{v[0] - u[0], v[1] - u[1], v[2] - u[2]}
}

func (v *Vec3[T]) Scale(s T) *Vec3[T] {
	v[0] *= s
	v[1] *= s
	v[2] *= s
	return v
}
func (v Vec3[T]) Scaled(s T) Vec3[T] {
	return Vec3[T]{v[0] * s, v[1] * s, v[2] * s}
}

func (v *Vec3[T]) Div(s T) *Vec3[T] {
	return v.Scale(1 / s)
}
func (v Vec3[T]) Divided(s T) Vec3[T] {
	return v.Scaled(1 / s)
}

func (v Vec3[T]) Len() T {
	return T(math.Sqrt(float64(v.LenSq())))
}
func (v Vec3[T]) LenSq() T {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

func (v Vec3[T]) Dot(u Vec3[T]) T {
	return v[0]*u[0] + v[1]*u[1] + v[2]*u[2]
}

func (v Vec3[T]) Cross(u Vec3[T]) Vec3[T] {
	return Vec3[T]{
		v[1]*u[2] - v[2]*u[1],
		v[2]*u[0] - v[0]*u[2],
		v[0]*u[1] - v[1]*u[0],
	}
}

func (v *Vec3[T]) Normalize() *Vec3[T] {
	return v.Div(v.Len())
}
func (v Vec3[T]) Normalized() Vec3[T] {
	return v.Divided(v.Len())
}

// ----------------------------------------------------------------------------

type Point3[T Float] = Vec3[T]

type RGB[T Float] = Vec3[T]

func (v Vec3[T]) RGBA() color.RGBA {
	// [0,1] → [0,255]
	r, g, b := v.X(), v.Y(), v.Z()
	intensity := NewInterval[T](0.000, 0.999)
	return color.RGBA{
		uint8(255.999 * intensity.Clamp(r)),
		uint8(255.999 * intensity.Clamp(g)),
		uint8(255.999 * intensity.Clamp(b)),
		255,
	}
}
