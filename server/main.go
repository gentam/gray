package main

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"gray"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	//go:embed index.html
	indexHTML []byte
	upgrader  = websocket.Upgrader{}
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write(indexHTML) })
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("serving on localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respCh := make(chan []byte)
	go writer(conn, respCh)

	world := makeWorld()
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		config := readConfig(p)
		fmt.Printf("canvas size: %d x %d\n", config.width, config.height)
		camera := makeCamera(config)
		go render(world, camera, respCh)
	}
}

func writer(conn *websocket.Conn, respCh <-chan []byte) {
	defer conn.Close()
	for {
		msg, ok := <-respCh
		if !ok {
			return
		}
		if err := conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
			return
		}
	}
}

type renderConfig struct {
	height, width uint32
}

func readConfig(p []byte) renderConfig {
	c := renderConfig{}
	numFields := 2
	if len(p) < numFields*4 {
		fmt.Println("invalid config payload")
		return c
	}
	c.width = binary.LittleEndian.Uint32(p[0:4])
	c.height = binary.LittleEndian.Uint32(p[4:8])
	return c
}

func pointsToBinary(pts []gray.Pixel) []byte {
	buf := make([]byte, len(pts)*7)
	for i, pt := range pts {
		offset := i * 7
		binary.LittleEndian.PutUint16(buf[offset:offset+2], uint16(pt.X))
		binary.LittleEndian.PutUint16(buf[offset+2:offset+4], uint16(pt.Y))
		buf[offset+4] = pt.R
		buf[offset+5] = pt.G
		buf[offset+6] = pt.B
	}
	return buf
}

func render(world gray.Hitter[float64], camera *gray.Camera[float64], sendCh chan []byte) {
	start := time.Now()
	pointsCh := make(chan []gray.Pixel)

	go camera.RenderStream(pointsCh, world)

	for points := range pointsCh {
		select {
		case sendCh <- pointsToBinary(points):
		default:
			fmt.Println("canceled")
			return
		}
	}
	fmt.Println("rendered in", time.Since(start))
}

func makeWorld() gray.Hitter[float64] {
	world := gray.NewHitterList[float64]()
	groundMaterial := gray.NewLambertian(rgb(0.5, 0.5, 0.5))
	world.Add(gray.NewSphere(point(0., -1000, 0), 1000, groundMaterial))

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := point(
				float64(a)+0.9*rand.Float64(),
				0.2,
				float64(b)+0.9*rand.Float64(),
			)
			if center.Subtracted(point(4, 0.2, 0)).Len() <= 0.9 {
				continue
			}

			var sphereMaterial gray.Material[float64]
			chooseMat := rand.Float64()
			switch {
			case chooseMat < 0.8: // diffuse
				albedo := gray.RandomVecIn(0.0, 1.0).Multiplied(gray.RandomVecIn(0.0, 1.0))
				sphereMaterial = gray.NewLambertian(albedo)
			case chooseMat < 0.95: // metal
				albedo := gray.RandomVecIn(0.5, 1)
				fuzz := gray.RandomFloatIn(0, 0.5)
				sphereMaterial = gray.NewMetal(albedo, fuzz)
			default: // glass
				sphereMaterial = gray.NewDielectric(1.5)
			}
			world.Add(gray.NewSphere(center, 0.2, sphereMaterial))
		}
	}

	material1 := gray.NewDielectric(1.5)
	world.Add(gray.NewSphere(point(0., 1, 0), 1.0, material1))

	material2 := gray.NewLambertian(rgb(0.4, 0.2, 0.1))
	world.Add(gray.NewSphere(point(-4., 1, 0), 1.0, material2))

	material3 := gray.NewMetal(rgb(0.7, 0.6, 0.5), 0.0)
	world.Add(gray.NewSphere(point(4., 1, 0), 1.0, material3))
	return world
}

func makeCamera(config renderConfig) *gray.Camera[float64] {
	camera := gray.NewCamera[float64]()
	camera.ImageWidth = int(config.width)
	camera.ImageHeight = int(config.height)
	camera.SamplesPerPixel = 1
	camera.MaxDepth = 50

	camera.VFOV = 20
	camera.LookFrom = point(13., 2., 3.)
	camera.LookAt = point(0., 0, 0)
	camera.VUp = point(0., 1, 0)

	camera.DefocusAngle = 0.6
	camera.FocusDistance = 10
	return camera
}

func point[T gray.Float](x, y, z T) gray.Point3[T] {
	return gray.Point3[T]{x, y, z}
}

func rgb[T gray.Float](r, g, b T) gray.RGB[T] {
	return gray.RGB[T]{r, g, b}
}
