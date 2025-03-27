package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"time"

	_ "embed"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"gonum.org/v1/gonum/mat"
)

//go:embed particle.png
var particlePng []byte

func main() {
	pixelgl.Run(run)
}

type Particle struct {
	Pos pixel.Vec
	Dir pixel.Vec
}

func run() {
	screenSize := pixel.V(800, 800)
	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, screenSize.X, screenSize.Y),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	speed := 100.0

	particles := make([]Particle, 20000)

	for i := range particles {
		particles[i].Pos = screenSize.Scaled(0.5).Add(pixel.V(rand.Float64()*200, 0).Rotated(rand.Float64() * 2 * math.Pi))
		particles[i].Dir = pixel.V(0, 1).Rotated(rand.Float64() * math.Pi * 2)
	}

	particleImg, err := png.Decode(bytes.NewReader(particlePng))
	if err != nil {
		panic(err)
	}
	particlePic := pixel.PictureDataFromImage(particleImg)

	particleBatch := pixel.NewBatch(&pixel.TrianglesData{}, particlePic)
	particleSprite := pixel.NewSprite(particlePic, particlePic.Rect)

	updatingTime := time.Duration(0)
	renderingTime := time.Duration(0)

	grid := NewGrid(int(screenSize.X), int(screenSize.Y), 4)

	step := 0

	for !win.Closed() {
		win.Update()

		startUpdating := time.Now()
		grid.Decay(0.99)
		for i := range particles {

			sensorLeftPosition := particles[i].Pos.Add(particles[i].Dir.Scaled(20).Rotated(-math.Pi / 5))
			sensorRightPosition := particles[i].Pos.Add(particles[i].Dir.Scaled(20).Rotated(math.Pi / 5))

			sensorLeftReading := grid.Sense(sensorLeftPosition, 0.1)
			sensorRightReading := grid.Sense(sensorRightPosition, 0.1)

			if sensorLeftReading > sensorRightReading {
				particles[i].Dir = particles[i].Dir.Rotated(0.4 * -math.Pi / 60.0)
			} else if sensorLeftReading < sensorRightReading {
				particles[i].Dir = particles[i].Dir.Rotated(0.4 * math.Pi / 60.0)
			}

			particles[i].Pos = particles[i].Pos.Add(particles[i].Dir.Scaled(1.0 / 60.0 * speed))

			if particles[i].Pos.X < 0 && particles[i].Dir.X < 0 {
				particles[i].Dir.X *= -1
			}
			if particles[i].Pos.Y < 0 && particles[i].Dir.Y < 0 {
				particles[i].Dir.Y *= -1
			}
			if particles[i].Pos.X > screenSize.X && particles[i].Dir.X > 0 {
				particles[i].Dir.X *= -1
			}
			if particles[i].Pos.Y > screenSize.Y && particles[i].Dir.Y > 0 {
				particles[i].Dir.Y *= -1
			}

			grid.AddPheremone(particles[i].Pos, 1.0/60)
		}
		updatingTime += time.Since(startUpdating)

		startRendering := time.Now()
		bgSprite, bgMat := grid.Sprite()

		particleBatch.Clear()
		for _, p := range particles {
			particleSprite.Draw(particleBatch, pixel.IM.Scaled(pixel.ZV, 0.05).Moved(p.Pos))
		}

		win.Clear(pixel.RGB(0.1, 0.1, 0.1))
		bgSprite.Draw(win, bgMat.Moved(pixel.V(100, 100)))
		particleBatch.Draw(win)
		renderingTime += time.Since(startRendering)

		if step%120 == 0 {
			fmt.Println("per frame:", "update", updatingTime/120, "render", renderingTime/120)
			updatingTime = 0
			renderingTime = 0
		}
		step++
	}
}

type Grid struct {
	chunkSie int
	grid     *mat.Dense
}

func NewGrid(sizeX, sizeY, chunkSize int) *Grid {
	if (sizeX/chunkSize)*chunkSize != sizeX || (sizeY/chunkSize)*chunkSize != sizeY {
		panic("not divisible")
	}

	return &Grid{
		chunkSie: chunkSize,
		grid:     mat.NewDense(sizeX/chunkSize+1, sizeY/chunkSize+1, nil),
	}
}

func (g *Grid) Decay(mult float64) {
	g.grid.Scale(mult, g.grid)
}

func (g *Grid) AddPheremone(at pixel.Vec, amount float64) {
	gpx := int(math.Round(at.X / float64(g.chunkSie)))
	gpy := int(math.Round(at.Y / float64(g.chunkSie)))
	if gpx >= 0 && gpy >= 0 && gpx < g.grid.RawMatrix().Cols && gpy < g.grid.RawMatrix().Rows {
		g.grid.Set(gpx, gpy, g.grid.At(gpx, gpy)+amount)
	}
}

func (g *Grid) Sense(at pixel.Vec, radius float64) float64 {
	gridRadius := int(math.Ceil(radius / float64(g.chunkSie)))
	gpx := int(math.Round(at.X / float64(g.chunkSie)))
	gpy := int(math.Round(at.Y / float64(g.chunkSie)))
	// TODO this is not a true circle
	total := 0.0
	n := 0
	for x := gpx - gridRadius; x <= gpx+gridRadius; x++ {
		if x < 0 || x >= g.grid.RawMatrix().Cols {
			continue
		}
		for y := gpy - gridRadius; y <= gpy+gridRadius; y++ {
			if y < 0 || y >= g.grid.RawMatrix().Rows {
				continue
			}
			total += g.grid.At(x, y)
			n++
		}
	}
	if n == 0 {
		return 0
	}
	return total / float64(n)
}

func (g *Grid) Sprite() (*pixel.Sprite, pixel.Matrix) {
	img := image.NewRGBA(image.Rect(0, 0, g.grid.RawMatrix().Rows, g.grid.RawMatrix().Cols))
	for x := range g.grid.RawMatrix().Rows {
		for y := range g.grid.RawMatrix().Cols {
			v := g.grid.At(x, y)
			if v > 1 {
				v = 1
			}
			img.Set(x, g.grid.RawMatrix().Cols-1-y, pixel.RGB(v, v, v))
		}
	}
	pic := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pic, pic.Bounds()), pixel.IM.Scaled(pixel.ZV, float64(g.chunkSie)).Moved(pixel.V(300, 300))
}
