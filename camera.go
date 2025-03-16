package gray

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/rand/v2"
	"os"
	"runtime"
	"sync"
	"time"
)

type Camera[T Float] struct {
	AspectRatio       T   // Ratio of image width over height
	ImageWidth        int // Rendered image width in pixel count
	ImageHeight       int // Rendered image height
	SamplesPerPixel   int // Number of random samples per pixel
	pixelSamplesScale T   // Color scale factor for a sum of pixel samples
	MaxDepth          int // Maximum number of ray bounces into scene

	VFOV     T         // Vertical view angle (field of view)
	LookFrom Point3[T] // Point camera is looking from
	LookAt   Point3[T] // Point camera is looking at
	VUp      Vec3[T]   // Camera-relative "up" direction
	center   Point3[T] // Camera center
	u, v, w  Vec3[T]   // Camera frame basis vectors

	DefocusAngle  T       // Variation angle of rays through each pixel
	FocusDistance T       // Distance from camera lookfrom point to plane of perfect focus
	defocusDiskU  Vec3[T] // Defocus disk horizontal radius
	defocusDiskV  Vec3[T] // Defocus disk vertical radius

	pixel00Loc  Point3[T] // Location of pixel 0, 0
	pixelDeltaU Vec3[T]   // Offset to pixel to the right
	pixelDeltaV Vec3[T]   // Offset to pixel below
}

func NewCamera[T Float]() *Camera[T] {
	return &Camera[T]{
		AspectRatio:     1.0,
		ImageWidth:      100,
		SamplesPerPixel: 10,
		MaxDepth:        10,

		VFOV:   90,
		LookAt: Point3[T]{0, 0, -1},
		VUp:    Vec3[T]{0, 1, 0},

		FocusDistance: 10,
	}
}

func (c *Camera[T]) init() {
	if c.ImageHeight == 0 {
		c.ImageHeight = max(1, int(T(c.ImageWidth)/c.AspectRatio))
	} else {
		c.AspectRatio = T(c.ImageWidth) / T(c.ImageHeight)
	}

	c.pixelSamplesScale = T(1.0) / T(c.SamplesPerPixel)

	c.center = c.LookFrom

	// Determine viewport dimensions.
	theta := degreesToRadians(c.VFOV)
	h := T(math.Tan(float64(theta / 2)))
	viewportHeight := 2 * h * c.FocusDistance
	viewPortWidth := viewportHeight * (T(c.ImageWidth) / T(c.ImageHeight))

	// Calculate the u,v,w unit basis vectors for the camera coordinate frame.
	c.w = c.LookFrom.Subtracted(c.LookAt).Normalized()
	c.u = c.VUp.Cross(c.w).Normalized()
	c.v = c.w.Cross(c.u)

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := c.u.Scaled(viewPortWidth)   // Vector across viewport horizontal edge
	viewportV := c.v.Scaled(-viewportHeight) // Vector down viewport vertical edge

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	c.pixelDeltaU = viewportU.Divided(T(c.ImageWidth))
	c.pixelDeltaV = viewportV.Divided(T(c.ImageHeight))

	// Calculate the location of the upper left pixel.
	viewportUpperLeft := c.center.
		Subtracted(c.w.Scaled(c.FocusDistance)).
		Subtracted(viewportU.Divided(2)).
		Subtracted(viewportV.Divided(2))
	c.pixel00Loc = viewportUpperLeft.Added(c.pixelDeltaU.Added(c.pixelDeltaV).Scaled(0.5))

	// Calculate the camera defocus disk basis vectors.
	defocusRadius := c.FocusDistance * T(math.Tan(float64(degreesToRadians(c.DefocusAngle/2))))
	c.defocusDiskU = c.u.Scaled(defocusRadius)
	c.defocusDiskV = c.v.Scaled(defocusRadius)
}

func (c *Camera[T]) RenderPNG(w io.Writer, world Hitter[T]) {
	c.init()
	rect := image.Rect(0, 0, c.ImageWidth, c.ImageHeight)
	img := image.NewRGBA(rect)

	type result struct {
		j   int
		row []color.RGBA
	}
	resultCh := make(chan result)

	for j := range c.ImageHeight {
		go func() {
			row := make([]color.RGBA, c.ImageWidth)
			for i := range c.ImageWidth {
				pixelColor := RGB[T]{}
				for range c.SamplesPerPixel {
					r := c.getRay(i, j)
					pixelColor.Add(c.rayColor(r, c.MaxDepth, world))
				}
				rgba := pixelColor.Scaled(c.pixelSamplesScale).RGBA()
				row[i] = rgba
			}
			resultCh <- result{j: j, row: row}
		}()
	}

	for j := range c.ImageHeight {
		fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d ", c.ImageHeight-j)
		res := <-resultCh
		for i := range c.ImageWidth {
			img.Set(i, res.j, res.row[i])
		}
	}

	if err := png.Encode(w, img); err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stderr, "\rDone.                 ")
}

