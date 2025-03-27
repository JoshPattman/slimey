package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/png"
	"math"
	"math/rand"
	"time"

	_ "embed"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

//go:embed particle.png
var particlePng []byte

// Simulation params
var (
	SCREEN_SIZE             pixel.Vec
	NUM_PARTICLES           int
	PARTICLE_SPEED          float64
	PARTICLE_ROTATION_SPEED float64
	PHEREMONE_CHUNK_SIZE    int
	PHEREMONE_RATE          float64
	PHEREMONE_DECAY         float64
	PROFILE                 bool
)

func parseFlags() {
	flag.Float64Var(&SCREEN_SIZE.X, "screen-x", 800, "Screen size X")
	flag.Float64Var(&SCREEN_SIZE.Y, "screen-y", 800, "Screen size Y")
	flag.IntVar(&NUM_PARTICLES, "num", 20000, "Number of particles to simulate")
	flag.Float64Var(&PARTICLE_SPEED, "particle-speed", 100.0, "The speed of each particle")
	flag.Float64Var(&PARTICLE_ROTATION_SPEED, "particle-rotation", 1.5, "The rotational speed of each particle")
	flag.IntVar(&PHEREMONE_CHUNK_SIZE, "chunk-size", 4, "The size in pixels of the pheremone chunks")
	flag.Float64Var(&PHEREMONE_RATE, "pher-rate", 1.0, "rate of dropping pheremones")
	flag.Float64Var(&PHEREMONE_DECAY, "pher-decay", 0.99, "decay rate of pheremones")
	flag.BoolVar(&PROFILE, "prof", false, "periodically print debug statistics out?")

	flag.Parse()
}

func main() {
	parseFlags()
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, SCREEN_SIZE.X, SCREEN_SIZE.Y),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	particles := make([]Particle, NUM_PARTICLES)

	for i := range particles {
		particles[i].Pos = SCREEN_SIZE.Scaled(0.5).Add(pixel.V(rand.Float64()*200, 0).Rotated(rand.Float64() * 2 * math.Pi))
		particles[i].Dir = pixel.V(0, 1).Rotated(rand.Float64() * math.Pi * 2)
		particles[i].Speed = PARTICLE_SPEED
		particles[i].RotationSpeed = PARTICLE_ROTATION_SPEED
	}

	particleImg, err := png.Decode(bytes.NewReader(particlePng))
	if err != nil {
		panic(err)
	}
	particlePic := pixel.PictureDataFromImage(particleImg)

	particleBatch := pixel.NewBatch(&pixel.TrianglesData{}, particlePic)
	particleSprite := pixel.NewSprite(particlePic, particlePic.Rect)

	updatingTime := time.Duration(0)
	drawingTime := time.Duration(0)
	renderingTime := time.Duration(0)

	grid := NewGrid(int(SCREEN_SIZE.X), int(SCREEN_SIZE.Y), PHEREMONE_CHUNK_SIZE)

	step := 0

	deltaTime := 1.0 / 60.0

	for !win.Closed() {
		win.Update()

		startUpdating := time.Now()
		grid.Decay(PHEREMONE_DECAY)
		for i := range particles {
			particles[i].UpdateDir(grid, deltaTime)
			particles[i].UpdatePos(deltaTime, SCREEN_SIZE.X, SCREEN_SIZE.Y)
			grid.AddPheremone(particles[i].Pos, PHEREMONE_RATE*deltaTime)
		}
		updatingTime += time.Since(startUpdating)

		startDrawing := time.Now()
		bgSprite, bgMat := grid.Sprite()

		particleBatch.Clear()
		for _, p := range particles {
			particleSprite.Draw(particleBatch, pixel.IM.Scaled(pixel.ZV, 0.05).Moved(p.Pos))
		}
		drawingTime += time.Since(startDrawing)

		startRendering := time.Now()
		win.Clear(pixel.RGB(0.1, 0.1, 0.1))
		bgSprite.Draw(win, bgMat.Moved(pixel.V(100, 100)))
		particleBatch.Draw(win)
		renderingTime += time.Since(startRendering)

		if step%120 == 0 {
			if PROFILE {
				fmt.Println("per frame:", "update", updatingTime/120, "draw", drawingTime/120, "render", renderingTime/120)
			}
			updatingTime = 0
			renderingTime = 0
			drawingTime = 0
		}
		step++
	}
}