type Pixel struct {
	X, Y    int
	R, B, G uint8
}

// RenderStream closes the streamCh on finish or when the context is canceled.
func (c *Camera[T]) RenderStream(ctx context.Context, streamCh chan<- []Pixel, world Hitter[T]) {
	if ctx.Err() != nil {
		return
	}
	c.init()

	// workers
	queue := make(chan [2]int)
	bufCh := make(chan Pixel)
	wg := sync.WaitGroup{}
	numCPU := runtime.NumCPU()
	wg.Add(numCPU)
	for range numCPU {
		go func() {
			defer wg.Done()
			for pos := range queue {
				select {
				case <-ctx.Done():
					return
				default:
				}

				i, j := pos[0], pos[1]
				pixelColor := RGB[T]{}
				for range c.SamplesPerPixel {
					r := c.getRay(i, j)
					pixelColor.Add(c.rayColor(r, c.MaxDepth, world))
				}
				r, g, b := pixelColor.Scaled(c.pixelSamplesScale).RGB()
				bufCh <- Pixel{X: i, Y: j, R: r, G: g, B: b}
			}
		}()
	}

	n := c.ImageWidth * c.ImageHeight
	pts := make([][2]int, n)
	for j := range c.ImageHeight {
		for i := range c.ImageWidth {
			pts[j*c.ImageWidth+i] = [2]int{i, j}
		}
	}
	rand.Shuffle(n, func(i int, j int) {
		pts[i], pts[j] = pts[j], pts[i]
	})
	go func() { // work distribution
		for _, pt := range pts {
			select {
			case <-ctx.Done():
				goto cleanup
			default:
			}
			queue <- pt
		}
	cleanup:
		close(queue)
		wg.Wait()
		close(bufCh)
	}()

	// buffer pixels for the interval and stream in batches
	const bufSize = 1 << 12
	buf := make([]Pixel, 0, bufSize)
	const interval = 10 * time.Millisecond
	timer := time.NewTimer(interval)
	for i := n; i > 0; {
		select {
		case <-ctx.Done():
			goto cleanup
		case pt, ok := <-bufCh:
			if !ok {
				goto cleanup
			}
			buf = append(buf, pt)
		case <-timer.C:
			if len(buf) > 0 {
				streamCh <- buf
				i -= len(buf)
				// NOTE: buf[:0] will cause data race
				buf = make([]Pixel, 0, bufSize)
			}
			timer.Reset(interval)
		}
	}
cleanup:
	timer.Stop()
	if len(buf) > 0 {
		streamCh <- buf
	}
	close(streamCh)
}

func (c *Camera[T]) rayColor(r *Ray[T], depth int, world Hitter[T]) RGB[T] {
	// If we've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return RGB[T]{}
	}

	rec := &HitRecord[T]{}
	// 0.001: ignore very close hits to fix shadow acne
	if world.Hit(r, NewInterval(0.001, T(math.Inf(1))), rec) {
		if ok, scattered, attenuation := rec.Material.Scatter(r, rec); ok {
			return c.rayColor(scattered, depth-1, world).Multiplied(attenuation)
		}
		return RGB[T]{}
	}

	unitDirection := r.Direction.Normalized()
	a := 0.5 * (unitDirection.Y() + 1.0)
	white := RGB[T]{1.0, 1.0, 1.0}
	blue := RGB[T]{0.5, 0.7, 1.0}
	return white.Scaled(1.0 - a).Added(blue.Scaled(a))
}

// getRay returns a camera ray originating from the defocus disk and directed at
// randomly sampled point around the pixel location i, j.
func (c *Camera[T]) getRay(i, j int) *Ray[T] {
	offset := c.sampleSquare()
	pixelSample := c.pixel00Loc.
		Added(c.pixelDeltaU.Scaled(T(i) + offset.X())).
		Added(c.pixelDeltaV.Scaled(T(j) + offset.Y()))
	rayOrigin := c.center
	if c.DefocusAngle > 0 {
		rayOrigin = c.defocusDiskSample()
	}
	rayDirection := pixelSample.Subtracted(rayOrigin)

	return NewRay(rayOrigin, rayDirection)
}

// sampleSquare returns the vector to a random point in the
// [-.5,-.5]-[+.5,+.5] unit square.
func (c *Camera[T]) sampleSquare() Vec3[T] {
	return Vec3[T]{
		T(rand.Float64() - 0.5),
		T(rand.Float64() - 0.5),
		0,
	}
}

// defocusDiskSample returns a random point in the camera defocus disk.
func (c *Camera[T]) defocusDiskSample() Vec3[T] {
	p := randomInUnitDisk[T]()
	return c.center.
		Added(c.defocusDiskU.Scaled(p[0])).
		Added(c.defocusDiskV.Scaled(p[1]))
}
